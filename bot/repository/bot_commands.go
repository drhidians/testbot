package repository

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	tg "github.com/drhidians/testbot/telegram"
)

var callCommand func(u *tg.Update) error

func (b *botRepository) sendToken() tg.CommandFunc {

	return func(c *tg.Command, u *tg.Update) error {

		mapClaims := jwt.MapClaims{"username": u.Message.From.Username,
			"id":   u.Message.From.ID,
			"name": u.Message.From.FirstName,
			"iat":  time.Now().Unix(),
		}

		if u.Message.From.Username != nil {
			mapClaims["username"] = *u.Message.From.Username
		}

		if u.Message.From.LastName != nil {
			mapClaims["name"] = u.Message.From.FirstName + *u.Message.From.LastName
		}

		_, tokenString, err := b.tokenAuth.Encode(mapClaims)

		if err != nil {
			return err
		}
		_, err = b.botAPI.SendMessage(context.Background(), &tg.TextMessage{
			ChatID: u.Message.Chat.ID,
			Text:   tokenString,
		})

		return err
	}
}
