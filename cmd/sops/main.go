package main //import "go.mozilla.org/sops/cmd/sops"

import (
	"go.mozilla.org/sops"

	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	encodingjson "encoding/json"
	"reflect"

	"go.mozilla.org/sops/aes"
	"go.mozilla.org/sops/json"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"go.mozilla.org/sops/yaml"
	"gopkg.in/urfave/cli.v1"
)

const (
	exitCouldNotReadInputFile                  int = 2
	exitCouldNotWriteOutputFile                int = 3
	exitErrorDumpingTree                       int = 4
	exitErrorReadingConfig                     int = 5
	exitErrorInvalidKMSEncryptionContextFormat int = 6
	exitErrorInvalidSetFormat                  int = 7
	exitErrorEncryptingMac                     int = 21
	exitErrorEncryptingTree                    int = 23
	exitErrorDecryptingMac                     int = 24
	exitErrorDecryptingTree                    int = 25
	exitCannotChangeKeysFromNonExistentFile    int = 49
	exitMacMismatch                            int = 51
	exitMacNotFound                            int = 52
	exitConfigFileNotFound                     int = 61
	exitKeyboardInterrupt                      int = 85
	exitInvalidTreePathFormat                  int = 91
	exitNoFileSpecified                        int = 100
	exitCouldNotRetrieveKey                    int = 128
	exitNoEncryptionKeyFound                   int = 111
	exitFileHasNotBeenModified                 int = 200
	exitNoEditorFound                          int = 201
	exitFailedToCompareVersions                int = 202
)

func loadPlainFile(c *cli.Context, store sops.Store, fileName string, fileBytes []byte) (tree sops.Tree, err error) {
	branch, err := store.Unmarshal(fileBytes)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree.Branch = branch
	ks, err := getKeySources(c, fileName)
	if err != nil {
		return tree, err
	}
	tree.Metadata = sops.Metadata{
		UnencryptedSuffix: c.String("unencrypted-suffix"),
		Version:           version,
		KeySources:        ks,
	}
	tree.GenerateDataKey()
	return
}

func loadEncryptedFile(c *cli.Context, store sops.Store, fileBytes []byte) (tree sops.Tree, err error) {
	metadata, err := store.UnmarshalMetadata(fileBytes)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error loading file metadata: %s", err), exitCouldNotReadInputFile)
	}
	branch, err := store.Unmarshal(fileBytes)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	return sops.Tree{
		Branch:   branch,
		Metadata: metadata,
	}, nil
}

func main() {
	cli.VersionPrinter = printVersion
	app := cli.NewApp()
	app.Name = "sops"
	app.Usage = "sops - encrypted file editor with AWS KMS and GPG support"
	app.ArgsUsage = "sops [options] file"
	app.Version = version
	app.Authors = []cli.Author{
		{Name: "Julien Vehent", Email: "jvehent@mozilla.com"},
		{Name: "Adrian Utrilla", Email: "adrianutrilla@gmail.com"},
	}
	app.UsageText = `sops is an editor of encrypted files that supports AWS KMS and PGP

   To encrypt or decrypt a document with AWS KMS, specify the KMS ARN
   in the -k flag or in the SOPS_KMS_ARN environment variable.
   (you need valid credentials in ~/.aws/credentials or in your env)

   To encrypt or decrypt using PGP, specify the PGP fingerprint in the
   -p flag or in the SOPS_PGP_FP environment variable.

   To use multiple KMS or PGP keys, separate them by commas. For example:
       $ sops -p "10F2...0A, 85D...B3F21" file.yaml

   The -p and -k flags are only used to encrypt new documents. Editing or
   decrypting existing documents can be done with "sops file" or
   "sops -d file" respectively. The KMS and PGP keys listed in the encrypted
   documents are used then. To manage master keys in existing documents, use
   the "add-{kms,pgp}" and "rm-{kms,pgp}" flags.

   To use a different GPG binary than the one in your PATH, set SOPS_GPG_EXEC.

   To select a different editor than the default (vim), set EDITOR.

   For more information, see the README at github.com/mozilla/sops`
	app.EnableBashCompletion = true
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
		cli.StringFlag{
			Name:  "encryption-context",
			Usage: "comma separated list of KMS encryption context key:value pairs",
		},
		cli.StringFlag{
			Name:  "set",
			Usage: `set a specific key or branch in the input JSON or YAML document. value must be a json encoded string. (edit mode only). eg. --set '["somekey"][0] {"somevalue":true}'`,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.NewExitError("Error: no file specified", exitNoFileSpecified)
		}
		fileName := c.Args()[0]
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("rm-kms") != "" || c.String("rm-pgp") != "" {
				return cli.NewExitError("Error: cannot add or remove keys on non-existent files, use `--kms` and `--pgp` instead.", 49)
			}
			if c.Bool("encrypt") || c.Bool("decrypt") || c.Bool("rotate") {
				return cli.NewExitError("Error: cannot operate on non-existent file", exitNoFileSpecified)
			}
		}

		inputStore := inputStore(c, fileName)
		outputStore := outputStore(c, fileName)

		// Instead of checking for errors after each operation here, we
		// just continue and let each command decide which errors can
		// be handled and which can't.
		fileBytes, readFileErr := ioutil.ReadFile(fileName)
		plainTree, loadPlainFileErr := loadPlainFile(c, inputStore, fileName, fileBytes)
		encryptedTree, loadEncryptedFileErr := loadEncryptedFile(c, inputStore, fileBytes)

		var output []byte
		var err error
		if c.Bool("encrypt") {
			if readFileErr != nil {
				return readFileErr
			}
			if loadPlainFileErr != nil {
				return loadPlainFileErr
			}
			output, err = encrypt(c, plainTree, outputStore)
		}

		if c.Bool("decrypt") {
			if readFileErr != nil {
				return readFileErr
			}
			if loadEncryptedFileErr != nil {
				return loadEncryptedFileErr
			}
			output, err = decrypt(c, encryptedTree, outputStore)
		}
		if c.Bool("rotate") {
			if readFileErr != nil {
				return readFileErr
			}
			if loadEncryptedFileErr != nil {
				return loadEncryptedFileErr
			}
			output, err = rotate(c, encryptedTree, outputStore)
		}

		isEditMode := !c.Bool("encrypt") && !c.Bool("decrypt") && !c.Bool("rotate")
		if isEditMode {
			output, err = edit(c, fileName, fileBytes)
		}

		if err != nil {
			return err
		}
		// We open the file *after* the operations on the tree have been
		// executed to avoid truncating it when there's errors
		var outputFile *os.File
		if c.Bool("in-place") || isEditMode {
			var err error
			outputFile, err = os.Create(fileName)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Could not open in-place file for writing: %s", err), exitCouldNotWriteOutputFile)
			}
			defer outputFile.Close()
		} else {
			outputFile = os.Stdout
		}
		outputFile.Write(output)
		return nil
	}
	app.Run(os.Args)
}

func runEditor(path string) error {
	editor := os.Getenv("EDITOR")
	var cmd *exec.Cmd
	if editor == "" {
		cmd = exec.Command("which", "vim", "nano")
		out, err := cmd.Output()
		if err != nil {
			panic("Could not find any editors")
		}
		cmd = exec.Command(strings.Split(string(out), "\n")[0], path)
	} else {
		cmd = exec.Command(editor, path)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func inputStore(context *cli.Context, path string) sops.Store {
	switch context.String("input-type") {
	case "yaml":
		return &yaml.Store{}
	case "json":
		return &json.Store{}
	default:
		return defaultStore(path)
	}
}
func outputStore(context *cli.Context, path string) sops.Store {
	switch context.String("output-type") {
	case "yaml":
		return &yaml.Store{}
	case "json":
		return &json.Store{}
	default:
		return defaultStore(path)
	}
}

func defaultStore(path string) sops.Store {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return &yaml.Store{}
	} else if strings.HasSuffix(path, ".json") {
		return &json.Store{}
	}
	return &json.BinaryStore{}
}

func decryptTree(tree sops.Tree, ignoreMac bool) (sops.Tree, map[string][]interface{}, error) {
	cipher := aes.Cipher{}
	stash := make(map[string][]interface{})
	key, err := tree.Metadata.GetDataKey()
	if err != nil {
		return tree, nil, cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	computedMac, err := tree.Decrypt(key, cipher, stash)
	if err != nil {
		return tree, nil, cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), exitErrorDecryptingTree)
	}
	fileMac, _, err := cipher.Decrypt(tree.Metadata.MessageAuthenticationCode, key, tree.Metadata.LastModified.Format(time.RFC3339))
	if fileMac != computedMac && !ignoreMac {
		outputMac := fileMac
		if outputMac == "" {
			outputMac = "no MAC"
		}
		return tree, nil, cli.NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", outputMac, computedMac), exitMacMismatch)
	}
	return tree, stash, nil
}

func decrypt(c *cli.Context, tree sops.Tree, outputStore sops.Store) ([]byte, error) {
	tree, _, err := decryptTree(tree, c.Bool("ignore-mac"))
	if err != nil {
		return nil, err
	}
	if c.String("extract") != "" {
		v, err := tree.Branch.Truncate(c.String("extract"))
		if err != nil {
			return nil, cli.NewExitError(err.Error(), exitInvalidTreePathFormat)
		}
		if newBranch, ok := v.(sops.TreeBranch); ok {
			tree.Branch = newBranch
		} else {
			bytes, err := outputStore.MarshalValue(v)
			if err != nil {
				return nil, cli.NewExitError(fmt.Sprintf("Error dumping tree: %s", err), exitErrorDumpingTree)
			}
			return bytes, nil
		}
	}
	out, err := outputStore.Marshal(tree.Branch)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error dumping file: %s", err), exitErrorDumpingTree)
	}
	return out, nil
}

func getKeySources(c *cli.Context, file string) ([]sops.KeySource, error) {
	var kmsKeys []sops.MasterKey
	var pgpKeys []sops.MasterKey
	kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
	if c.String("encryption-context") != "" && kmsEncryptionContext == nil {
		return nil, cli.NewExitError("Invalid KMS encryption context format", exitErrorInvalidKMSEncryptionContextFormat)
	}
	if c.String("kms") != "" {
		for _, k := range kms.MasterKeysFromArnString(c.String("kms"), kmsEncryptionContext) {
			kmsKeys = append(kmsKeys, k)
		}
	}
	if c.String("pgp") != "" {
		for _, k := range pgp.MasterKeysFromFingerprintString(c.String("pgp")) {
			pgpKeys = append(pgpKeys, k)
		}
	}
	var err error
	if c.String("kms") == "" && c.String("pgp") == "" {
		var confBytes []byte
		if c.String("config") != "" {
			confBytes, err = ioutil.ReadFile(c.String("config"))
			if err != nil {
				return nil, cli.NewExitError(fmt.Sprintf("Error loading config file: %s", err), exitErrorReadingConfig)
			}
		}
		kmsString, pgpString, err := yaml.MasterKeyStringsForFile(file, confBytes)
		if err == nil {
			for _, k := range pgp.MasterKeysFromFingerprintString(pgpString) {
				pgpKeys = append(pgpKeys, k)
			}
			for _, k := range kms.MasterKeysFromArnString(kmsString, kmsEncryptionContext) {
				kmsKeys = append(kmsKeys, k)
			}
		}
	}
	kmsKs := sops.KeySource{Name: "kms", Keys: kmsKeys}
	pgpKs := sops.KeySource{Name: "pgp", Keys: pgpKeys}
	return []sops.KeySource{kmsKs, pgpKs}, nil
}

func encryptTree(tree sops.Tree, stash map[string][]interface{}) (sops.Tree, error) {
	cipher := aes.Cipher{}
	key, err := tree.Metadata.GetDataKey()
	if err != nil {
		return tree, cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	computedMac, err := tree.Encrypt(key, cipher, stash)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	tree.Metadata.LastModified = time.Now().UTC()
	encryptedMac, err := cipher.Encrypt(computedMac, key, tree.Metadata.LastModified.Format(time.RFC3339), nil)
	if err != nil {
		return tree, cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), exitErrorEncryptingMac)
	}
	tree.Metadata.MessageAuthenticationCode = encryptedMac
	return tree, nil
}

func encrypt(c *cli.Context, tree sops.Tree, outputStore sops.Store) ([]byte, error) {
	tree, err := encryptTree(tree, nil)
	if err != nil {
		return nil, err
	}
	out, err := outputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return out, err
}

func rotate(c *cli.Context, tree sops.Tree, outputStore sops.Store) ([]byte, error) {
	tree, _, err := decryptTree(tree, c.Bool("ignore-mac"))
	if err != nil {
		return nil, err
	}
	kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
	if c.String("encryption-context") != "" && kmsEncryptionContext == nil {
		return nil, cli.NewExitError("Invalid KMS encryption context format", exitErrorInvalidKMSEncryptionContextFormat)
	}
	tree.Metadata.AddKMSMasterKeys(c.String("add-kms"), kmsEncryptionContext)
	tree.Metadata.AddPGPMasterKeys(c.String("add-pgp"))
	tree.Metadata.RemoveKMSMasterKeys(c.String("rm-kms"))
	tree.Metadata.RemovePGPMasterKeys(c.String("rm-pgp"))
	_, errs := tree.GenerateDataKey()
	if len(errs) > 0 {
		return nil, cli.NewExitError(fmt.Sprintf("Error encrypting the data key with one or more master keys: %s", errs), exitCouldNotRetrieveKey)
	}
	tree, err = encryptTree(tree, nil)
	if err != nil {
		return nil, err
	}
	out, err := outputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return out, nil
}

func hashFile(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}

var exampleTree = sops.TreeBranch{
	sops.TreeItem{
		Key:   "hello",
		Value: `Welcome to SOPS! Edit this file as you please!`,
	},
	sops.TreeItem{
		Key:   "example_key",
		Value: "example_value",
	},
	sops.TreeItem{
		Key: "example_array",
		Value: []interface{}{
			"example_value1",
			"example_value2",
		},
	},
	sops.TreeItem{
		Key:   "example_number",
		Value: 1234.56789,
	},
	sops.TreeItem{
		Key:   "example_booleans",
		Value: []interface{}{true, false},
	},
}

func loadExample(c *cli.Context, file string) (sops.Tree, error) {
	var in []byte
	var tree sops.Tree
	fileStore := inputStore(c, file)
	if _, ok := fileStore.(*json.BinaryStore); ok {
		// Get the value under the first key
		in = []byte(exampleTree[0].Value.(string))
	} else {
		var err error
		in, err = fileStore.Marshal(exampleTree)
		if err != nil {
			return tree, err
		}
	}
	branch, _ := fileStore.Unmarshal(in)
	tree.Branch = branch
	ks, err := getKeySources(c, file)
	if err != nil {
		return tree, err
	}
	tree.Metadata.UnencryptedSuffix = c.String("unencrypted-suffix")
	tree.Metadata.Version = version
	tree.Metadata.KeySources = ks
	key, errs := tree.GenerateDataKey()
	if len(errs) > 0 {
		return tree, cli.NewExitError(fmt.Sprintf("Error encrypting the data key with one or more master keys: %s", errs), exitCouldNotRetrieveKey)
	}
	tree.Metadata.UpdateMasterKeys(key)
	return tree, nil
}

func jsonValueToTreeInsertableValue(jsonValue string) (interface{}, error) {
	var valueToInsert interface{}
	err := encodingjson.Unmarshal([]byte(jsonValue), &valueToInsert)
	if err != nil {
		return nil, cli.NewExitError("Value for --set is not valid JSON", exitErrorInvalidSetFormat)
	}
	// Check if decoding it as json we find a single value
	// and not a map or slice, in which case we can't marshal
	// it to a sops.TreeBranch
	kind := reflect.ValueOf(valueToInsert).Kind()
	if kind == reflect.Map || kind == reflect.Slice {
		var err error
		valueToInsert, err = (&json.Store{}).Unmarshal([]byte(jsonValue))
		if err != nil {
			return nil, cli.NewExitError("Invalid --set value format", exitErrorInvalidSetFormat)
		}
	}
	return valueToInsert, nil
}

func extractSetArguments(set string) (path string, key string, valueToInsert interface{}, err error) {
	// Set is a string with the format "python-dict-index json-value"
	// Since python-dict-index has to end with ], we split at "] " to get the two parts
	pathValuePair := strings.SplitAfterN(set, "] ", 2)
	if len(pathValuePair) < 2 {
		return "", "", nil, cli.NewExitError("Invalid --set format", exitErrorInvalidSetFormat)
	}
	fullPath := strings.TrimRight(pathValuePair[0], " ")
	jsonValue := pathValuePair[1]
	valueToInsert, err = jsonValueToTreeInsertableValue(jsonValue)
	splitPath := strings.Split(fullPath, "[")

	// The path is the full path except the last entry
	path = strings.Join(splitPath[0:len(splitPath)-1], "")

	// The key is the last entry in the full path
	key, err = sops.TrimTreePathComponent(splitPath[len(splitPath)-1])
	if err != nil {
		return "", "", nil, cli.NewExitError("Invalid --set path format", exitErrorInvalidSetFormat)
	}
	return path, key, valueToInsert, nil
}

func edit(c *cli.Context, file string, fileBytes []byte) ([]byte, error) {
	var tree sops.Tree
	var stash map[string][]interface{}
	var err error
	if fileBytes == nil {
		tree, err = loadExample(c, file)
	} else {
		tree, err = loadEncryptedFile(c, inputStore(c, file), fileBytes)
		if err != nil {
			return nil, err
		}
		tree, stash, err = decryptTree(tree, c.Bool("ignore-mac"))
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not load file: %s", err), exitCouldNotReadInputFile)
	}
	if c.String("set") != "" {
		path, key, value, err := extractSetArguments(c.String("set"))
		if err != nil {
			return nil, err
		}
		parent, err := tree.Branch.Truncate(path)
		if err != nil {
			return nil, cli.NewExitError("Could not truncate tree to the provided path", exitErrorInvalidSetFormat)
		}
		branch := parent.(sops.TreeBranch)
		tree.Branch = branch.InsertOrReplaceValue(key, value)
	} else {
		tmpdir, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Could not create temporary directory: %s", err), exitCouldNotWriteOutputFile)
		}
		defer os.RemoveAll(tmpdir)
		tmpfile, err := os.Create(path.Join(tmpdir, path.Base(file)))
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Could not create temporary file: %s", err), exitCouldNotWriteOutputFile)
		}
		var out []byte
		if c.Bool("show-master-keys") {
			out, err = outputStore(c, file).MarshalWithMetadata(tree.Branch, tree.Metadata)
		} else {
			out, err = outputStore(c, file).Marshal(tree.Branch)
		}
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
		}
		_, err = tmpfile.Write(out)
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Could not write output file: %s", err), exitCouldNotWriteOutputFile)
		}
		origHash, err := hashFile(tmpfile.Name())
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Could not hash file: %s", err), exitCouldNotReadInputFile)
		}
		for {
			err = runEditor(tmpfile.Name())
			if err != nil {
				return nil, cli.NewExitError(fmt.Sprintf("Could not run editor: %s", err), exitNoEditorFound)
			}
			newHash, err := hashFile(tmpfile.Name())
			if err != nil {
				return nil, cli.NewExitError(fmt.Sprintf("Could not hash file: %s", err), exitCouldNotReadInputFile)
			}
			if bytes.Equal(newHash, origHash) {
				return nil, cli.NewExitError("File has not changed, exiting.", exitFileHasNotBeenModified)
			}
			edited, err := ioutil.ReadFile(tmpfile.Name())
			if err != nil {
				return nil, cli.NewExitError(fmt.Sprintf("Could not read edited file: %s", err), exitCouldNotReadInputFile)
			}
			newBranch, err := inputStore(c, file).Unmarshal(edited)
			if err != nil {
				fmt.Printf("Could not load tree: %s\nProbably invalid syntax. Press a key to return to the editor, or Ctrl+C to exit.", err)
				bufio.NewReader(os.Stdin).ReadByte()
				continue
			}
			if c.Bool("show-master-keys") {
				metadata, err := inputStore(c, file).UnmarshalMetadata(edited)
				if err != nil {
					fmt.Printf("sops branch is invalid: %s.\nPress a key to return to the editor, or Ctrl+C to exit.", err)
					bufio.NewReader(os.Stdin).ReadByte()
					continue
				}
				tree.Metadata = metadata
			}
			tree.Branch = newBranch
			needVersionUpdated, err := AIsNewerThanB(version, tree.Metadata.Version)
			if err != nil {
				return nil, cli.NewExitError(fmt.Sprintf("Failed to compare document version %q with program version %q: %v", tree.Metadata.Version, version, err), exitFailedToCompareVersions)
			}
			if needVersionUpdated {
				tree.Metadata.Version = version
			}
			if tree.Metadata.MasterKeyCount() == 0 {
				fmt.Println("No master keys were provided, so sops can't encrypt the file.\nPress a key to return to the editor, or Ctrl+C to exit.")
				bufio.NewReader(os.Stdin).ReadByte()
				continue
			}
			break
		}
	}
	tree, err = encryptTree(tree, stash)
	if err != nil {
		return nil, err
	}
	out, err := outputStore(c, file).MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return out, nil
}
