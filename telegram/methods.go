package telegram

import (
	"context"
	"errors"
)

var (
	ErrNotDeleted  = errors.New("telegram: message not deleted")
	ErrNotEdited   = errors.New("telegram: message not edited")
	ErrNotAnswered = errors.New("telegram: query not answered")
)

// https://core.telegram.org/bots/api#getupdates
func (b *bot) GetUpdates(ctx context.Context, opts ...UpdatesOption) ([]*Update, error) {
	uo := new(updatesOptions)
	for _, opt := range opts {
		opt(uo)
	}
	var v []*Update
	if err := b.do(ctx, "getUpdates", uo, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// NewWebhook
func (b *bot) NewWebhook(link string, max_connection int) WebhookConfig {

	return WebhookConfig{
		URL:            link,
		MaxConnections: max_connection,
	}
}

// SetWebhook
func (b *bot) SetWebhook(ctx context.Context, config WebhookConfig) error {

	return b.do(ctx, "setWebhook", config, nil)
}

// DeleteWebhook

// GetWebhookInfo
func (b *bot) GetWebhookInfo(ctx context.Context) (*WebhookInfo, error) {

	var w *WebhookInfo
	err := b.do(ctx, "getWebhookInfo", nil, &w)
	if err != nil {
		return nil, err
	}

	return w, err
}

// WebhookInfo

// https://core.telegram.org/bots/api#getme
func (b *bot) GetMe(ctx context.Context) (*User, error) {
	var u *User
	if err := b.do(ctx, "getMe", nil, &u); err != nil {
		return nil, err
	}
	return u, nil
}

// https://core.telegram.org/bots/api#sendchataction
// func (b *bot) SendChatAction(ctx context.Context, action *ChatAction) error {
// 	var ok bool
// 	if err := b.do(ctx, "answerCallbackQuery", a, &ok); err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return ErrNotAnswered
// 	}
// 	return nil
// }

// GetUserProfilePhotos
func (b *bot) GetUserProfilePhotos(ctx context.Context, config UserProfilePhotosConfig) (*UserProfilePhotos, error) {

	var up *UserProfilePhotos
	err := b.do(ctx, "getUserProfilePhotos", config, &up)
	if err != nil {
		return nil, err
	}

	return up, err
}

// GetFile
func (b *bot) GetFile(ctx context.Context, config FileConfig) (*File, error) {
	var f *File
	err := b.do(ctx, "getFile", config, &f)
	if err != nil {
		return nil, err
	}

	return f, err
}

// GetFileDirectURL
func (b *bot) GetFileDirectURL(ctx context.Context, config FileConfig) (*string, error) {
	file, err := b.GetFile(ctx, config)

	if err != nil {
		return nil, err
	}

	if file.FilePath == nil {
		//TO DO
		return nil, errors.New("empty file path")
	}

	directURL := b.fileurl + "/" + *file.FilePath
	return &directURL, nil
}

// KickChatMember
// UnbanChatMember
// RestrictChatMember
// PromoteChatMember
// ExportChatInviteLink
// SetChatPhoto
// SetChatDescription
// PinChatMessage
// UnpinChatMessage
// LeaveChat
// GetChat
// GetChatAdministrators
// GetChatMembersCount
// GetChatMembers

// TODO: What does True mean for edit* methods?
// > On success, if edited message is sent by the bot, the edited Message is
// > returned, otherwise True is returned.
// https://core.telegram.org/bots/api#editmessagetext
