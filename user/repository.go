package user

import (
	"context"

	"github.com/drhidians/testbot/models"
)

// Repository represent the user's repository contract
type Repository interface {
	Store(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByTelegramID(ctx context.Context, id int) (*models.User, error)
}
