package bot

import (
	"context"

	"github.com/drhidians/testbot/models"
)

// Usecase represent the bot's usecases
type Usecase interface {
	Store(ctx context.Context, bot *models.Bot) error
	Get(ctx context.Context) error
}
