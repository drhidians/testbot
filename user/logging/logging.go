package user

import (
	"github.com/go-kit/kit/log"
)

type logging struct {
	logger log.Logger
	next   Usecase
}

// NewLogging returns a new instance of a logging Service.
func NewLogging(logger log.Logger, s Usecase) Usecase {
	return &logging{logger, s}
}


func (s *logging) GetByID(ctx context.Context, id UserID) (err error){
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

