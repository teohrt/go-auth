package app

import (
	"log"
	"net/http"
	"os"
	"recollection/handlers/authHandler"
	"recollection/handlers/healthHandler"
	"recollection/services/authService"

	"github.com/caarlos0/env/v10"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func Start() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Error().Err(err).Msg("Failed to load .env")
		return
	}

	PORT := os.Getenv("SERVER_PORT")
	if PORT == "" {
		PORT = "80"
	}

	cognitoConfig := authService.Config{}
	if err := env.Parse(&cognitoConfig); err != nil {
		logger.Error().Err(err).Msg("Failed retrieving AWS Cognito environment variable requirements")
		return
	}

	cognito, err := authService.New(&cognitoConfig, &logger)
	if err != nil {
		logger.Error().Err(err).Msg("Cognito initialization failed")
		return
	}

	router := chi.NewRouter()
	router.Route("/v1", func(subRouter chi.Router) {
		subRouter.Get("/health", healthHandler.Handler())
		subRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/register", authHandler.RegistrationHandler(cognito, &logger))
			authRouter.Post("/login", authHandler.LoginHandler(cognito, &logger))
			authRouter.Post("/confirm", authHandler.ConfirmRegistrationHandler(cognito, &logger))
		})
	})

	logger.Debug().Msg("Server running locally and listening on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
