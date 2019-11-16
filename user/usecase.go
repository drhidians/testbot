package user

import (
	"context"

	"github.com/drhidians/testbot/models"
)

// Usecase represent the user's usecases
type Usecase interface {
	Store(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
}
