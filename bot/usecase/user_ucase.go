package usecase

import (
	"context"
	"time"

	"github.com/dhidians/testbot/models"

	"github.com/dhidians/testbot/user"

	"github.com/dhidians/testbot/bot"
)

type botUsecase struct {
	userRepo       user.Repository
	botRepo        bot.Repository
	contextTimeout time.Duration
}

// NewBotUsecase will create new an botUsecase object representation of bot.Usecase interface
func NewBotUsecase(ur bot.Repository, br bot.Repository, timeout time.Duration) bot.Usecase {
	return &botUsecase{
		botRepo:        ur,
		botRepo:        br,
		contextTimeout: timeout,
	}
}

func (a *botUsecase) Store(c context.Context, u *models.Bot) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.botRepo.Store(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (a *botUsecase) Get(c context.Context) (*models.Bot, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	resBot, err := a.botRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	res.Bot = *resBot
	return res, nil
}
