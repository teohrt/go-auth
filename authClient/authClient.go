package authClient

import (
	"fmt"
	"recollection/entity"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type Config struct {
	AWSRegion       string
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
}

type Client interface {
	Register(input *entity.RegistrationInputBody)
	Login(input *entity.LoginInputBody) *string
	ConfirmRegistration(input *entity.RegistrationConfirmationInputBody)
}

type clientImpl struct {
	cognito *cognito.CognitoIdentityProvider
	config  *Config
}

func New(config *Config) clientImpl {
	baseConfig := &aws.Config{
		Region: aws.String(config.AWSRegion),
	}
	sess, err := session.NewSession(baseConfig)
	if err != nil {
		fmt.Println("session broke")
		panic(err)
	}

	return clientImpl{
		cognito: cognito.New(sess),
		config:  config,
	}
}

func (client clientImpl) Register(input *entity.RegistrationInputBody) {
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
		fmt.Println("sign up broke")
		panic(err)
	}
}

func (client clientImpl) Login(input *entity.LoginInputBody) *string {
	flow := aws.String("USER_PASSWORD_AUTH") // magic aws string designating login type
	params := map[string]*string{
		"USERNAME": aws.String(input.Username),
		"PASSWORD": aws.String(input.Password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(client.config.AppClientID),
	}

	res, err := client.cognito.InitiateAuth(authTry)
	if err != nil {
		fmt.Println("InitiateAuth broke")
		panic(err)
	}

	jwt := res.AuthenticationResult.AccessToken
	fmt.Println(*jwt)
	return jwt
}

func (client clientImpl) ConfirmRegistration(input *entity.RegistrationConfirmationInputBody) {
	user := &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(input.Code),
		Username:         aws.String(input.Username),
		ClientId:         aws.String(client.config.AppClientID),
	}

	_, err := client.cognito.ConfirmSignUp(user)
	if err != nil {
		fmt.Println("ConfirmSignUp broke")
		panic(err)
	}
}
