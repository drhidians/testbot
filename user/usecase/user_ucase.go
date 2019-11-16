package usecase

import (
	"context"
	"time"

	"github.com/dhidians/testbot/models"

	"github.com/dhidians/testbot/user"

	"github.com/dhidians/testbot/bot"
)

type userUsecase struct {
	userRepo       user.Repository
	botRepo        bot.Repository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewUserUsecase(ur user.Repository, br bot.Repository, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:       ur,
		botRepo:        br,
		contextTimeout: timeout,
	}
}

func (a *userUsecase) Store(c context.Context, u *models.User) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.userRepo.Store(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (a *userUsecase) GetByID(c context.Context, id int64) (*models.User, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	resUser, err := a.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res.User = *resUser
	return res, nil
}
