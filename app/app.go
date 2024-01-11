package app

import (
	"log"
	"net/http"
	"os"
	"recollection/handlers/authHandler"
	"recollection/handlers/healthHandler"
	"recollection/services/authService"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

func Start() {
	PORT := os.Getenv("SERVER_PORT")
	if PORT == "" {
		PORT = "80"
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cognito, err := authService.New(&authService.Config{
		AWSRegion:   os.Getenv("AWS_REGION"),
		UserPoolID:  os.Getenv("AWS_COGNITO_USER_POOL_ID"),
		AppClientID: os.Getenv("AWS_COGNITO_APP_CLIENT_ID"),
	}, &logger)
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
