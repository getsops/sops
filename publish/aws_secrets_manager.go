package publish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
	"github.com/getsops/sops/v3/logging"
	"github.com/sirupsen/logrus"
)

var awsSecretsLog *logrus.Logger

func init() {
	awsSecretsLog = logging.NewLogger("PUBLISH")
}

// AWSSecretsManagerDestination is the AWS Secrets Manager implementation of the Destination interface
type AWSSecretsManagerDestination struct {
	region     string
	secretName string
}

// NewAWSSecretsManagerDestination is the constructor for an AWS Secrets Manager Destination
func NewAWSSecretsManagerDestination(region, secretName string) *AWSSecretsManagerDestination {
	return &AWSSecretsManagerDestination{region, secretName}
}

// Path returns the AWS Secrets Manager path/ARN of a secret
func (awssmsd *AWSSecretsManagerDestination) Path(fileName string) string {
	if awssmsd.secretName != "" {
		return fmt.Sprintf("arn:aws:secretsmanager:%s:*:secret:%s", awssmsd.region, awssmsd.secretName)
	}
	return fmt.Sprintf("arn:aws:secretsmanager:%s:*:secret:%s", awssmsd.region, fileName)
}

// Returns NotImplementedError
func (awssmsd *AWSSecretsManagerDestination) Upload(fileContents []byte, fileName string) error {
	return &NotImplementedError{"AWS Secrets Manager does not support uploading encrypted sops files directly. Use UploadUnencrypted instead."}
}

// UploadUnencrypted uploads unencrypted data to AWS Secrets Manager as JSON
func (awssmsd *AWSSecretsManagerDestination) UploadUnencrypted(data map[string]interface{}, fileName string) error {
	ctx := context.TODO()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awssmsd.region))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	// Determine secret name - use configured name or derive from filename
	secretName := awssmsd.secretName
	if secretName == "" {
		secretName = fileName
	}

	// Convert data to JSON string for storage
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	secretString := string(jsonData)

	// Check if secret metadata exists first
	_, err = client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretName),
	})

	secretExists := true
	hasValue := false
	var getSecretOutput *secretsmanager.GetSecretValueOutput

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ResourceNotFoundException" {
			secretExists = false
			awsSecretsLog.Infof("Secret %s does not exist, will create new secret", secretName)
		} else {
			awsSecretsLog.Warnf("Cannot check if destination secret already exists in %s. New version will be created even if the data has not been changed.", secretName)
		}
	} else {
		// Secret exists, now check if it has a value
		getSecretOutput, err = client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ResourceNotFoundException" {
				hasValue = false
				awsSecretsLog.Infof("Secret %s exists but has no value, will add initial value", secretName)
			} else {
				awsSecretsLog.Warnf("Cannot retrieve current value of secret %s: %v", secretName, err)
				hasValue = false
			}
		} else {
			hasValue = true
		}
	}

	// If secret exists and has value, check if content is identical
	if secretExists && hasValue && getSecretOutput.SecretString != nil {
		if *getSecretOutput.SecretString == secretString {
			awsSecretsLog.Infof("Secret %s is already up-to-date.", secretName)
			return nil
		}
	}

	// Create or update secret
	if secretExists {
		// Update existing secret
		_, err = client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(secretName),
			SecretString: aws.String(secretString),
		})
		if err != nil {
			return fmt.Errorf("failed to update secret %s: %w", secretName, err)
		}
		awsSecretsLog.Infof("Successfully updated secret %s", secretName)
	} else {
		// Create new secret
		_, err = client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secretName),
			SecretString: aws.String(secretString),
			Description:  aws.String("Secret created by SOPS publish command"),
		})
		if err != nil {
			return fmt.Errorf("failed to create secret %s: %w", secretName, err)
		}
		awsSecretsLog.Infof("Successfully created secret %s", secretName)
	}

	return nil
}
