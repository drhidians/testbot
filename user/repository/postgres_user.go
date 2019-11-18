package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

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
	//TO DO move to utils directory
	stmt, err := Conn.Prepare(`CREATE Table IF NOT EXISTS users(id SERIAL PRIMARY KEY, external_id int, name varchar(50), username varchar(30), language varchar(10), avatar varchar(255), created_at TIMESTAMP with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP with time zone);`)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	return &postgresUserRepository{Conn}
}

func (p *postgresUserRepository) GetByID(ctx context.Context, id int) (user *models.User, err error) {

	user = new(models.User)
	query := `SELECT * FROM users WHERE ID = $1`

	err = p.Conn.QueryRowContext(ctx, query, id).Scan(&user)

	if err != nil {
		return nil, err
	}

	return
}

func (p *postgresUserRepository) GetByTelegramID(ctx context.Context, tgID int) (u *models.User, err error) {

	u = new(models.User)

	query := `SELECT id, Name, created_at, updated_at, external_id, username, avatar, language FROM users WHERE external_id = $1`

	var updatedAt pq.NullTime
	var username sql.NullString
	var avatar sql.NullString
	var language sql.NullString

	var createdAt pq.NullTime
	err = p.Conn.QueryRowContext(ctx, query, tgID).Scan(&u.ID, &u.Name, &createdAt, &updatedAt, &u.ExternalID, &username, &avatar, &language)

	if createdAt.Valid {
		u.CreatedAt = createdAt.Time.Unix()
	}
	if updatedAt.Valid {
		unixT := updatedAt.Time.Unix()
		u.UpdatedAt = &unixT
	}
	if username.Valid {
		u.Username = &username.String
	}
	if avatar.Valid {
		u.Avatar = &avatar.String
	}
	if language.Valid {
		u.Language = &language.String
	}
	return
}

func newNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func (p *postgresUserRepository) Store(ctx context.Context, u *models.User) (err error) {

	_, err = p.GetByTelegramID(ctx, u.ExternalID)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil {
		query := `UPDATE users SET updated_at = CURRENT_TIMESTAMP, language = $2, avatar = $3 WHERE external_id = $1`
		stmt, err := p.Conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		language := newNullString(u.Language)
		avatar := newNullString(u.Avatar)

		_, err = stmt.ExecContext(ctx, u.ExternalID, language, avatar)

		return err
	}

	query := `INSERT INTO users (external_id, username, name, avatar, language, created_at) VALUES ($1 , $2 , $3, $4 , $5, CURRENT_TIMESTAMP)`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	language := newNullString(u.Language)
	avatar := newNullString(u.Avatar)
	username := newNullString(u.Username)
	_, err = stmt.ExecContext(ctx, u.ExternalID, username, u.Name, avatar, language)

	if err != nil {
		return err
	}

	/*lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = lastID*/
	return
}
