package app

import (
	"log"
	"net/http"
	"os"
	"recollection/handlers/authHandler"
	"recollection/handlers/healthHandler"

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

	deps, err := getDependencies(&logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize dependencies")
		return
	}

	router := chi.NewRouter()
	router.Route("/v1", func(subRouter chi.Router) {
		subRouter.Get("/health", healthHandler.Handler())
		subRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/register", authHandler.RegistrationHandler(deps.authClient, deps.userService, &logger))
			authRouter.Post("/login", authHandler.LoginHandler(deps.authClient, &logger))
			authRouter.Post("/confirm", authHandler.ConfirmRegistrationHandler(deps.authClient, &logger))
		})
	})

	PORT := os.Getenv("SERVER_PORT")
	logger.Debug().Msg("Server running locally and listening on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
