package bot

import (
	"context"

	tg "github.com/drhidians/testbot"
	"github.com/drhidians/testbot/models"
)

// Repository represent the bot's repository contract
type Repository interface {
	Get(context.Context) (*models.Bot, error)
	Update(context.Context, tg.Update) (*models.User, error)
	GetFile(context.Context, string) ([]byte, error)
}
