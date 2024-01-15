package app

import (
	"context"
	"recollection/services/authService"
	"recollection/services/userService"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
)

type dependencies struct {
	userService userService.Service
	authClient  authService.Client
}

func getDependencies(ctx *context.Context, logger *zerolog.Logger) (*dependencies, error) {
	cognitoConfig := authService.Config{}
	if err := env.Parse(&cognitoConfig); err != nil {
		logger.Error().Err(err).Msg("Failed retrieving AWS Cognito environment variable requirements")
		return nil, err
	}

	authClient, err := authService.New(ctx, &cognitoConfig, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Cognito initialization failed")
		return nil, err
	}

	userService, err := userService.New(ctx, authClient, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Cognito initialization failed")
		return nil, err
	}
	return &dependencies{
		userService,
		authClient,
	}, nil
}
