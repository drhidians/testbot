package usecase

import (
	"context"
	"time"

	"github.com/drhidians/testbot/models"
	"github.com/drhidians/testbot/user"

	"github.com/drhidians/testbot/bot"
	tg "github.com/drhidians/testbot/telegram"
)

// Service represent the bot's usecases
type Service interface {
	Update(context.Context, tg.Update) error
	Get(context.Context) (*models.Bot, error)
}

type botUsecase struct {
	userRepo       user.Repository
	botRepo        bot.Repository
	contextTimeout time.Duration
}

// NewBotService will create new an botUsecase object representation of bot.Usecase interface
func NewBotService(ur user.Repository, br bot.Repository, timeout time.Duration) Service {

	b := &botUsecase{
		userRepo:       ur,
		botRepo:        br,
		contextTimeout: timeout,
	}

	return b

}

func (b *botUsecase) Get(c context.Context) (bot *models.Bot, err error) {

	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	return b.botRepo.Get(ctx)
}

func (b *botUsecase) Update(c context.Context, upd tg.Update) (err error) {

	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	user, err := b.botRepo.Update(ctx, upd)

	if err != nil {
		return err
	}

	err = b.userRepo.Store(ctx, user)
	if err != nil {
		return err
	}

	return err
}
