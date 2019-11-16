package usecase

import (
	"context"
	"time"

	"github.com/drhidians/testbot/models"

	"github.com/drhidians/testbot/user"

	"github.com/drhidians/testbot/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Service represent the bot's usecases
type Service interface {
	Store(ctx context.Context, bot *models.Bot) error
	Get(ctx context.Context) (*models.Bot, error)
}

type botUsecase struct {
	userRepo       user.Repository
	botRepo        bot.Repository
	contextTimeout time.Duration
	botAPI         *tgbotapi.BotAPI
}

// NewBotService will create new an botUsecase object representation of bot.Usecase interface
func NewBotService(ur user.Repository, br bot.Repository, bAPI *tgbotapi.BotAPI, timeout time.Duration) Service {

	return &botUsecase{
		userRepo:       ur,
		botRepo:        br,
		contextTimeout: timeout,
		botAPI:         bAPI,
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

	_, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	resBot, err := a.botAPI.GetMe()

	if err != nil {
		return nil, err
	}

	var b *models.Bot
	b.ID = resBot.ID
	b.Name = resBot.FirstName + " " + resBot.LastName
	b.Username = resBot.UserName

	return b, nil
}
