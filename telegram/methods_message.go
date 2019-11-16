package telegram

// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT CHANGE IT

import "context"

// https://core.telegram.org/bots/api#sendmessage
func (b *bot) SendMessage(ctx context.Context, m *TextMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendMessage", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#forwardmessage
func (b *bot) ForwardMessage(ctx context.Context, m *ForwardedMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "forwardMessage", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#sendphoto
func (b *bot) SendPhoto(ctx context.Context, m *PhotoMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendPhoto", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#sendaudio
func (b *bot) SendAudio(ctx context.Context, m *AudioMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendAudio", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#senddocument
func (b *bot) SendDocument(ctx context.Context, m *DocumentMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendDocument", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#sendsticker
func (b *bot) SendSticker(ctx context.Context, m *StickerMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendSticker", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#editmessagetext
func (b *bot) EditMessageText(ctx context.Context, m *MessageText) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "editMessageText", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#editmessagecaption
func (b *bot) EditMessageCaption(ctx context.Context, m *MessageCaption) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "editMessageCaption", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// https://core.telegram.org/bots/api#editmessagereplymarkup
func (b *bot) EditMessageReplyMarkup(ctx context.Context, m *MessageReplyMarkup) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "editMessageReplyMarkup", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}
