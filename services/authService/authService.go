package authService

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/rs/zerolog"
)

type Config struct {
	AWSRegion   string
	UserPoolID  string
	AppClientID string
}

type Client interface {
	Register(input *RegistrationInputBody) error
	Login(input *LoginInputBody) (*cognito.AuthenticationResultType, error)
	ConfirmRegistration(input *RegistrationConfirmationInputBody) error
}

type clientImpl struct {
	cognito *cognito.CognitoIdentityProvider
	config  *Config
	logger  *zerolog.Logger
}

func New(config *Config, logger *zerolog.Logger) (Client, error) {
	baseConfig := &aws.Config{
		Region: aws.String(config.AWSRegion),
	}
	sess, err := session.NewSession(baseConfig)
	if err != nil {
		logger.Error().Err(err).Msg("AWS session initialization failed")
		return nil, err
	}

	return clientImpl{
		cognito: cognito.New(sess),
		config:  config,
		logger:  logger,
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
