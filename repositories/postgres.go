package repositories

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type PostgresRepo interface {
	CreateUser(input DBUser) error
}

type repoImpl struct {
	Client *pgx.Conn
	Logger *zerolog.Logger
}

func NewPostgresRepo(logger *zerolog.Logger) (PostgresRepo, error) {
	dbClient, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error().Err(err).Msg("Database connection failed")
		return nil, err
	}
	defer dbClient.Close(context.Background())

	return repoImpl{
		Client: dbClient,
		Logger: logger,
	}, nil
}

func (repo repoImpl) CreateUser(input DBUser) error {
	// TODO
	return nil
}
