package bot

import (
	"context"

	"github.com/drhidians/testbot/models"
)

// Repository represent the bot's repository contract
type Repository interface {
	Store(ctx context.Context, bot *models.Bot) error
	Get(ctx context.Context) (*models.Bot, error)
}
