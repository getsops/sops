package publish

import (
	"fmt"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

type VaultDestination struct {
	vaultAddress string
	vaultPath    string
}

func NewVaultDestination(vaultAddress, vaultPath string) *VaultDestination {
	if !strings.HasSuffix(vaultPath, "/") {
		vaultPath = vaultPath + "/"
	}
	return &VaultDestination{vaultAddress, vaultPath}
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
	return fmt.Sprintf("secret/data/%s%s", vaultd.vaultPath, fileName)
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

	secretsData := make(map[string]interface{})
	secretsData["data"] = data

	_, err = client.Logical().Write(vaultd.secretsPath(fileName), secretsData)
	if err != nil {
		return err
	}

	return nil
}
