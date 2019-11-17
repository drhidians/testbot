package models

// Bot represent the bot model
type Bot struct {
	ID       int     `json:"id"`
	Username *string `json:"username"`
	Name     string  `json:"name"`
}
