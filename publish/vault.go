package publish

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	vault "github.com/hashicorp/vault/api"
	"go.mozilla.org/sops/v3/logging"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("PUBLISH")
}

type VaultDestination struct {
	vaultAddress string
	vaultPath    string
	kvMountName  string
	kvVersion    int
}

func NewVaultDestination(vaultAddress, vaultPath, kvMountName string, kvVersion int) *VaultDestination {
	if !strings.HasSuffix(vaultPath, "/") {
		vaultPath = vaultPath + "/"
	}
	if kvMountName == "" {
		kvMountName = "secret/"
	}
	if !strings.HasSuffix(kvMountName, "/") {
		kvMountName = kvMountName + "/"
	}
	if kvVersion != 1 && kvVersion != 2 {
		kvVersion = 2
	}
	return &VaultDestination{vaultAddress, vaultPath, kvMountName, kvVersion}
}

func (vaultd *VaultDestination) getAddress() string {
	if vaultd.vaultAddress != "" {
		return vaultd.vaultAddress
	}
	return vault.DefaultConfig().Address
}

func (vaultd *VaultDestination) Path(fileName string) string {
	return fmt.Sprintf("%s/v1/%s", vaultd.getAddress(), vaultd.secretsPath(fileName))
}

func (vaultd *VaultDestination) secretsPath(fileName string) string {
	if vaultd.kvVersion == 1 {
		return fmt.Sprintf("%s%s%s", vaultd.kvMountName, vaultd.vaultPath, fileName)
	}
	return fmt.Sprintf("%sdata/%s%s", vaultd.kvMountName, vaultd.vaultPath, fileName)
}

// Returns NotImplementedError
func (vaultd *VaultDestination) Upload(fileContents []byte, fileName string) error {
	return &NotImplementedError{"Vault does not support uploading encrypted sops files directly."}
}

func (vaultd *VaultDestination) UploadUnencrypted(data map[string]interface{}, fileName string) error {
	client, err := vault.NewClient(nil)
	if err != nil {
		return err
	}
	if vaultd.vaultAddress != "" {
		err = client.SetAddress(vaultd.vaultAddress)
		if err != nil {
			return err
		}
	}

	secretsPath := vaultd.secretsPath(fileName)
	existingSecret, err := client.Logical().Read(secretsPath)
	if err != nil {
		log.Warnf("Cannot check if destination secret already exists in %s. New version will be created even if the data has not been changed.", secretsPath)
	}
	if existingSecret != nil && cmp.Equal(data, existingSecret.Data["data"]) {
		log.Infof("Secret in %s is already up-to-date.\n", secretsPath)
		return nil
	}

	secretsData := make(map[string]interface{})

	if vaultd.kvVersion == 1 {
		secretsData = data
	} else if vaultd.kvVersion == 2 {
		secretsData["data"] = data
	}

	_, err = client.Logical().Write(secretsPath, secretsData)
	if err != nil {
		return err
	}

	return nil
}
