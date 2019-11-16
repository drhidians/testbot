package user

import (
	"context"

	"github.com/drhidians/testbot/models"
)

// Logging represent the user's usecases
type Logging interface {
	Store(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) error
}
