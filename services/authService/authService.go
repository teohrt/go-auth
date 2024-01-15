package authService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rs/zerolog"
)

type Config struct {
	AWSRegion   string `env:"AWS_REGION"`
	UserPoolID  string `env:"AWS_COGNITO_USER_POOL_ID"`
	AppClientID string `env:"AWS_COGNITO_APP_CLIENT_ID"`
}

type Client interface {
	Register(input *RegistrationInputBody) error
	Login(input *LoginInputBody) (*cognito.AuthenticationResultType, error)
	ConfirmRegistration(input *RegistrationConfirmationInputBody) error
	ParseJWTClaims(tokenString string) (*Claims, error)
}

type clientImpl struct {
	cognito      *cognito.CognitoIdentityProvider
	publicKeySet jwk.Set
	config       *Config
	logger       *zerolog.Logger
}

func New(ctx *context.Context, config *Config, logger *zerolog.Logger) (Client, error) {
	baseConfig := &aws.Config{
		Region: aws.String(config.AWSRegion),
	}
	sess, err := session.NewSession(baseConfig)
	if err != nil {
		logger.Error().Err(err).Msg("AWS session initialization failed")
		return nil, err
	}

	publicKeysURL := "https://cognito-idp." + config.AWSRegion + ".amazonaws.com/" + config.UserPoolID + "/.well-known/jwks.json"
	publicKeySet, err := jwk.Fetch(*ctx, publicKeysURL)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve Cognito public key set")
		return nil, err
	}

	return clientImpl{
		cognito:      cognito.New(sess),
		publicKeySet: publicKeySet,
		config:       config,
		logger:       logger,
	}, nil
}

func (client clientImpl) Register(input *RegistrationInputBody) error {
	user := &cognito.SignUpInput{
		Username: aws.String(input.Username),
		Password: aws.String(input.Password),
		ClientId: aws.String(client.config.AppClientID),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Email),
			},
		},
	}

	_, err := client.cognito.SignUp(user)
	if err != nil {
		client.logger.Error().Err(err).Msg("AWS Cognito SignUp failed")
		return err
	}
	return nil
}

func (client clientImpl) Login(input *LoginInputBody) (*cognito.AuthenticationResultType, error) {
	params := map[string]*string{
		"USERNAME": aws.String(input.Username),
		"PASSWORD": aws.String(input.Password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       aws.String("USER_PASSWORD_AUTH"), // magic aws string designating login type
		AuthParameters: params,
		ClientId:       aws.String(client.config.AppClientID),
	}

	res, err := client.cognito.InitiateAuth(authTry)
	if err != nil {
		client.logger.Error().Err(err).Msg("AWS Cognito InitiateAuth failed")
		return nil, err
	}

	return res.AuthenticationResult, nil
}

func (client clientImpl) ConfirmRegistration(input *RegistrationConfirmationInputBody) error {
	user := &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(input.Code),
		Username:         aws.String(input.Username),
		ClientId:         aws.String(client.config.AppClientID),
	}

	_, err := client.cognito.ConfirmSignUp(user)
	if err != nil {
		client.logger.Error().Err(err).Msg("AWS Cognito ConfirmSignUp failed")
		return err
	}
	return nil
}

func (client clientImpl) ParseJWTClaims(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		keys, ok := client.publicKeySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key with specified kid is not present in jwks")
		}
		var publickey interface{}
		err := keys.Raw(&publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}
		return publickey, nil
	})
	if err != nil {
		client.logger.Error().Err(err).Msg("failed parsing JWT")
		return nil, err
	}

	// convert map to json
	jsonString, err := json.Marshal(token.Claims.(jwt.MapClaims))
	if err != nil {
		client.logger.Error().Err(err).Msg("failed parsing marshalling json")
		return nil, err
	}

	claims := Claims{}
	json.Unmarshal(jsonString, &claims)

	// TODO - verify client id, expiration

	return &claims, nil
}
