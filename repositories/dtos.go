package repositories

import "github.com/jackc/pgx/v5/pgtype"

type DBUser struct {
	UserID     pgtype.Text      `db:"user_id"`
	Username   pgtype.Text      `db:"username"`
	AuthID     pgtype.Text      `db:"cognito_id"`
	Email      pgtype.Text      `db:"email"`
	Created_at pgtype.Timestamp `db:"created_at"`
	Updated_at pgtype.Timestamp `db:"updated_at"`
}

type CreateUserInput struct {
	AuthID   string
	Username string
}
