package repository

import (
	"context"
	"database/sql"

	"github.com/drhidians/testbot/models"
	"github.com/drhidians/testbot/bot"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type postgresBotRepository struct {
	Conn *sql.DB
}

// NewPostgresBotRepository will create an object that represent the bot.Repository interface
func NewPostgresBotRepository(Conn *sql.DB) bot.Repository {
	return &postgresBotRepository{Conn}
}

func (p *postgresBotRepository) Get(ctx context.Context) (bot *models.Bot, err error) {
	query := `SELECT * FROM bot WHERE ID = 1`

	err := p.Conn.QueryContext(ctx, query).Scan(&bot)

	if err != nil {
		return nil, err
	}

	return
}

func (p *postgresBotRepository) Store(ctx context.Context, b *models.Bot)  err error {
	query := `INSERT  bot SET id=?, username="?", name=?`
	
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	res, err := stmt.ExecContext(ctx, b.ID, b.Name, b.Username)

	return 
}
