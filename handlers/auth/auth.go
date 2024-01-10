package auth

import (
	"encoding/json"
	"net/http"
	"recollection/authClient"
	"recollection/entity"
	"recollection/utils"

	"github.com/rs/zerolog"
)

func RegistrationHandler(auth authClient.Client, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(entity.RegistrationInputBody)
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

func LoginHandler(auth authClient.Client, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(entity.LoginInputBody)
		if err := decoder.Decode(input); err != nil {
			msg := "JSON decode failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		res, err := auth.Login(input)
		if err != nil {
			msg := "auth login failed"
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

func ConfirmRegistrationHandler(auth authClient.Client, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		input := new(entity.RegistrationConfirmationInputBody)
		if err := decoder.Decode(input); err != nil {
			msg := "JSON decode failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		if err := auth.ConfirmRegistration(input); err != nil {
			msg := "confirm registration failed"
			logger.Error().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
