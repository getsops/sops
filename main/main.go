package main

import (
	"go.mozilla.org/sops"

	"crypto/rand"
	"fmt"
	"go.mozilla.org/sops/json"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"go.mozilla.org/sops/yaml"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
			return cli.NewExitError("Error: no file specified", 1)
		}

		file := c.Args()[0]
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("rm-kms") != "" || c.String("rm-pgp") != "" {
				return cli.NewExitError("Error: cannot add or remove keys on non-existent files, use `--kms` and `--pgp` instead.", 49)
			}
			if c.Bool("encrypt") || c.Bool("decrypt") || c.Bool("rotate") {
				return cli.NewExitError("Error: cannot operate on non-existent file", 100)
			}
		}
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return cli.NewExitError("Error: could not read file", 2)
		}
		if c.Bool("encrypt") {
			return encrypt(c, file, fileBytes)
		} else if c.Bool("decrypt") {
			return decrypt(c, file, fileBytes)
		} else if c.Bool("rotate") {

		} else {

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
		return &yaml.YAMLStore{}
	} else if strings.HasSuffix(path, ".json") {
		return &json.JSONStore{}
	}
	panic("Unknown file type for file " + path)
}

func findKey(keysources []sops.KeySource) (string, error) {
	for _, ks := range keysources {
		fmt.Println("Trying keysource: ", ks.Name)
		for _, k := range ks.Keys {
			fmt.Println("Trying key: ", k.ToString())
			key, err := k.Decrypt()
			if err == nil {
				return key, nil
			}
		}
	}
	return "", fmt.Errorf("Could not get master key")
}

func decrypt(c *cli.Context, file string, fileBytes []byte) error {
	store := store(file)
	err := store.LoadMetadata(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), 3)
	}
	key, err := findKey(store.Metadata().KeySources)
	if err != nil {
		return cli.NewExitError(err.Error(), 4)
	}
	err = store.Load(string(fileBytes), key)
	fmt.Println(err == sops.MacMismatch)
	if err == sops.MacMismatch && !c.Bool("ignore-mac") {
		return cli.NewExitError("MAC mismatch", 5)
	} else if err != sops.MacMismatch {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), 6)
	}
	s, err := store.DumpUnencrypted()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error dumping file: %s", err), 7)
	}
	fmt.Print(s)
	return nil
}

func encrypt(c *cli.Context, file string, fileBytes []byte) error {
	store := store(file)
	err := store.LoadUnencrypted(string(fileBytes))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), 4)
	}
	var metadata sops.Metadata
	metadata.UnencryptedSuffix = c.String("unencrypted-suffix")
	metadata.Version = "2.0.0"
	var kmsKeys []sops.MasterKey
	if c.String("kms") != "" {
		for _, k := range kms.KMSMasterKeysFromArnString(c.String("kms")) {
			kmsKeys = append(kmsKeys, &k)
		}
	}
	metadata.KeySources = append(metadata.KeySources, sops.KeySource{Name: "kms", Keys: kmsKeys})

	var pgpKeys []sops.MasterKey
	if c.String("pgp") != "" {
		for _, k := range pgp.GPGMasterKeysFromFingerprintString(c.String("pgp")) {
			pgpKeys = append(pgpKeys, &k)
		}
	}
	metadata.KeySources = append(metadata.KeySources, sops.KeySource{Name: "pgp", Keys: pgpKeys})
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not generate random key: %s", err), 8)
	}
	for _, ks := range metadata.KeySources {
		for _, k := range ks.Keys {
			err = k.Encrypt(string(key))
		}
	}

	store.SetMetadata(metadata)
	out, err := store.Dump(string(key))
	fmt.Println(out)
	return nil
}
