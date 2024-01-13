package userService

import (
	"context"
	"recollection/repositories"
	"recollection/services/authService"

	"github.com/rs/zerolog"
)

type Service interface {
	ConfirmUserRegistration(ctx *context.Context, input *authService.RegistrationConfirmationInputBody) error
}

type serviceImpl struct {
	dbRepo repositories.PostgresRepo
	auth   authService.Client
	logger *zerolog.Logger
}

func New(ctx *context.Context, auth authService.Client, logger *zerolog.Logger) (Service, error) {
	dbRepo, err := repositories.NewPostgresRepo(ctx, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed initializing postgres repository")
		return nil, err
	}

	return &serviceImpl{
		dbRepo,
		auth,
		logger,
	}, nil
}

func (svc serviceImpl) ConfirmUserRegistration(ctx *context.Context, input *authService.RegistrationConfirmationInputBody) error {
	if err := svc.auth.ConfirmRegistration(input); err != nil {
		svc.logger.Error().Err(err).Msg("failed confirming registration")
		return err
	}

	dbInput := repositories.CreateUserInput{
		Username: input.Username,
	}
	if err := svc.dbRepo.CreateUser(ctx, dbInput); err != nil {
		svc.logger.Error().Err(err).Msg("failed adding user to db")
		return err
	}
	return nil
}
