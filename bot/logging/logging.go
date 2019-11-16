package user

import (
	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	next   Usecase
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Usecase) Usecase {
	return &loggingService{logger, s}
}


func (s *loggingService) GetByID(ctx context.Context, id UserID) (err error){
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_user_by_id",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetByID(ctx context.Context, id UserID)
}

