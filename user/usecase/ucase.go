package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/drhidians/testbot/models"

	"github.com/drhidians/testbot/user"

	"github.com/drhidians/testbot/bot"
)

var ErrUserNotFound = errors.New("invalid credentials")

// Service represent the user's usecases
type Service interface {
	Store(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetBot(c context.Context) (*models.Bot, error)
}

type userUsecase struct {
	userRepo       user.Repository
	botRepo        bot.Repository
	contextTimeout time.Duration
}

// NewUserService will create new an userUsecase object representation of user.Usecase interface
func NewUserService(ur user.Repository, br bot.Repository, timeout time.Duration) Service {
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

	return resUser, nil
}

func (a *userUsecase) GetBot(c context.Context) (*models.Bot, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	resUser, err := a.botRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return resUser, nil
}
