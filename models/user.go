package models

// User represent the user model
type User struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	CreatedAt  int64   `json:"joinedAt"`
	UpdatedAt  *int64  `json:"updatedAt"`
	ExternalID int     `json:"externalId"`
	Username   *string `json:"username"`
	Avatar     *string `json:"avatar"`
	Language   *string `json:"language"`
}
