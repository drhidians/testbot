package telegram

// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT CHANGE IT

import "context"

// https://core.telegram.org/bots/api#answercallbackquery
func (b *bot) AnswerCallbackQuery(ctx context.Context, v *CallbackQueryAnswer) error {
	var ok bool
	if err := b.do(ctx, "answerCallbackQuery", v, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotAnswered
	}
	return nil
}

// https://core.telegram.org/bots/api#deletemessage
func (b *bot) DeleteMessage(ctx context.Context, v *DeletedMessage) error {
	var ok bool
	if err := b.do(ctx, "deleteMessage", v, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotAnswered
	}
	return nil
}
