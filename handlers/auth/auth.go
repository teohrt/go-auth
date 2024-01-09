package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recollection/authClient"
	"recollection/entity"
)

func RegistrationHandler(auth authClient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(entity.RegistrationInputBody)
		if err := decoder.Decode(input); err != nil {
			fmt.Println("json broke")
			panic(err)
		}

		auth.Register(input)
		w.Write([]byte("😎"))
	}
}

func LoginHandler(auth authClient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(entity.LoginInputBody)
		if err := decoder.Decode(input); err != nil {
			fmt.Println("json broke")
			panic(err)
		}

		jwt := auth.Login(input)
		w.Write([]byte(*jwt))
	}
}

func ConfirmRegistrationHandler(auth authClient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(entity.RegistrationConfirmationInputBody)
		if err := decoder.Decode(input); err != nil {
			fmt.Println("json broke")
			panic(err)
		}

		auth.ConfirmRegistration(input)
		w.Write([]byte("😎"))
	}
}
