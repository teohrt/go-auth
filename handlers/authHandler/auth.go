package authHandler

import (
	"encoding/json"
	"net/http"
	"recollection/services/authService"
	"recollection/services/userService"
	"recollection/utils"

	"github.com/rs/zerolog"
)

func RegistrationHandler(auth authService.Client, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(authService.RegistrationInputBody)
		if err := decoder.Decode(input); err != nil {
			msg := "JSON decode failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		if err := auth.Register(input); err != nil {
			msg := "Auth Register Failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(svc userService.Service, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		decoder := json.NewDecoder(r.Body)

		input := new(authService.LoginInputBody)
		if err := decoder.Decode(input); err != nil {
			msg := "JSON decode failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		res, err := svc.Login(&ctx, input)
		if err != nil {
			msg := "Failed login"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
			return
		}

		jsonObj, err := json.Marshal(res)
		if err != nil {
			msg := "Failed marshalling json"
			logger.Warn().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonObj)
	}
}

func ConfirmRegistrationHandler(svc userService.Service, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		decoder := json.NewDecoder(r.Body)

		input := new(authService.RegistrationConfirmationInputBody)
		if err := decoder.Decode(input); err != nil {
			msg := "JSON decode failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		err := svc.ConfirmUserRegistration(&ctx, input)
		if err != nil {
			msg := "Failed User Registration"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
