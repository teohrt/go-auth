package userService

import (
	"context"
	"recollection/repositories"
	"recollection/services/authService"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/rs/zerolog"
)

type Service interface {
	ConfirmUserRegistration(ctx *context.Context, input *authService.RegistrationConfirmationInputBody) error
	Login(ctx *context.Context, input *authService.LoginInputBody) (*cognito.AuthenticationResultType, error)
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

func (svc serviceImpl) Login(ctx *context.Context, input *authService.LoginInputBody) (*cognito.AuthenticationResultType, error) {
	auth_res, err := svc.auth.Login(input)
	if err != nil {
		svc.logger.Error().Err(err).Msg("failed auth login")
		return nil, err
	}
	// token := *res.AccessToken

	user, err := svc.dbRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		svc.logger.Error().Err(err).Msg("failed retrieving user by username")
		return nil, err
	}

	svc.logger.Info().Msgf("%+v", user)

	// check if this is first login (query user auth data to check if exists)
	// if not, parse token and update email and authid

	return auth_res, nil
}
