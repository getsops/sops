package publish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
	"github.com/getsops/sops/v3/logging"
	"github.com/sirupsen/logrus"
)

var parameterLog *logrus.Logger

func init() {
	parameterLog = logging.NewLogger("PUBLISH")
}

// AWSParameterStoreDestination is the AWS Parameter Store implementation of the Destination interface
type AWSParameterStoreDestination struct {
	region        string
	parameterPath string
	parameterType string
}

// NewAWSParameterStoreDestination is the constructor for an AWS Parameter Store Destination
func NewAWSParameterStoreDestination(region, parameterPath, parameterType string) *AWSParameterStoreDestination {
	// Default to SecureString if not specified
	if parameterType == "" {
		parameterType = "SecureString"
	}

	// Ensure parameter path starts with /
	if parameterPath != "" && !strings.HasPrefix(parameterPath, "/") {
		parameterPath = "/" + parameterPath
	}

	return &AWSParameterStoreDestination{region, parameterPath, parameterType}
}

// Path returns the AWS Parameter Store path
func (awspsd *AWSParameterStoreDestination) Path(fileName string) string {
	if awspsd.parameterPath != "" {
		// If path ends with /, append filename; otherwise use path as-is
		if strings.HasSuffix(awspsd.parameterPath, "/") {
			return awspsd.parameterPath + fileName
		}
		return awspsd.parameterPath
	}
	// Default: use filename as parameter name
	if !strings.HasPrefix(fileName, "/") {
		return "/" + fileName
	}
	return fileName
}

// Upload uploads encrypted file contents to AWS Parameter Store
// This stores the entire encrypted file as a parameter value
func (awspsd *AWSParameterStoreDestination) Upload(fileContents []byte, fileName string) error {
	ctx := context.TODO()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awspsd.region))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := ssm.NewFromConfig(cfg)
	parameterName := awspsd.Path(fileName)

	// Check if parameter already exists and compare content
	getParamOutput, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true), // Decrypt for comparison if it's a SecureString
	})

	parameterExists := true
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ParameterNotFound" {
			parameterExists = false
			parameterLog.Infof("Parameter %s does not exist, will create new parameter", parameterName)
		} else {
			parameterLog.Warnf("Cannot check if destination parameter already exists in %s. New version will be created even if the data has not been changed.", parameterName)
		}
	}

	// If parameter exists, check if content is identical
	currentValue := string(fileContents)
	if parameterExists && getParamOutput.Parameter.Value != nil {
		if *getParamOutput.Parameter.Value == currentValue {
			parameterLog.Infof("Parameter %s is already up-to-date.", parameterName)
			return nil
		}
	}

	// Determine parameter type
	var paramType types.ParameterType
	switch strings.ToLower(awspsd.parameterType) {
	case "string":
		paramType = types.ParameterTypeString
	case "stringlist":
		paramType = types.ParameterTypeStringList
	case "securestring":
		paramType = types.ParameterTypeSecureString
	default:
		paramType = types.ParameterTypeSecureString // Default to SecureString for security
	}

	// Put parameter (creates or updates)
	_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:        aws.String(parameterName),
		Value:       aws.String(currentValue),
		Type:        paramType,
		Overwrite:   aws.Bool(true),
		Description: aws.String("Parameter created/updated by SOPS publish command"),
	})

	if err != nil {
		return fmt.Errorf("failed to put parameter %s: %w", parameterName, err)
	}

	if parameterExists {
		parameterLog.Infof("Successfully updated parameter %s", parameterName)
	} else {
		parameterLog.Infof("Successfully created parameter %s", parameterName)
	}

	return nil
}

// UploadUnencrypted uploads unencrypted data to AWS Parameter Store as JSON
func (awspsd *AWSParameterStoreDestination) UploadUnencrypted(data map[string]interface{}, fileName string) error {
	ctx := context.TODO()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awspsd.region))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := ssm.NewFromConfig(cfg)
	parameterName := awspsd.Path(fileName)

	// Convert data to JSON string for storage
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	parameterValue := string(jsonData)

	// Check if parameter already exists and compare content
	getParamOutput, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true), // Decrypt for comparison if it's a SecureString
	})

	parameterExists := true
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ParameterNotFound" {
			parameterExists = false
			parameterLog.Infof("Parameter %s does not exist, will create new parameter", parameterName)
		} else {
			parameterLog.Warnf("Cannot check if destination parameter already exists in %s. New version will be created even if the data has not been changed.", parameterName)
		}
	}

	// If parameter exists, check if content is identical
	if parameterExists && getParamOutput.Parameter.Value != nil {
		if *getParamOutput.Parameter.Value == parameterValue {
			parameterLog.Infof("Parameter %s is already up-to-date.", parameterName)
			return nil
		}
	}

	// Determine parameter type
	var paramType types.ParameterType
	switch strings.ToLower(awspsd.parameterType) {
	case "string":
		paramType = types.ParameterTypeString
	case "stringlist":
		paramType = types.ParameterTypeStringList
	case "securestring":
		paramType = types.ParameterTypeSecureString
	default:
		paramType = types.ParameterTypeSecureString // Default to SecureString for security
	}

	// Put parameter (creates or updates)
	_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:        aws.String(parameterName),
		Value:       aws.String(parameterValue),
		Type:        paramType,
		Overwrite:   aws.Bool(true),
		Description: aws.String("Parameter created/updated by SOPS publish command"),
	})

	if err != nil {
		return fmt.Errorf("failed to put parameter %s: %w", parameterName, err)
	}

	if parameterExists {
		parameterLog.Infof("Successfully updated parameter %s", parameterName)
	} else {
		parameterLog.Infof("Successfully created parameter %s", parameterName)
	}

	return nil
}
