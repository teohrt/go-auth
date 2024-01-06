package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"recollection/handlers"

	"github.com/go-chi/chi"
)

func Start() {
	PORT := os.Getenv("SERVER_PORT")

	router := chi.NewRouter()
	router.Route("/v1", func(subRouter chi.Router) {
		subRouter.Get("/health", handlers.Health())
	})

	fmt.Println("Server running locally and listening on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
