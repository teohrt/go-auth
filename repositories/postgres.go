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
	status, err := repo.client.Exec(*ctx, query, args)
	if err != nil {
		repo.logger.Error().Err(err).Msg("Failed inserting row")
		return err
	}
	repo.logger.Debug().Msgf("Inserted user: %+v", status)
	return nil
}

func (repo repoImpl) ReadUser(input DBUser) error {
	// TODO
	return nil
}

func (repo repoImpl) UpdateUser(input DBUser) error {
	// TODO
	return nil
}

func (repo repoImpl) DeleteUser(input DBUser) error {
	// TODO
	return nil
}
