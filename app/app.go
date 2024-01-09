package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	authClient "recollection/authClient"
	"recollection/handlers/auth"
	"recollection/handlers/health"

	"github.com/go-chi/chi"
)

func Start() {
	PORT := os.Getenv("SERVER_PORT")
	if PORT == "" {
		PORT = "80"
	}

	cognito := authClient.New(&authClient.Config{
		AWSRegion:   os.Getenv("AWS_REGION"),
		UserPoolID:  os.Getenv("AWS_COGNITO_USER_POOL_ID"),
		AppClientID: os.Getenv("AWS_COGNITO_APP_CLIENT_ID"),
	})

	router := chi.NewRouter()
	router.Route("/v1", func(subRouter chi.Router) {
		subRouter.Get("/health", health.Handler())
		subRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/register", auth.RegistrationHandler(cognito))
			authRouter.Post("/login", auth.LoginHandler(cognito))
			authRouter.Post("/confirm", auth.ConfirmRegistrationHandler(cognito))
		})
	})

	fmt.Println("Server running locally and listening on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
