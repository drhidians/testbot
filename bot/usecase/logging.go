package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/drhidians/testbot/models"
	tg "github.com/drhidians/testbot/telegram"
	"github.com/go-kit/kit/log"
)

type logging struct {
	logger log.Logger
	next   Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &logging{logger, s}
}

func (s *logging) Update(ctx context.Context, upd tg.Update) (err error) {

	defer func(begin time.Time) {

		var message *tg.Message
		var text string
		var messageID int

		if upd.Message != nil {
			message = upd.Message
			messageID = message.MessageID
		}

		if message.Text != nil {
			text = *upd.Message.Text
		}

		s.logger.Log(
			"method", "ineract_bot",
			"mesage", text,
			"mesageID", messageID,
			"userID", strconv.FormatInt(upd.Message.Chat.ID, 10),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, upd)
}

func (s *logging) Get(ctx context.Context) (bot *models.Bot, err error) {

	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_bot",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.next.Get(ctx)
}
