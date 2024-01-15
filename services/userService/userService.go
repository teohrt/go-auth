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

	user, err := svc.dbRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		svc.logger.Error().Err(err).Msg("failed retrieving user by username")
		return nil, err
	}

	isFirstLogin := !user.AuthID.Valid // || !user.Email.Valid - no way to get email at this point
	svc.logger.Info().Msgf("isFirstLogin %v", isFirstLogin)
	if isFirstLogin {
		svc.logger.Info().Msg("This is the user's first login")
		jwtToken := *auth_res.AccessToken
		_, err := svc.auth.ParseJWTClaims(jwtToken)
		if err != nil {
			svc.logger.Error().Err(err).Msg("failed parsing token")
			return nil, err
		}

		// TODO - update auth id, updated_at timestamp
	}

	return auth_res, nil
}
