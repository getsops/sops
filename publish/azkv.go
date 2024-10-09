package publish

import (
	"context"
	"fmt"
	"strings"

	"go.mozilla.org/sops/v3/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/sirupsen/logrus"
)

var azkvLog *logrus.Logger

func init() {
	azkvLog = logging.NewLogger("AZKVPublish")
}

type AzureKeyVaultDestination struct {
	azureKeyVaultURL string
	publishSuffix    string
}

func NewAzureKeyVaultDestination(azureKeyVaultURL string, publishSuffix string) *AzureKeyVaultDestination {
	return &AzureKeyVaultDestination{azureKeyVaultURL, publishSuffix}
}

func (azkvd *AzureKeyVaultDestination) Path(fileName string) string {
	return fmt.Sprintf("%s", azkvd.azureKeyVaultURL)
}

// Returns NotImplementedError
func (azkvd *AzureKeyVaultDestination) Upload(fileContents []byte, fileName string) error {
	return &NotImplementedError{"Azure Key Vault does not support uploading encrypted sops files directly."}
}

func (azkvd *AzureKeyVaultDestination) UploadUnencrypted(data map[string]interface{}, fileName string) error {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to retrieve Azure credential %v", err)
		return err
	}

	client, err := azsecrets.NewClient(azkvd.azureKeyVaultURL, credential, nil)
	if err != nil {
		log.Fatalf("failed to create azsecrets client %v", err)
		return err
	}

	for secretName, secretElement := range data {

		if !strings.HasSuffix(secretName, azkvd.publishSuffix) {
			continue
		}

		secretValue := fmt.Sprintf("%v", secretElement)

		sanitizedSecretName := strings.ReplaceAll(secretName, "_", "-")

		azkvLog.Infof("Uploading Secret Name: %v, Sanitized: %v", secretName, sanitizedSecretName)

		// Create a secret
		params := azsecrets.SetSecretParameters{Value: &secretValue}
		_, err = client.SetSecret(context.TODO(), sanitizedSecretName, params, nil)
		if err != nil {
			log.Fatalf("failed to create a secret %v: %v", secretName, err)
		}
	}

	return nil
}
