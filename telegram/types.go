package telegram

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

// Getting updates
// https://core.telegram.org/bots/api#getting-updates

// https://core.telegram.org/bots/api#update
type Update struct {
	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message"`
	EditedMessage     *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
	// InlineQuery
	// ChosenInlineResult
	CallbackQuery *CallbackQuery `json:"callback_query"`
	// ShippingQuery
	// PreCheckoutQuery
}

// https://core.telegram.org/bots/api#webhookinfo
type WebhookInfo struct {
	URL                  string   `json:"url"`
	HasCustomCertificate bool     `json:"has_custom_certificate"`
	PendingUpdateCount   int      `json:"pending_update_count"`
	LastErrorDate        int      `json:"last_error_date"`
	LastErrorMessage     string   `json:"last_error_message"`
	MaxConnections       int      `json:"max_connections"`
	AllowedUpdates       []string `json:"allowed_updates"`
}

// Available types
// https://core.telegram.org/bots/api#available-types

// https://core.telegram.org/bots/api#user
type User struct {
	ID           int     `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     *string `json:"last_name"`
	Username     *string `json:"username"`
	LanguageCode *string `json:"language_code"`
}

// https://core.telegram.org/bots/api#chat
type Chat struct {
	ID                          int64      `json:"id"`
	Type                        string     `json:"type"`
	Title                       *string    `json:"title"`
	Username                    *string    `json:"username"`
	FirstName                   *string    `json:"first_name"`
	LastName                    *string    `json:"last_name"`
	AllMembersAreAdministrators *bool      `json:"all_members_are_administrators"`
	ChatPhoto                   *ChatPhoto `json:"chat_photo"`
	Description                 *string    `json:"description"`
	InviteLink                  *string    `json:"invite_link"`
}

func (c *Chat) IsPrivate() bool    { return c.Type == "private" }
func (c *Chat) IsGroup() bool      { return c.Type == "group" }
func (c *Chat) IsSupergroup() bool { return c.Type == "supergroup" }
func (c *Chat) IsChannel() bool    { return c.Type == "channel" }

// https://core.telegram.org/bots/api#message
type Message struct {
	MessageID       int              `json:"message_id"`
	From            *User            `json:"from"`
	Date            int              `json:"date"`
	Chat            Chat             `json:"chat"`
	ForwardFrom     *User            `json:"forward_from"`
	ForwardFromChat *Chat            `json:"forward_from_chat"`
	ForwardDate     *int             `json:"forward_date"`
	ReplyToMessage  *Message         `json:"reply_to_message"`
	EditDate        *int             `json:"edit_date"`
	Text            *string          `json:"text"`
	Entities        []*MessageEntity `json:"entities"`
	Audio           *Audio           `json:"audio"`
	Document        *Document        `json:"document"`
	// Game
	Photo                 []*PhotoSize `json:"photo"`
	Sticker               *Sticker     `json:"sticker"`
	Video                 *Video       `json:"video"`
	Voice                 *Voice       `json:"voice"`
	VideoNote             *VideoNote   `json:"video_note"`
	NewChatMembers        []*User      `json:"new_chat_members"`
	Caption               *string      `json:"caption"`
	Contact               *Contact     `json:"contact"`
	Location              *Location    `json:"location"`
	Venue                 *Venue       `json:"venue"`
	NewChatMember         *User        `json:"new_chat_member"`
	LeftChatMember        *User        `json:"left_chat_member"`
	NewChatTitle          *string      `json:"new_chat_title"`
	NewChatPhoto          []*PhotoSize `json:"new_chat_photo"`
	DeleteChatPhoto       *bool        `json:"delete_chat_photo"`
	GroupChatCreated      *bool        `json:"group_chat_created"`
	SupergroupChatCreated *bool        `json:"supergroup_chat_created"`
	ChannelChatCreated    *bool        `json:"channel_chat_created"`
	MigrateToChatID       *int64       `json:"migrate_to_chat_id"`
	MigrateFromChatID     *int64       `json:"migrate_from_chat_id"`
	PinnedMessage         *Message     `json:"pinned_message"`
	// Invoice
	// SuccessfulPayment
}

// https://core.telegram.org/bots/api#messageentity
type MessageEntity struct {
	Type   string  `json:"type"`
	Offset int     `json:"offset"`
	Length int     `json:"length"`
	URL    *string `json:"url"`
	User   *User   `json:"user"`
}

func (e *MessageEntity) IsMention() bool    { return e.Type == "mention" }
func (e *MessageEntity) IsHashtag() bool    { return e.Type == "hastag" }
func (e *MessageEntity) IsBotCommand() bool { return e.Type == "bot_command" }
func (e *MessageEntity) IsURL() bool        { return e.Type == "url" }
func (e *MessageEntity) IsEmail() bool      { return e.Type == "email" }

// https://core.telegram.org/bots/api#photosize
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize *int   `json:"file_size"`
}

// https://core.telegram.org/bots/api#audio
type Audio struct {
	FileID    string  `json:"file_id"`
	Duration  int     `json:"duration"`
	Performer *string `json:"performer"`
	Title     *string `json:"title"`
	MimeType  *string `json:"mime_type"`
	FileSize  *int    `json:"file_size"`
}

// https://core.telegram.org/bots/api#document
type Document struct {
	FileID   string     `json:"file_id"`
	Thumb    *PhotoSize `json:"thumb"`
	FileName *string    `json:"file_name"`
	MimeType *string    `json:"mime_type"`
	FileSize *int       `json:"file_size"`
}

// https://core.telegram.org/bots/api#video
type Video struct {
	FileID   string     `json:"file_id"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb"`
	MimeType *string    `json:"mime_type"`
	FileSize *int       `json:"file_size"`
}

// https://core.telegram.org/bots/api#voice
type Voice struct {
	FileID   string  `json:"file_id"`
	Duration int     `json:"duration"`
	MimeType *string `json:"mime_type"`
	FileSize *int    `json:"file_size"`
}

// https://core.telegram.org/bots/api#videonote
type VideoNote struct {
	FileID   string     `json:"file_id"`
	Length   int        `json:"length"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb"`
	FileSize *int       `json:"file_size"`
}

// https://core.telegram.org/bots/api#contact
type Contact struct {
	PhoneNumber string
	FirstName   string
	LastName    *string
	UserID      *int
}

// https://core.telegram.org/bots/api#location
type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// https://core.telegram.org/bots/api#venue
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID *string  `json:"foursquare_id"`
}

// https://core.telegram.org/bots/api#userprofilephotos
type UserProfilePhotos struct {
	TotalCount int            `json:"total_count"`
	Photos     [][]*PhotoSize `json:"photos"`
}

// https://core.telegram.org/bots/api#file
type File struct {
	FileID   string  `json:"file_id"`
	FileSize *int    `json:"file_size"`
	FilePath *string `json:"file_path"`
}

// https://core.telegram.org/bots/api#replykeyboardmarkup
type ReplyKeyboardMarkup struct {
	Keyboard        [][]*KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool                `json:"one_time_keyboard,omitempty"`
	Selective       bool                `json:"selective,omitempty"`
}

// https://core.telegram.org/bots/api#keyboardbutton
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

// https://core.telegram.org/bots/api#replykeyboardremove
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard,omitempty"`
	Selective      bool `json:"selective,omitempty"`
}

// https://core.telegram.org/bots/api#inlinekeyboardmarkup
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

// https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url,omitempty"`
	CallbackData                 string `json:"callback_data,omitempty"`
	SwitchInlineQuery            string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
	// CallbackGame
	// Pay
}

// https://core.telegram.org/bots/api#callbackquery
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID *string  `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            *string  `json:"data"`
	// GameShortName *string
}

// https://core.telegram.org/bots/api#forcereply
type ForceReply struct {
	ForeceReply bool `json:"force_reply"`
	Selective   bool `json:"selective"`
}

// https://core.telegram.org/bots/api#chatphoto
type ChatPhoto struct {
	SmallFileID string `json:"small_file_id"`
	BigFileID   string `json:"big_file_id"`
}

// https://core.telegram.org/bots/api#chatmember
type ChatMember struct {
	User      User   `json:"user"`
	Status    string `json:"status"`
	UntilDate *int   `json:"until_date"`
	// Administrators only.
	CanBeEdited           *bool `json:"can_be_edited"`
	CanChangeInfo         *bool `json:"can_change_info"`
	CanPostMessages       *bool `json:"can_post_messages"`
	CanEditMessages       *bool `json:"can_edit_messages"`
	CanDeleteMessages     *bool `json:"can_delete_messages"`
	CanInviteUsers        *bool `json:"can_invite_users"`
	CanRestrictMembers    *bool `json:"can_restrict_members"`
	CanPinMessages        *bool `json:"can_pin_messages"`
	CanPromoteMembers     *bool `json:"can_promote_members"`
	CanSendMessages       *bool `json:"can_send_messages"`
	CanSendMediaMessages  *bool `json:"can_send_media_messages"`
	CanSendOtherMessages  *bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews *bool `json:"can_add_web_page_previews"`
}

// https://core.telegram.org/bots/api#responseparameters
type ResponseParameters struct {
	MigrateToChatID *int64 `json:"migrate_to_chat_id"`
	RetryAfter      *int   `json:"retry_after"`
}

// https://core.telegram.org/bots/api#inputfile
type InputFile interface {
	io.Reader
	Name() string
}

// Parse modes.
const (
	ModeDefault  ParseMode = 0
	ModeMarkdown           = 1
	ModeHTML               = 2
)

type ParseMode int

// MarshalJSON implements json.Marshaler interface.
func (m ParseMode) MarshalJSON() (b []byte, err error) {
	switch m {
	case ModeDefault:
		// ModeDefault should be evaluated as empty by json package and skipped
		// in message marshalling. But empty string was a simplest solution to
		// leave ParseMode a simple type.
		b = []byte(`""`)
	case ModeMarkdown:
		b = []byte(`"Markdown"`)
	case ModeHTML:
		b = []byte(`"HTML"`)
	}
	return
}

// https://core.telegram.org/bots/api#sendmessage
type TextMessage struct {
	ChatID                int64     `json:"chat_id"`
	Text                  string    `json:"text"`
	ParseMode             ParseMode `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool      `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool      `json:"disable_notification,omitempty"`
	ReplyToMessageID      int       `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           Markup    `json:"reply_markup,omitempty"`
}

type Markup interface {
	json.Marshaler
	json.Unmarshaler
}

var _ Markup = (*ReplyKeyboardMarkup)(nil)
var _ Markup = (*ReplyKeyboardRemove)(nil)
var _ Markup = (*InlineKeyboardMarkup)(nil)
var _ Markup = (*ForceReply)(nil)

// https://core.telegram.org/bots/api#forwardmessage
type ForwardedMessage struct {
	ChatID              int64 `json:"chat_id"`
	FromChatID          int64 `json:"from_chat_id"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	MessageID           int   `json:"message_id"`
}

// https://core.telegram.org/bots/api#sendphoto
type PhotoMessage struct {
	ChatID              int64     `json:"chat_id"`
	Photo               InputFile `json:"-"`
	PhotoID             string    `json:"photo,omitempty"`
	Caption             string    `json:"caption,omitempty"`
	DisableNotification bool      `json:"disable_notification,omitempty"`
	ReplyToMessageID    int       `json:"reply_to_message_id,omitempty"`
}

// Multipart implements Multiparter interface.
func (m *PhotoMessage) Multipart() *Multipart {
	if m.Photo == nil {
		return nil
	}
	return &Multipart{
		Files: map[string]InputFile{"photo": m.Photo},
		// TODO: Do not include optional fields.
		Form: url.Values{
			"chat_id":              {strconv.FormatInt(m.ChatID, 10)},
			"caption":              {m.Caption},
			"disable_notification": {strconv.FormatBool(m.DisableNotification)},
			"reply_to_message_id":  {strconv.FormatInt(int64(m.ReplyToMessageID), 10)},
		},
	}
}

// https://core.telegram.org/bots/api#sendaudio
type AudioMessage struct {
	ChatID              int64     `json:"chat_id"`
	Audio               InputFile `json:"-"`
	AudioID             string    `json:"audio"`
	Caption             string    `json:"caption,omitempty"`
	Duration            int       `json:"duration,omitempty"`
	Performer           string    `json:"performer,omitempty"`
	Title               string    `json:"title,omitempty"`
	DisableNotification bool      `json:"disable_notification,omitempty"`
	ReplyToMessageID    int       `json:"reply_to_message_id,omitempty"`
	ReplyMarkup         Markup    `json:"reply_markup,omitempty"`
}

// Multipart implements Multiparter interface.
func (m *AudioMessage) Multipart() *Multipart {
	if m.Audio == nil {
		return nil
	}
	var markup string
	if m.ReplyMarkup != nil {
		if b, err := m.ReplyMarkup.MarshalJSON(); err == nil {
			markup = string(b)
		}
	}
	return &Multipart{
		Files: map[string]InputFile{"audio": m.Audio},
		// TODO: Do not include optional fields.
		Form: url.Values{
			"chat_id":              {strconv.FormatInt(m.ChatID, 10)},
			"caption":              {m.Caption},
			"duration":             {strconv.FormatInt(int64(m.Duration), 10)},
			"performer":            {m.Performer},
			"title":                {m.Title},
			"disable_notification": {strconv.FormatBool(m.DisableNotification)},
			"reply_to_message_id":  {strconv.FormatInt(int64(m.ReplyToMessageID), 10)},
			"reply_markup":         {markup},
		},
	}
}

// https://core.telegram.org/bots/api#senddocument
type DocumentMessage struct {
	ChatID              int64     `json:"chat_id"`
	Document            InputFile `json:"document"`
	Caption             string    `json:"caption,omitempty"`
	DisableNotification bool      `json:"disable_notification,omitempty"`
	ReplyToMessageID    int       `json:"reply_to_message_id,omitempty"`
	ReplyMarkup         Markup    `json:"reply_markup,omitempty"`
}

// Multipart implements Multiparter interface.
func (m *DocumentMessage) Multipart() *Multipart {
	if m.Document == nil {
		return nil
	}
	var markup string
	if m.ReplyMarkup != nil {
		if b, err := m.ReplyMarkup.MarshalJSON(); err == nil {
			markup = string(b)
		}
	}
	return &Multipart{
		Files: map[string]InputFile{"document": m.Document},
		// TODO: Do not include optional fields.
		Form: url.Values{
			"chat_id":              {strconv.FormatInt(m.ChatID, 10)},
			"caption":              {m.Caption},
			"disable_notification": {strconv.FormatBool(m.DisableNotification)},
			"reply_to_message_id":  {strconv.FormatInt(int64(m.ReplyToMessageID), 10)},
			"reply_markup":         {markup},
		},
	}
}

// VideoMessage
// VoiceMessage
// VideoNoteMessage

var _ = (Multiparter)((*PhotoMessage)(nil))
var _ = (Multiparter)((*AudioMessage)(nil))
var _ = (Multiparter)((*DocumentMessage)(nil))

// var _ = (Multiparter)((*VideoMessage)(nil))
// var _ = (Multiparter)((*VoiceMessage)(nil))
// var _ = (Multiparter)((*VideoNoteMessage)(nil))

// LocationMessage
// VenueMessage
// ContactMessage

// ChatActionMessage

// UserProfilePhotosMessage

// KickChatMemberMessage
// UnbanChatMemberMessage
// RestrictChatMemberMessage
// PromoteChatMemberMessage
// ExportChatInviteLinkMessage
// SetChatPhotoMessage
// DeleteChatPhotoMessage
// SetChatTitleMessage
// SetChatDescriptionMessage
// PinChatMessageMessage
// UnpinChatMessageMessage
// LeaveChatMessage

// getChat
// getChatAdministrators
// getChatMembersCount
// getChatMember

// TODO: Replace
type ChatID int64

// https://core.telegram.org/bots/api#answercallbackquery
type CallbackQueryAnswer struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

// Updating messages
// https://core.telegram.org/bots/api#updating-messages

// https://core.telegram.org/bots/api#editmessagetext
type MessageText struct {
	ChatID                int64                 `json:"chat_id,omitempty"`
	MessageID             int                   `json:"message_id,omitempty"`
	InlineMessageID       int                   `json:"inline_message_id,omitempty"`
	Text                  string                `json:"text"`
	ParseMode             ParseMode             `json:"parse_mode"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// https://core.telegram.org/bots/api#editmessagecaption
type MessageCaption struct {
	ChatID          int64                 `json:"chat_id,omitempty"`
	MessageID       int                   `json:"message_id,omitempty"`
	InlineMessageID int                   `json:"inline_message_id,omitempty"`
	Caption         string                `json:"caption,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// https://core.telegram.org/bots/api#editmessagereplymarkup
type MessageReplyMarkup struct {
	ChatID          int64                 `json:"chat_id,omitempty"`
	MessageID       int                   `json:"message_id,omitempty"`
	InlineMessageID int                   `json:"inline_message_id,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// https://core.telegram.org/bots/api#deletemessage
type DeletedMessage struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}

// Stickers
// https://core.telegram.org/bots/api#stickers

// https://core.telegram.org/bots/api#sticker
type Sticker struct {
	FileID       string        `json:"file_id"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	Thumb        *PhotoSize    `json:"thumb"`
	Emoji        *string       `json:"emoji"`
	SetName      *string       `json:"set_name"`
	MaskPosition *MaskPosition `json:"mask_position"`
	FileSize     *int          `json:"file_size"`
}

// https://core.telegram.org/bots/api#stickerset
type StickerSet struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	ContainsMasks bool       `json:"contains_masks"`
	Stickers      []*Sticker `json:"stickers"`
}

// https://core.telegram.org/bots/api#maskposition
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

// https://core.telegram.org/bots/api#sendsticker
type StickerMessage struct {
	ChatID              int64     `json:"chat_id"`
	Sticker             InputFile `json:"-"`
	StickerID           string    `json:"sticker"`
	DisableNotification bool      `json:"disalbe_notification,omitempty"`
	ReplyToMessageID    int       `json:"reply_to_message_id,omitempty"`
	ReplyMarkup         Markup    `json:"reply_markup,omitempty"`
}

// Inline mode
// https://core.telegram.org/bots/api#inline-mode
// TODO: Add types and methods.
