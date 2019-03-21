package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	wordwrap "github.com/mitchellh/go-wordwrap"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/stores/dotenv"
	"go.mozilla.org/sops/stores/ini"
	"go.mozilla.org/sops/stores/json"
	"go.mozilla.org/sops/stores/yaml"
	"go.mozilla.org/sops/version"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
)

type ExampleFileEmitter interface {
	EmitExample() []byte
}

type Store interface {
	sops.Store
	ExampleFileEmitter
}

// DecryptTreeOpts are the options needed to decrypt a tree
type DecryptTreeOpts struct {
	// Tree is the tree to be decrypted
	Tree *sops.Tree
	// KeyServices are the key services to be used for decryption of the data key
	KeyServices []keyservice.KeyServiceClient
	// IgnoreMac is whether or not to ignore the Message Authentication Code included in the SOPS tree
	IgnoreMac bool
	// Cipher is the cryptographic cipher to use to decrypt the values inside the tree
	Cipher sops.Cipher
}

// DecryptTree decrypts the tree passed in through the DecryptTreeOpts and additionally returns the decrypted data key
func DecryptTree(opts DecryptTreeOpts) (dataKey []byte, err error) {
	dataKey, err = opts.Tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return nil, NewExitError(err, codes.CouldNotRetrieveKey)
	}
	computedMac, err := opts.Tree.Decrypt(dataKey, opts.Cipher)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), codes.ErrorDecryptingTree)
	}
	fileMac, err := opts.Cipher.Decrypt(opts.Tree.Metadata.MessageAuthenticationCode, dataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339))
	if !opts.IgnoreMac {
		if fileMac != computedMac {
			// If the file has an empty MAC, display "no MAC" instead of not displaying anything
			if fileMac == "" {
				fileMac = "no MAC"
			}
			return nil, NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", fileMac, computedMac), codes.MacMismatch)
		}
	}
	return dataKey, nil
}

// EncryptTreeOpts are the options needed to encrypt a tree
type EncryptTreeOpts struct {
	// Tree is the tree to be encrypted
	Tree *sops.Tree
	// Cipher is the cryptographic cipher to use to encrypt the values inside the tree
	Cipher sops.Cipher
	// DataKey is the key the cipher should use to encrypt the values inside the tree
	DataKey []byte
}

// EncryptTree encrypts the tree passed in through the EncryptTreeOpts
func EncryptTree(opts EncryptTreeOpts) error {
	unencryptedMac, err := opts.Tree.Encrypt(opts.DataKey, opts.Cipher)
	if err != nil {
		return NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), codes.ErrorEncryptingTree)
	}
	opts.Tree.Metadata.LastModified = time.Now().UTC()
	opts.Tree.Metadata.MessageAuthenticationCode, err = opts.Cipher.Encrypt(unencryptedMac, opts.DataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339))
	if err != nil {
		return NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), codes.ErrorEncryptingMac)
	}
	return nil
}

// LoadEncryptedFile loads an encrypted SOPS file, returning a SOPS tree
func LoadEncryptedFile(loader sops.EncryptedFileLoader, inputPath string) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	path, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, err
	}
	tree, err := loader.LoadEncryptedFile(fileBytes)
	tree.FilePath = path
	return &tree, err
}

func NewExitError(i interface{}, exitCode int) *cli.ExitError {
	if userErr, ok := i.(sops.UserError); ok {
		return NewExitError(userErr.UserError(), exitCode)
	}
	return cli.NewExitError(i, exitCode)
}

func IsYAMLFile(path string) bool {
	return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
}

func IsJSONFile(path string) bool {
	return strings.HasSuffix(path, ".json")
}

func IsEnvFile(path string) bool {
	return strings.HasSuffix(path, ".env")
}

func IsIniFile(path string) bool {
	return strings.HasSuffix(path, ".ini")
}

func DefaultStoreForPath(path string) Store {
	if IsYAMLFile(path) {
		return &yaml.Store{}
	} else if IsJSONFile(path) {
		return &json.Store{}
	} else if IsEnvFile(path) {
		return &dotenv.Store{}
	} else if IsIniFile(path) {
		return &ini.Store{}
	}
	return &json.BinaryStore{}
}

const KMS_ENC_CTX_BUG_FIXED_VERSION = "3.3.0"

func DetectKMSEncryptionContextBug(tree *sops.Tree) (bool, error) {
	versionCheck, err := version.AIsNewerThanB(KMS_ENC_CTX_BUG_FIXED_VERSION, tree.Metadata.Version)
	if err != nil {
		return false, err
	}

	if versionCheck {
		_, _, key := GetKMSKeyWithEncryptionCtx(tree)
		if key != nil {
			return true, nil
		}
	}

	return false, nil
}

func GetKMSKeyWithEncryptionCtx(tree *sops.Tree) (keyGroupIndex int, keyIndex int, key *kms.MasterKey) {
	for i, kg := range tree.Metadata.KeyGroups {
		for n, k := range kg {
			kmsKey, ok := k.(*kms.MasterKey)
			if ok {
				if kmsKey.EncryptionContext != nil && len(kmsKey.EncryptionContext) >= 2 {
					duplicateValues := map[string]int{}
					for _, v := range kmsKey.EncryptionContext {
						duplicateValues[*v] = duplicateValues[*v] + 1
					}
					if len(duplicateValues) > 1 {
						return i, n, kmsKey
					}
				}
			}
		}
	}
	return 0, 0, nil
}

type GenericDecryptOpts struct {
	Cipher      sops.Cipher
	InputStore  sops.Store
	InputPath   string
	IgnoreMAC   bool
	KeyServices []keyservice.KeyServiceClient
}

// LoadEncryptedFileWithBugFixes is a wrapper around LoadEncryptedFile which includes
// check for the issue described in https://github.com/mozilla/sops/pull/435
func LoadEncryptedFileWithBugFixes(opts GenericDecryptOpts) (*sops.Tree, error) {
	tree, err := LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	encCtxBug, err := DetectKMSEncryptionContextBug(tree)
	if err != nil {
		return nil, err
	}
	if encCtxBug {
		tree, err = FixAWSKMSEncryptionContextBug(opts, tree)
		if err != nil {
			return nil, err
		}
	}

	return tree, nil
}

// FixAWSKMSEncryptionContextBug is used to fix the issue described in https://github.com/mozilla/sops/pull/435
func FixAWSKMSEncryptionContextBug(opts GenericDecryptOpts, tree *sops.Tree) (*sops.Tree, error) {
	message := "Up until version 3.3.0 of sops there was a bug surrounding the " +
		"use of encryption context with AWS KMS." +
		"\nYou can read the full description of the issue here:" +
		"\nhttps://github.com/mozilla/sops/pull/435" +
		"\n\nIf a TTY is detected, sops will ask you if you'd like for this issue to be " +
		"automatically fixed, which will require re-encrypting the data keys used by " +
		"each key." +
		"\n\nIf you are not using a TTY, sops will fix the issue for this run.\n\n"
	fmt.Println(wordwrap.WrapString(message, 75))

	persistFix := false

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		var response string
		for response != "y" && response != "n" {
			fmt.Println("Would you like sops to automatically fix this issue? (y/n): ")
			_, err := fmt.Scanln(&response)
			if err != nil {
				return nil, err
			}
		}
		if response == "n" {
			return nil, fmt.Errorf("Exiting. User responded no.")
		} else {
			persistFix = true
		}
	}

	dataKey := []byte{}
	// If there is another key, then we should be able to just decrypt
	// without having to try different variations of the encryption context.
	dataKey, err := DecryptTree(DecryptTreeOpts{
		Cipher:      opts.Cipher,
		IgnoreMac:   opts.IgnoreMAC,
		Tree:        tree,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		dataKey = RecoverDataKeyFromBuggyKMS(opts, tree)
	}

	if dataKey == nil {
		return nil, NewExitError(fmt.Sprintf("Failed to decrypt, meaning there is likely another problem from the encryption context bug: %s", err), codes.ErrorDecryptingTree)
	}

	errs := tree.Metadata.UpdateMasterKeysWithKeyServices(dataKey, opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not re-encrypt data key: %s", errs)
		return nil, err
	}

	err = EncryptTree(EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    tree,
		Cipher:  opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	// If we are not going to persist the fix, just return the re-encrypted tree.
	if !persistFix {
		return tree, nil
	}

	encryptedFile, err := opts.InputStore.EmitEncryptedFile(*tree)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}

	file, err := os.Create(opts.InputPath)
	defer file.Close()
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Could not open file for writing: %s", err), codes.CouldNotWriteOutputFile)
	}
	_, err = file.Write(encryptedFile)
	if err != nil {
		return nil, err
	}

	newTree, err := LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	return newTree, nil
}

// RecoverDataKeyFromBuggyKMS loops through variations on Encryption Context to
// recover the datakey. This is used to fix the issue described in https://github.com/mozilla/sops/pull/435
func RecoverDataKeyFromBuggyKMS(opts GenericDecryptOpts, tree *sops.Tree) []byte {
	kgndx, kndx, originalKey := GetKMSKeyWithEncryptionCtx(tree)

	keyToEdit := *originalKey

	encCtxVals := map[string]interface{}{}
	for _, v := range keyToEdit.EncryptionContext {
		encCtxVals[*v] = ""
	}

	encCtxVariations := []map[string]*string{}
	for ctxVal := range encCtxVals {
		encCtxVariation := map[string]*string{}
		for key := range keyToEdit.EncryptionContext {
			val := ctxVal
			encCtxVariation[key] = &val
		}
		encCtxVariations = append(encCtxVariations, encCtxVariation)
	}

	for _, encCtxVar := range encCtxVariations {
		keyToEdit.EncryptionContext = encCtxVar
		tree.Metadata.KeyGroups[kgndx][kndx] = &keyToEdit
		dataKey, err := DecryptTree(DecryptTreeOpts{
			Cipher:      opts.Cipher,
			IgnoreMac:   opts.IgnoreMAC,
			Tree:        tree,
			KeyServices: opts.KeyServices,
		})
		if err == nil {
			tree.Metadata.KeyGroups[kgndx][kndx] = originalKey
			tree.Metadata.Version = version.Version
			return dataKey
		}
	}

	return nil
}
