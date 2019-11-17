package models

import (
	"time"
)

// User represent the user model
type User struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"joinedAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	ExternalID int        `json:"externalId"`
	Username   *string    `json:"username"`
	Avatar     *string    `json:"avatar"`
	Language   *string    `json:"language"`
}
