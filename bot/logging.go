package bot

import (
	"context"

	"github.com/drhidians/testbot/models"
)

// Logging represent the bot's usecases
type Logging interface {
	Store(ctx context.Context, bot *models.Bot) error
	Get(ctx context.Context) (*models.Bot, error)
}
