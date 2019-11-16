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

func (s *logging) Store(ctx context.Context, u *models.Bot) (err error) {

	defer func(begin time.Time) {
		s.logger.Log(
			"method", "store_bot",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Store(ctx, u)
}

func (s *logging) Get(ctx context.Context) (user *models.Bot, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_bot",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Get(ctx)
}
