package repositories

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresRepo interface {
	CreateUser(ctx *context.Context, input CreateUserInput) error
	GetUserByUsername(ctx *context.Context, username string) (*DBUser, error)
}

type repoImpl struct {
	client *pgxpool.Pool
	logger *zerolog.Logger
}

func NewPostgresRepo(ctx *context.Context, logger *zerolog.Logger) (PostgresRepo, error) {
	dbClient, err := pgxpool.New(*ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error().Err(err).Msg("Failed creating connection pool")
		return nil, err
	}

	return repoImpl{
		client: dbClient,
		logger: logger,
	}, nil
}

func (repo repoImpl) CreateUser(ctx *context.Context, input CreateUserInput) error {
	query := `
	INSERT INTO Users (username) 
	VALUES (@username)
	`
	args := pgx.NamedArgs{
		"username": input.Username,
	}
	_, err := repo.client.Exec(*ctx, query, args)
	if err != nil {
		repo.logger.Error().Err(err).Msg("Failed inserting user")
		return err
	}
	return nil
}

func (repo repoImpl) GetUserByUsername(ctx *context.Context, username string) (*DBUser, error) {
	query := `
	SELECT * FROM Users
	WHERE username = @username
	`
	args := pgx.NamedArgs{
		"username": username,
	}
	row, err := repo.client.Query(*ctx, query, args)
	if err != nil {
		repo.logger.Error().Err(err).Msg("Failed query - GetUserByUsername")
		return nil, err
	}
	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[DBUser])
	if err != nil {
		repo.logger.Error().Err(err).Msg("Failed to convert row to struct")
		return nil, err
	}
	return &user, nil
}

func (repo repoImpl) UpdateUser(input DBUser) error {
	// TODO
	return nil
}

func (repo repoImpl) DeleteUser(input DBUser) error {
	// TODO
	return nil
}
