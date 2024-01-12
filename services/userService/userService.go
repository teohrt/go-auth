package userService

import (
	"recollection/repositories"

	"github.com/rs/zerolog"
)

type Service interface {
}

type serviceImpl struct {
	dbRepo repositories.PostgresRepo
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) (Service, error) {
	dbRepo, err := repositories.NewPostgresRepo(logger)
	if err != nil {
		logger.Error().Err(err).Msg("Postgres repository initialization failed")
		return nil, err
	}

	return &serviceImpl{
		dbRepo,
		logger,
	}, nil
}

func (svc serviceImpl) CreateUser() {
	// todo
}

func (svc serviceImpl) UpdateUserAuthId() {
	// todo
}

func (svc serviceImpl) UpdateUserEmail() {
	// todo
}
