package usecase

import (
	"context"
	"time"

	"github.com/drhidians/testbot/models"
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

func (s *logging) Store(ctx context.Context, u *models.User) (err error) {

	defer func(begin time.Time) {
		s.logger.Log(
			"method", "store_user",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Store(ctx, u)
}

func (s *logging) GetByTelegramID(ctx context.Context, tgID int) (user *models.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_user_by_tg_id",
			"id", tgID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetByTelegramID(ctx, tgID)
}

func (s *logging) GetBot(ctx context.Context) (bot *models.Bot, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_bot",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetBot(ctx)
}

func (s *logging) GetAvatar(ctx context.Context, avatarID string) (b []byte, err error) {

	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_user_avatar",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetAvatar(ctx, avatarID)
}
