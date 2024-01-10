package app

import (
	"log"
	"net/http"
	"os"
	authClient "recollection/authClient"
	"recollection/handlers/auth"
	"recollection/handlers/health"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

func Start() {
	PORT := os.Getenv("SERVER_PORT")
	if PORT == "" {
		PORT = "80"
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cognito, _ := authClient.New(&authClient.Config{
		AWSRegion:   os.Getenv("AWS_REGION"),
		UserPoolID:  os.Getenv("AWS_COGNITO_USER_POOL_ID"),
		AppClientID: os.Getenv("AWS_COGNITO_APP_CLIENT_ID"),
	}, &logger)

	router := chi.NewRouter()
	router.Route("/v1", func(subRouter chi.Router) {
		subRouter.Get("/health", health.Handler())
		subRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/register", auth.RegistrationHandler(cognito, &logger))
			authRouter.Post("/login", auth.LoginHandler(cognito, &logger))
			authRouter.Post("/confirm", auth.ConfirmRegistrationHandler(cognito, &logger))
		})
	})

	logger.Debug().Msg("Server running locally and listening on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
