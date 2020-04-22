package qorexGoHelpers

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"net/http"
)

const (
	BadCredentialRequest = 400
	BadCredentialRequestMessage = "Malformed credential request"
	BadCredentialRequestError = "SecretManager error, both SecretName and Region must be provided"

	ErrorUnknown = "An unknown error occurred"
)



// SecretsConfig is used to set up the connection
type SecretsConfig struct {
	SecretName string // Instance name to pass in
	Region string // Region of Secret to pass in
	Response *Response
}

type DBConnectionConfig struct {
	User string `json:"username"`
	Password string `json:"password"`
	Host string `json:"host"`
}

func (s *SecretsConfig) GetSecret() DBConnectionConfig {

	if s.Region == "" || s.SecretName == "" {
		err := errors.New(BadCredentialRequestError)
		s.Response.SetStatus(BadCredentialRequest, BadCredentialRequestMessage, err)
	}
	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(s.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(s.SecretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			err := errors.New(aerr.Error())
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				s.Response.SetStatus(
					http.StatusInternalServerError,
					secretsmanager.ErrCodeDecryptionFailure,
					err)

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				s.Response.SetStatus(
					http.StatusInternalServerError,
					secretsmanager.ErrCodeInternalServiceError,
					err)

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				s.Response.SetStatus(
					http.StatusInternalServerError,
					secretsmanager.ErrCodeInvalidParameterException,
					err)

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				s.Response.SetStatus(
					http.StatusInternalServerError,
					secretsmanager.ErrCodeInvalidRequestException,
					err)

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				s.Response.SetStatus(
					http.StatusInternalServerError,
					secretsmanager.ErrCodeResourceNotFoundException,
					err)

			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			s.Response.SetStatus(
				http.StatusInternalServerError,
				ErrorUnknown,
				err)
		}
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			s.Response.SetStatus(
				http.StatusInternalServerError,
				"Base64 Decode Error",
				err)
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
	}

	_ = decodedBinarySecret

	dbConn := DBConnectionConfig{}
	err = json.Unmarshal([]byte(secretString), &dbConn)
	if err != nil {
		s.Response.SetStatus(
			http.StatusInternalServerError,
			"Error",
			err)
	}

	return dbConn
}