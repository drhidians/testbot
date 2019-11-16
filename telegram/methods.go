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

// SetWebhook
func (b *bot) SetWebhook(ctx context.Context) (*User, error) {
	var u *User
	if err := b.do(ctx, "getMe", nil, &u); err != nil {
		return nil, err
	}
	return u, nil
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
// GetFile
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
