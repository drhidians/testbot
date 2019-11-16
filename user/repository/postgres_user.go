package repository

import (
	"context"
	"database/sql"

	"github.com/drhidians/testbot/models"
	"github.com/drhidians/testbot/user"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type postgresUserRepository struct {
	Conn *sql.DB
}

// NewPostgresUserRepository will create an object that represent the user.Repository interface
func NewPostgresUserRepository(Conn *sql.DB) user.Repository {
	return &postgresUserRepository{Conn}
}

func (p *postgresUserRepository) GetByID(ctx context.Context, id int64) (user *models.User, err error) {
	query := `SELECT * FROM user WHERE ID = ?`

	err := p.Conn.QueryContext(ctx, query, id).Scan(&user)

	if err != nil {
		return nil, err
	}

	return
}

func (p *postgresUserRepository) Store(ctx context.Context, u *models.User)  err error {
	query := `INSERT  user SET externalId=? , username=? , name=?, avatar=? , language=?, created_at=?, updated_at=?`
	
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	res, err := stmt.ExecContext(ctx,u.ExternalID, u.Username,u.Name,u.Avatar,u.Language,u.JoinedAt,u.UpdatedAt)

	return 
}