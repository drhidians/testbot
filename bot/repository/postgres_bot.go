package repository

import (
	"context"
	"database/sql"

	"github.com/drhidians/testbot/bot"
	"github.com/drhidians/testbot/models"
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

	err = p.Conn.QueryRowContext(ctx, query).Scan(&bot)

	if err != nil {
		return nil, err
	}

	return
}

func (p *postgresBotRepository) Store(ctx context.Context, b *models.Bot) (err error) {
	query := `INSERT  bot SET id=?, username="?", name=?`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, b.ID, b.Name, b.Username)

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	b.ID = int(lastID)
	return
}
