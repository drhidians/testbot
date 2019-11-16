package telegram

//go:generate python methods_bool.py
//go:generate python methods_message.py
//go:generate python types_keyboards.py
//go:generate gofmt -w .

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const jsonContentType = "application/json;chartset=utf-8"

const (
	defaultURL         = "https://api.telegram.org/bot"
	defaultErrTimeout  = 5 * time.Second
	defaultPollTimeout = time.Minute
)

var ErrEmptyToken = errors.New("telegram: empty token")

type Bot interface {
	Username() string
	Updates() <-chan *Update
	Errors() <-chan error

	GetMe(context.Context) (*User, error)
	GetUpdates(context.Context, ...UpdatesOption) ([]*Update, error)

	SendMessage(context.Context, *TextMessage) (*Message, error)
	ForwardMessage(context.Context, *ForwardedMessage) (*Message, error)
	// SendPhoto(context.Context, *PhotoMessage) (*Message, error)
	// SendAudio(context.Context, *AudioMessage) (*Message, error)
	// SendDocument(context.Context, *DocumentMessage) (*Message, error)
	// SendSticker(context.Context, *StickerMessage) (*Message, error)
	// SendVideo(context.Context, *VideoMessage) (*Message, error)
	// SendVoice(context.Context, *VoiceMessage) (*Message, error)
	// SendVoiceNote(context.Context, *VoiceNoteMessage) (*Message, error)
	// SendLocation(context.Context, *LocationMessage) (*Message, error)
	// SendVenue(context.Context, *VenueMessage) (*Message, error)
	// SendContact(context.Context, *ContactMessage) (*Message, error)

	EditMessageText(context.Context, *MessageText) (*Message, error)
	EditMessageCaption(context.Context, *MessageCaption) (*Message, error)
	EditMessageReplyMarkup(context.Context, *MessageReplyMarkup) (*Message, error)
	DeleteMessage(context.Context, *DeletedMessage) error
}

func NewBot(ctx context.Context, token string, opts ...BotOption) (Bot, error) {
	if token == "" {
		return nil, ErrEmptyToken
	}
	b := newBot(ctx, token, opts...)
	if err := b.getUsername(); err != nil {
		return nil, err
	}
	if !b.noUpdates {
		go b.listenToUpdates()
	}
	return b, nil
}

type SOCKS5 struct {
	Address  string // host:port
	User     string
	Password string
}

type botOptions struct {
	Username    string
	URL         string
	ErrTimeout  time.Duration
	PollTimeout time.Duration
	NoUpdates   bool
}

type BotOption func(*botOptions)

func WithUsername(s string) BotOption {
	return func(o *botOptions) {
		o.Username = s
	}
}

func withURL(url string) BotOption {
	return func(o *botOptions) {
		o.URL = url
	}
}

func WithErrTimeout(t time.Duration) BotOption {
	return func(o *botOptions) {
		o.ErrTimeout = t
	}
}

func WithPollTimeout(t time.Duration) BotOption {
	return func(o *botOptions) {
		o.PollTimeout = t
	}
}

func WithoutUpdates() BotOption {
	return func(o *botOptions) {
		o.NoUpdates = true
	}
}

type bot struct {
	username    string
	url         string
	ctx         context.Context
	client      *http.Client
	errTimeout  time.Duration
	pollTimeout time.Duration
	noUpdates   bool
	updatec     chan *Update
	errorc      chan error
}

func newBot(ctx context.Context, token string, opts ...BotOption) *bot {
	o := &botOptions{URL: defaultURL, ErrTimeout: defaultErrTimeout, PollTimeout: defaultPollTimeout}
	for _, opt := range opts {
		opt(o)
	}

	var client *http.Client

	b := &bot{
		username:    o.Username,
		url:         o.URL + token,
		ctx:         ctx,
		client:      client,
		errTimeout:  o.ErrTimeout,
		pollTimeout: o.PollTimeout,
		noUpdates:   o.NoUpdates,
		updatec:     make(chan *Update),
		errorc:      make(chan error),
	}
	if b.noUpdates {
		close(b.updatec)
		close(b.errorc)
	}
	return b
}

func (b *bot) getUsername() error {
	if b.username != "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.errTimeout)
	defer cancel()
	me, err := b.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("telegram: could not get name: %s", err)
	}
	b.username = *me.Username
	return nil
}

func (b *bot) listenToUpdates() {
	var offset int
	donec := b.ctx.Done()
loop:
	for {
		u, err := b.GetUpdates(b.ctx, WithOffset(offset), WithTimeout(b.pollTimeout))
		// Handle context errors differently - shutdown gracefully.
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			break loop
		}

		if err != nil {
			select {
			case b.errorc <- err:
				sleepctx(b.ctx, b.errTimeout)
				continue
			case <-donec:
				break
			}
		}
		// No updates this time - repeat the loop and wait for another pack.
		if len(u) == 0 {
			continue
		}
		// Increment offset according to the last update id. Next time updates
		// pack will not contain updates up to this last one.
		offset = u[len(u)-1].UpdateID + 1

		for _, up := range u {
			select {
			case b.updatec <- up:
				continue
			case <-donec:
				break loop
			}
		}
	}

	// TODO: How to ensure updatesc and errorc to be drained?

	// Don't forget to close channels.
	close(b.updatec)
	close(b.errorc)
}

// sleepctx pauses for at lease t duration. It returns early if ctx is cancelled or
// its deadline is exceeded.
func sleepctx(ctx context.Context, t time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(t):
	}
}

func (b *bot) Username() string { return b.username }

func (b *bot) Updates() <-chan *Update { return b.updatec }
func (b *bot) Errors() <-chan error    { return b.errorc }

// call issues HTTP request to API for the method with form values and decodes
// received data in v. It returns error otherwise.
func (b *bot) do(ctx context.Context, method string, data interface{}, v interface{}) error {
	url := b.url + "/" + method

	body, contentType, err := b.encode(data)
	if err != nil {
		return err
	}
	resp, err := post(ctx, b.client, url, contentType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := new(apiResponse)
	if err := json.Unmarshal(bdata, r); err != nil {
		return err
	}

	if !r.OK {
		return &Error{
			ErrorCode:   r.ErrorCode,
			Description: r.Description,
			Parameters:  r.Parameters,
		}
	}

	return json.Unmarshal([]byte(r.Result), v)
}

// post issues a POST request via the do function.
//
// Copied from golang.org/x/net/context/ctxhttp
func post(ctx context.Context, client *http.Client, url, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return do(ctx, client, req)
}

// do sends an HTTP request with the provided http.Client and returns
// an HTTP response.
//
// If the client is nil, http.DefaultClient is used.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
//
// Copied from golang.org/x/net/context/ctxhttp
func do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req.WithContext(ctx))
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}
	return resp, err
}

func (b *bot) encode(data interface{}) (io.Reader, string, error) {
	if m, ok := data.(Multiparter); ok {
		if v := m.Multipart(); v != nil {
			return v.Encode()
		}
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, "", err
	}
	return buf, jsonContentType, nil
}

// Multiparter is an interface for messages that may be converted to a multipart
// form (e.g. photo, document, video). *Multipart may be nil meaning unavailable
// conversion.
type Multiparter interface {
	Multipart() *Multipart
}

type Multipart struct {
	Form  url.Values
	Files map[string]InputFile
}

// Encode encodes Multipart to multipart/form-data. It returns io.Reader for
// content, content type with boundary and error. In case of failed encoding
// io.Reader is nil, content type is an empty string.
func (m *Multipart) Encode() (io.Reader, string, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for key := range m.Form {
		if err := w.WriteField(key, m.Form.Get(key)); err != nil {
			return nil, "", err
		}
	}

	for key := range m.Files {
		file := m.Files[key]
		if dest, err := w.CreateFormFile(key, file.Name()); err != nil {
			return nil, "", err
		} else {
			if _, err := io.Copy(dest, file); err != nil {
				return nil, "", err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return buf, w.FormDataContentType(), nil
}

type updatesOptions struct {
	Offset  int
	Limit   int
	Timeout time.Duration
	// AllowedUpdates []string
}

// MarshalJSON implements json.Marshaler interface.
func (o *updatesOptions) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if o.Offset > 0 {
		m["offset"] = o.Offset
	}
	if o.Limit > 0 {
		m["limit"] = o.Limit
	}
	if o.Timeout > 0 {
		m["timeout"] = int(o.Timeout.Seconds())
	}
	return json.Marshal(m)
}

type UpdatesOption func(*updatesOptions)

// WithOffset sets id of the first expected update in response. Usually offset
// should equal last update's id + 1.
func WithOffset(offset int) UpdatesOption {
	return func(o *updatesOptions) {
		o.Offset = offset
	}
}

// WithLimit modifies updates request to limit the number of updates in response.
func WithLimit(limit int) UpdatesOption {
	return func(o *updatesOptions) {
		if limit < 1 {
			limit = 1
		}
		if limit > 100 {
			limit = 100
		}
		o.Limit = limit
	}
}

// WithTimeout modifies timeout of updates request. 0 duration means short
// polling (for testing only).
func WithTimeout(t time.Duration) UpdatesOption {
	return func(o *updatesOptions) {
		o.Timeout = t
	}
}

// Error represents an error returned by API. It satisfies error interface.
type Error struct {
	ErrorCode   int
	Description string
	Parameters  *ResponseParameters
}

// Error returns an error string.
func (e *Error) Error() string {
	return fmt.Sprintf("telegram: %d %s", e.ErrorCode, e.Description)
}

// apiResponse represents API response. When OK is false then ErrorCode and
// Description defines the error situation.
type apiResponse struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	// error part
	ErrorCode   int                 `json:"error_code,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}
