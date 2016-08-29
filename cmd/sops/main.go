package main

import (
	"go.mozilla.org/sops"

	"fmt"
	"go.mozilla.org/sops/aes"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"go.mozilla.org/sops/yaml"
	"gopkg.in/urfave/cli.v1"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	exitCouldNotReadInputFile               int = 2
	exitCouldNotWriteOutputFile             int = 3
	exitErrorDumpingTree                    int = 4
	exitErrorReadingConfig                  int = 5
	exitErrorEncryptingTree                 int = 23
	exitErrorDecryptingTree                 int = 23
	exitCannotChangeKeysFromNonExistentFile int = 49
	exitMacMismatch                         int = 51
	exitConfigFileNotFound                  int = 61
	exitKeyboardInterrupt                   int = 85
	exitInvalidTreePathFormat               int = 91
	exitNoFileSpecified                     int = 100
	exitCouldNotRetrieveKey                 int = 128
	exitNoEncryptionKeyFound                int = 111
	exitFileHasNotBeenModified              int = 200
	exitNoEditorFound                       int = 201
)

func main() {
	app := cli.NewApp()
	app.Name = "sops"
	app.Usage = "sops - encrypted file editor with AWS KMS and GPG support"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "decrypt, d",
			Usage: "decrypt a file and output the result to stdout",
		},
		cli.BoolFlag{
			Name:  "encrypt, e",
			Usage: "encrypt a file and output the result to stdout",
		},
		cli.BoolFlag{
			Name:  "rotate, r",
			Usage: "generate a new data encryption key and reencrypt all values with the new key",
		},
		cli.StringFlag{
			Name:   "kms, k",
			Usage:  "comma separated list of KMS ARNs",
			EnvVar: "SOPS_KMS_ARN",
		},
		cli.StringFlag{
			Name:   "pgp, p",
			Usage:  "comma separated list of PGP fingerprints",
			EnvVar: "SOPS_PGP_FP",
		},
		cli.BoolFlag{
			Name:  "in-place, i",
			Usage: "write output back to the same file instead of stdout",
		},
		cli.StringFlag{
			Name:  "extract",
			Usage: "extract a specific key or branch from the input document. Decrypt mode only. Example: --extract '[\"somekey\"][0]'",
		},
		cli.StringFlag{
			Name:  "input-type",
			Usage: "currently json and yaml are supported. If not set, sops will use the file's extension to determine the type",
		},
		cli.StringFlag{
			Name:  "output-type",
			Usage: "currently json and yaml are supported. If not set, sops will use the input file's extension to determine the output format",
		},
		cli.BoolFlag{
			Name:  "show-master-keys, s",
			Usage: "display master encryption keys in the file during editing",
		},
		cli.StringFlag{
			Name:  "add-kms",
			Usage: "add the provided comma-separated list of KMS ARNs to the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "rm-kms",
			Usage: "remove the provided comma-separated list of KMS ARNs from the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "add-pgp",
			Usage: "add the provided comma-separated list of PGP fingerprints to the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "rm-pgp",
			Usage: "remove the provided comma-separated list of PGP fingerprints from the list of master keys on the given file",
		},
		cli.BoolFlag{
			Name:  "ignore-mac",
			Usage: "ignore Message Authentication Code during decryption",
		},
		cli.StringFlag{
			Name:  "unencrypted-suffix",
			Usage: "override the unencrypted key suffix. default: unencrypted_",
			Value: sops.DefaultUnencryptedSuffix,
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "path to sops' config file. If set, sops will not search for the config file recursively.",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.NewExitError("Error: no file specified", exitNoFileSpecified)
		}
		file := c.Args()[0]
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("rm-kms") != "" || c.String("rm-pgp") != "" {
				return cli.NewExitError("Error: cannot add or remove keys on non-existent files, use `--kms` and `--pgp` instead.", 49)
			}
			if c.Bool("encrypt") || c.Bool("decrypt") || c.Bool("rotate") {
				return cli.NewExitError("Error: cannot operate on non-existent file", exitNoFileSpecified)
			}
		}
		fileBytes, err := ioutil.ReadFile(file)

		var output *os.File
		if c.Bool("in-place") {
			var err error
			output, err = os.Create(file)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Could not open in-place file for writing: %s", err), exitCouldNotWriteOutputFile)
			}
			defer output.Close()
		} else {
			output = os.Stdout
		}
		if err != nil {
			return cli.NewExitError("Error: could not read file", exitCouldNotReadInputFile)
		}
		if c.Bool("encrypt") {
			return encrypt(c, file, fileBytes, output)
		} else if c.Bool("decrypt") {
			return decrypt(c, file, fileBytes, output)
		} else if c.Bool("rotate") {
			return rotate(c, file, fileBytes, output)
		}
		return nil
	}
	app.Run(os.Args)
}

func runEditor(path string) {
	editor := os.Getenv("EDITOR")
	var cmd *exec.Cmd
	if editor == "" {
		cmd := exec.Command("which", "vim", "nano")
		out, err := cmd.Output()
		if err != nil {
			panic("Could not find any editors")
		}
		cmd = exec.Command(strings.Split(string(out), "\n")[0], path)
	} else {
		cmd = exec.Command(editor, path)
	}
	cmd.Run()
}

func store(path string) sops.Store {
	if strings.HasSuffix(path, ".yaml") {
		return &yaml.Store{}
	} else if strings.HasSuffix(path, ".json") {
		// return &json.JSONStore{}
	}
	panic("Unknown file type for file " + path)
}

func decryptFile(store sops.Store, fileBytes []byte, ignoreMac bool) (sops.Tree, error) {
	var tree sops.Tree
	metadata, err := store.UnmarshalMetadata(fileBytes)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	key, err := metadata.GetDataKey()
	if err != nil {
		return tree, cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	branch, err := store.Unmarshal(fileBytes)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree = sops.Tree{Branch: branch, Metadata: metadata}
	cipher := aes.Cipher{}
	mac, err := tree.Decrypt(key, cipher)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), exitErrorDecryptingTree)
	}
	originalMac, err := cipher.Decrypt(metadata.MessageAuthenticationCode, key, []byte(metadata.LastModified.Format(time.RFC3339)))
	if originalMac != mac && !ignoreMac {
		return tree, cli.NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", originalMac, mac), 9)
	}
	return tree, nil
}

func decrypt(c *cli.Context, file string, fileBytes []byte, output io.Writer) error {
	store := store(file)
	tree, err := decryptFile(store, fileBytes, c.Bool("ignore-mac"))
	if err != nil {
		return err
	}
	if c.String("extract") != "" {
		v, err := tree.Branch.Truncate(c.String("extract"))
		if err != nil {
			return cli.NewExitError(err.Error(), exitInvalidTreePathFormat)
		}
		if newBranch, ok := v.(sops.TreeBranch); ok {
			tree.Branch = newBranch
		} else {
			bytes, err := sops.ToBytes(v)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error dumping tree: %s", err), exitErrorDumpingTree)
			}
			output.Write(bytes)
			return nil
		}
	}
	out, err := store.Marshal(tree.Branch)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error dumping file: %s", err), exitErrorDumpingTree)
	}
	_, err = output.Write([]byte(out))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not write to output stream: %s", err), exitCouldNotWriteOutputFile)
	}
	return nil
}

func encrypt(c *cli.Context, file string, fileBytes []byte, output io.Writer) error {
	store := store(file)
	branch, err := store.Unmarshal(fileBytes)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	var metadata sops.Metadata
	metadata.UnencryptedSuffix = c.String("unencrypted-suffix")
	metadata.Version = "2.0.0"
	var kmsKeys []sops.MasterKey
	var pgpKeys []sops.MasterKey

	if c.String("kms") != "" {
		for _, k := range kms.MasterKeysFromArnString(c.String("kms")) {
			kmsKeys = append(kmsKeys, &k)
		}
	}
	if c.String("pgp") != "" {
		for _, k := range pgp.MasterKeysFromFingerprintString(c.String("pgp")) {
			pgpKeys = append(pgpKeys, &k)
		}
	}

	if c.String("kms") == "" && c.String("pgp") == "" {
		var confBytes []byte
		if c.String("config") != "" {
			confBytes, err = ioutil.ReadFile(c.String("config"))
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error loading config file: %s", err), exitErrorReadingConfig)
			}
		}
		kmsString, pgpString, err := yaml.MasterKeyStringsForFile(file, confBytes)
		if err == nil {
			for _, k := range pgp.MasterKeysFromFingerprintString(pgpString) {
				pgpKeys = append(pgpKeys, &k)
			}
			for _, k := range kms.MasterKeysFromArnString(kmsString) {
				kmsKeys = append(kmsKeys, &k)
			}
		}
	}
	kmsKs := sops.KeySource{Name: "kms", Keys: kmsKeys}
	pgpKs := sops.KeySource{Name: "pgp", Keys: pgpKeys}
	metadata.KeySources = append(metadata.KeySources, kmsKs)
	metadata.KeySources = append(metadata.KeySources, pgpKs)
	tree := sops.Tree{Branch: branch, Metadata: metadata}
	key, err := tree.GenerateDataKey()
	if err != nil {
		return cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	cipher := aes.Cipher{}
	mac, err := tree.Encrypt(key, cipher)
	encryptedMac, err := cipher.Encrypt(mac, key, []byte(metadata.LastModified.Format(time.RFC3339)))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), exitErrorEncryptingTree)
	}
	metadata.MessageAuthenticationCode = encryptedMac
	out, err := store.MarshalWithMetadata(tree.Branch, metadata)
	_, err = output.Write([]byte(out))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not write to output stream: %s", err), exitCouldNotWriteOutputFile)
	}
	return nil
}

func rotate(c *cli.Context, file string, fileBytes []byte, output io.Writer) error {
	store := store(file)
	tree, err := decryptFile(store, fileBytes, c.Bool("ignore-mac"))
	if err != nil {
		return err
	}
	newKey, err := tree.GenerateDataKey()
	if err != nil {
		return cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	cipher := aes.Cipher{}
	_, err = tree.Encrypt(newKey, cipher)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	tree.Metadata.AddKMSMasterKeys(c.String("add-kms"))
	tree.Metadata.AddPGPMasterKeys(c.String("add-pgp"))
	tree.Metadata.RemoveKMSMasterKeys(c.String("rm-kms"))
	tree.Metadata.RemovePGPMasterKeys(c.String("rm-pgp"))
	tree.Metadata.UpdateMasterKeys(newKey)
	out, err := store.MarshalWithMetadata(tree.Branch, tree.Metadata)

	_, err = output.Write([]byte(out))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not write to output stream: %s", err), exitCouldNotWriteOutputFile)
	}
	return nil
}
