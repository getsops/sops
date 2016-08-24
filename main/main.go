package main

import (
	"go.mozilla.org/sops"

	"crypto/rand"
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
	exitErrorEncryptingTree                 int = 23
	exitErrorDecryptingTree                 int = 23
	exitCannotChangeKeysFromNonExistentFile int = 49
	exitMacMismatch                         int = 51
	exitConfigFileNotFound                  int = 61
	exitKeyboardInterrupt                   int = 85
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

func findKey(keysources []sops.KeySource) ([]byte, error) {
	for _, ks := range keysources {
		for _, k := range ks.Keys {
			key, err := k.Decrypt()
			if err == nil {
				return key, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not get master key")
}

func decrypt(c *cli.Context, file string, fileBytes []byte, output io.Writer) error {
	store := store(file)
	metadata, err := store.LoadMetadata(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	key, err := findKey(metadata.KeySources)
	if err != nil {
		return cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	branch, err := store.Load(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree := sops.Tree{Branch: branch, Metadata: metadata}
	cipher := aes.Cipher{}
	mac, err := tree.Decrypt(key, cipher)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), exitErrorDecryptingTree)
	}
	originalMac, err := cipher.Decrypt(metadata.MessageAuthenticationCode, key, []byte(metadata.LastModified.Format(time.RFC3339)))
	if originalMac != mac && !c.Bool("ignore-mac") {
		return cli.NewExitError("MAC mismatch.", 9)
	}
	out, err := store.Dump(tree.Branch)
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
	branch, err := store.Load(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	var metadata sops.Metadata
	metadata.UnencryptedSuffix = c.String("unencrypted-suffix")
	metadata.Version = "2.0.0"
	var kmsKeys []sops.MasterKey
	if c.String("kms") != "" {
		for _, k := range kms.MasterKeysFromArnString(c.String("kms")) {
			kmsKeys = append(kmsKeys, &k)
		}
	}
	metadata.KeySources = append(metadata.KeySources, sops.KeySource{Name: "kms", Keys: kmsKeys})

	var pgpKeys []sops.MasterKey
	if c.String("pgp") != "" {
		for _, k := range pgp.MasterKeysFromFingerprintString(c.String("pgp")) {
			pgpKeys = append(pgpKeys, &k)
		}
	}
	metadata.KeySources = append(metadata.KeySources, sops.KeySource{Name: "pgp", Keys: pgpKeys})
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not generate random key: %s", err), exitCouldNotRetrieveKey)
	}
	for _, ks := range metadata.KeySources {
		for _, k := range ks.Keys {
			err = k.Encrypt(key)
		}
	}
	tree := sops.Tree{Branch: branch, Metadata: metadata}
	cipher := aes.Cipher{}
	mac, err := tree.Encrypt(key, cipher)
	metadata.MessageAuthenticationCode = mac
	out, err := store.DumpWithMetadata(tree.Branch, metadata)
	_, err = output.Write([]byte(out))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not write to output stream: %s", err), exitCouldNotWriteOutputFile)
	}
	return nil
}

func rotate(c *cli.Context, file string, fileBytes []byte, output io.Writer) error {
	store := store(file)
	metadata, err := store.LoadMetadata(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	key, err := findKey(metadata.KeySources)
	if err != nil {
		return cli.NewExitError(err.Error(), 4)
	}
	branch, err := store.Load(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree := sops.Tree{Branch: branch, Metadata: metadata}
	cipher := aes.Cipher{}
	mac, err := tree.Decrypt(key, cipher)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), 8)
	}
	originalMac, err := cipher.Decrypt(metadata.MessageAuthenticationCode, key, []byte(metadata.LastModified.Format(time.RFC3339)))
	if originalMac != mac && !c.Bool("ignore-mac") {
		return cli.NewExitError("MAC mismatch.", 9)
	}
	newKey := make([]byte, 32)
	_, err = rand.Read(newKey)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not generate random key: %s", err), exitCouldNotRetrieveKey)
	}
	for _, ks := range metadata.KeySources {
		for _, k := range ks.Keys {
			k.Encrypt(newKey)
		}
	}
	_, err = tree.Encrypt(newKey, cipher)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	fmt.Println(metadata.KeySources)
	metadata.AddKMSMasterKeys(c.String("add-kms"))
	metadata.AddPGPMasterKeys(c.String("add-pgp"))
	metadata.RemoveKMSMasterKeys(c.String("rm-kms"))
	metadata.RemovePGPMasterKeys(c.String("rm-pgp"))
	metadata.UpdateMasterKeys(newKey)
	fmt.Println(metadata.KeySources)
	out, err := store.DumpWithMetadata(tree.Branch, metadata)

	_, err = output.Write([]byte(out))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not write to output stream: %s", err), exitCouldNotWriteOutputFile)
	}
	return nil
}
