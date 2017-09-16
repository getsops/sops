package main //import "go.mozilla.org/sops/cmd/sops"

import (
	"log"
	"net"
	"net/url"

	"google.golang.org/grpc"

	"go.mozilla.org/sops"

	"fmt"
	"os"
	"strings"
	"time"

	encodingjson "encoding/json"
	"reflect"

	"strconv"

	"go.mozilla.org/sops/aes"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/subcommand/groups"
	keyservicecmd "go.mozilla.org/sops/cmd/sops/subcommand/keyservice"
	"go.mozilla.org/sops/config"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"go.mozilla.org/sops/stores/json"
	yamlstores "go.mozilla.org/sops/stores/yaml"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	cli.VersionPrinter = printVersion
	app := cli.NewApp()

	keyserviceFlags := []cli.Flag{
		cli.BoolTFlag{
			Name:  "enable-local-keyservice",
			Usage: "use local key service",
		},
		cli.StringSliceFlag{
			Name:  "keyservice",
			Usage: "Specify the key services to use in addition to the local one. Can be specified more than once. Syntax: protocol://address. Example: tcp://myserver.com:5000",
		},
	}
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
	app.Commands = []cli.Command{
		{
			Name:  "keyservice",
			Usage: "start a SOPS key service server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "network, net",
					Usage: "network to listen on, e.g. 'tcp' or 'unix'",
					Value: "tcp",
				},
				cli.StringFlag{
					Name:  "address, addr",
					Usage: "address to listen on, e.g. '127.0.0.1:5000' or '/tmp/sops.sock'",
					Value: "127.0.0.1:5000",
				},
			},
			Action: func(c *cli.Context) error {
				return keyservicecmd.Run(keyservicecmd.Opts{
					Network: c.String("network"),
					Address: c.String("address"),
				})
			},
		},
		{
			Name:  "groups",
			Usage: "modify the groups on a SOPS file",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new group to a SOPS file",
					Flags: append([]cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Usage: "the file to add the group to",
						},
						cli.StringSliceFlag{
							Name:  "pgp",
							Usage: "the PGP fingerprints the new group should contain. Can be specified more than once",
						},
						cli.StringSliceFlag{
							Name:  "kms",
							Usage: "the KMS ARNs the new group should contain. Can be specified more than once",
						},
						cli.BoolFlag{
							Name:  "in-place, i",
							Usage: "write output back to the same file instead of stdout",
						},
						cli.IntFlag{
							Name:  "shamir-secret-sharing-threshold",
							Usage: "the number of master keys required to retrieve the data key with shamir",
						},
						cli.StringFlag{
							Name:  "encryption-context",
							Usage: "comma separated list of KMS encryption context key:value pairs",
						},
					}, keyserviceFlags...),
					Action: func(c *cli.Context) error {
						pgpFps := c.StringSlice("pgp")
						kmsArns := c.StringSlice("kms")
						var group sops.KeyGroup
						for _, fp := range pgpFps {
							group = append(group, pgp.NewMasterKeyFromFingerprint(fp))
						}
						for _, arn := range kmsArns {
							group = append(group, kms.NewMasterKeyFromArn(arn, kms.ParseKMSContext(c.String("encryption-context"))))
						}
						return groups.Add(groups.AddOpts{
							InputPath:      c.String("file"),
							InPlace:        c.Bool("in-place"),
							InputStore:     inputStore(c, c.String("file")),
							OutputStore:    outputStore(c, c.String("file")),
							Group:          group,
							GroupThreshold: c.Int("shamir-secret-sharing-threshold"),
							KeyServices:    keyservices(c),
						})
					},
				},
				{
					Name:  "delete",
					Usage: "delete a key group from a SOPS file",
					Flags: append([]cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Usage: "the file to add the group to",
						},
						cli.BoolFlag{
							Name:  "in-place, i",
							Usage: "write output back to the same file instead of stdout",
						},
						cli.IntFlag{
							Name:  "shamir-secret-sharing-threshold",
							Usage: "the number of master keys required to retrieve the data key with shamir",
						},
					}, keyserviceFlags...),
					ArgsUsage: `[index]`,
					Action: func(c *cli.Context) error {
						group, err := strconv.ParseUint(c.Args().First(), 10, 32)
						if err != nil {
							return fmt.Errorf("failed to parse [index] argument: %s", err)
						}

						return groups.Delete(groups.DeleteOpts{
							InputPath:      c.String("file"),
							InPlace:        c.Bool("in-place"),
							InputStore:     inputStore(c, c.String("file")),
							OutputStore:    outputStore(c, c.String("file")),
							Group:          uint(group),
							GroupThreshold: c.Int("shamir-secret-sharing-threshold"),
							KeyServices:    keyservices(c),
						})
					},
				},
			},
		},
	}
	app.Flags = append([]cli.Flag{
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
		cli.IntFlag{
			Name:  "shamir-secret-sharing-threshold",
			Usage: "the number of master keys required to retrieve the data key with shamir",
		},
	}, keyserviceFlags...)

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.NewExitError("Error: no file specified", codes.NoFileSpecified)
		}
		fileName := c.Args()[0]
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("rm-kms") != "" || c.String("rm-pgp") != "" {
				return cli.NewExitError("Error: cannot add or remove keys on non-existent files, use `--kms` and `--pgp` instead.", 49)
			}
			if c.Bool("encrypt") || c.Bool("decrypt") || c.Bool("rotate") {
				return cli.NewExitError("Error: cannot operate on non-existent file", codes.NoFileSpecified)
			}
		}

		inputStore := inputStore(c, fileName)
		outputStore := outputStore(c, fileName)
		svcs := keyservices(c)

		var output []byte
		var err error
		if c.Bool("encrypt") {
			keyGroups, err := keyGroups(c, fileName)
			if err != nil {
				return err
			}
			shamirThreshold, err := shamirThreshold(c, fileName)
			if err != nil {
				return err
			}
			output, err = encrypt(encryptOpts{
				OutputStore:       outputStore,
				InputStore:        inputStore,
				InputPath:         fileName,
				Cipher:            aes.NewCipher(),
				UnencryptedSuffix: c.String("unencrypted-suffix"),
				KeyServices:       svcs,
				KeyGroups:         keyGroups,
				GroupThreshold:    shamirThreshold,
			})
			if err != nil {
				return err
			}
		}

		if c.Bool("decrypt") {
			extract, err := parseTreePath(c.String("extract"))
			if err != nil {
				return cli.NewExitError(fmt.Errorf("error parsing --extract path: %s", err), codes.InvalidTreePathFormat)
			}
			output, err = decrypt(decryptOpts{
				OutputStore: outputStore,
				InputStore:  inputStore,
				InputPath:   fileName,
				Cipher:      aes.NewCipher(),
				Extract:     extract,
				KeyServices: svcs,
				IgnoreMAC:   c.Bool("ignore-mac"),
			})
			if err != nil {
				return err
			}
		}
		if c.Bool("rotate") {
			var addMasterKeys []keys.MasterKey
			kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
			for _, k := range kms.MasterKeysFromArnString(c.String("add-kms"), kmsEncryptionContext) {
				addMasterKeys = append(addMasterKeys, k)
			}
			for _, k := range pgp.MasterKeysFromFingerprintString(c.String("add-pgp")) {
				addMasterKeys = append(addMasterKeys, k)
			}

			var rmMasterKeys []keys.MasterKey
			for _, k := range kms.MasterKeysFromArnString(c.String("add-kms"), kmsEncryptionContext) {
				rmMasterKeys = append(rmMasterKeys, k)
			}
			for _, k := range pgp.MasterKeysFromFingerprintString(c.String("add-pgp")) {
				rmMasterKeys = append(rmMasterKeys, k)
			}
			output, err = rotate(rotateOpts{
				OutputStore:      outputStore,
				InputStore:       inputStore,
				InputPath:        fileName,
				Cipher:           aes.NewCipher(),
				KeyServices:      svcs,
				IgnoreMAC:        c.Bool("ignore-mac"),
				AddMasterKeys:    addMasterKeys,
				RemoveMasterKeys: rmMasterKeys,
			})
			if err != nil {
				return err
			}
		}

		if c.String("set") != "" {
			path, value, err := extractSetArguments(c.String("set"))
			if err != nil {
				return err
			}
			output, err = set(setOpts{
				OutputStore: outputStore,
				InputStore:  inputStore,
				InputPath:   fileName,
				Cipher:      aes.NewCipher(),
				KeyServices: svcs,
				IgnoreMAC:   c.Bool("ignore-mac"),
				Value:       value,
				TreePath:    path,
			})
			if err != nil {
				return err
			}
		}

		isEditMode := !c.Bool("encrypt") && !c.Bool("decrypt") && !c.Bool("rotate") && c.String("set") == ""
		if isEditMode {
			_, statErr := os.Stat(fileName)
			fileExists := statErr == nil
			opts := editOpts{
				OutputStore:    outputStore,
				InputStore:     inputStore,
				InputPath:      fileName,
				Cipher:         aes.NewCipher(),
				KeyServices:    svcs,
				IgnoreMAC:      c.Bool("ignore-mac"),
				ShowMasterKeys: c.Bool("show-master-keys"),
			}
			if fileExists {
				output, err = edit(opts)
			} else {
				// File doesn't exist, edit the example file instead
				keyGroups, err := keyGroups(c, fileName)
				if err != nil {
					return err
				}
				shamirThreshold, err := shamirThreshold(c, fileName)
				if err != nil {
					return err
				}
				output, err = editExample(editExampleOpts{
					editOpts:          opts,
					UnencryptedSuffix: c.String("unencrypted-suffix"),
					KeyGroups:         keyGroups,
					GroupThreshold:    shamirThreshold,
				})
			}
		}

		if err != nil {
			return err
		}
		// We open the file *after* the operations on the tree have been
		// executed to avoid truncating it when there's errors
		var outputFile *os.File
		if c.Bool("in-place") || isEditMode || c.String("set") != "" {
			var err error
			outputFile, err = os.Create(fileName)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Could not open in-place file for writing: %s", err), codes.CouldNotWriteOutputFile)
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

func keyservices(c *cli.Context) (svcs []keyservice.KeyServiceClient) {
	if c.Bool("enable-local-keyservice") {
		svcs = append(svcs, keyservice.NewLocalClient())
	}
	uris := c.StringSlice("keyservice")
	for _, uri := range uris {
		url, err := url.Parse(uri)
		if err != nil {
			log.Printf("Error parsing keyservice URI %s, skipping", uri)
			continue
		}
		addr := url.Host
		if url.Scheme == "unix" {
			addr = url.Path
		}
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout(url.Scheme, addr, timeout)
			}),
		}
		log.Printf("Connecting to key service %s://%s", url.Scheme, addr)
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		svcs = append(svcs, keyservice.NewKeyServiceClient(conn))
	}
	return
}

func inputStore(context *cli.Context, path string) sops.Store {
	switch context.String("input-type") {
	case "yaml":
		return &yamlstores.Store{}
	case "json":
		return &json.Store{}
	default:
		return defaultStore(path)
	}
}
func outputStore(context *cli.Context, path string) sops.Store {
	switch context.String("output-type") {
	case "yaml":
		return &yamlstores.Store{}
	case "json":
		return &json.Store{}
	default:
		return defaultStore(path)
	}
}

func defaultStore(path string) sops.Store {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return &yamlstores.Store{}
	} else if strings.HasSuffix(path, ".json") {
		return &json.Store{}
	}
	return &json.BinaryStore{}
}

func parseTreePath(arg string) ([]interface{}, error) {
	var path []interface{}
	components := strings.Split(arg, "[")
	for _, component := range components {
		if component == "" {
			continue
		}
		if component[len(component)-1] != ']' {
			return nil, fmt.Errorf("component %s doesn't end with ]", component)
		}
		component = component[:len(component)-1]
		if component[0] == byte('"') || component[0] == byte('\'') {
			// The component is a string
			component = component[1 : len(component)-1]
			path = append(path, component)
		} else {
			// The component must be a number
			i, err := strconv.Atoi(component)
			if err != nil {
				return nil, err
			}
			path = append(path, i)
		}
	}
	return path, nil
}

func keyGroups(c *cli.Context, file string) ([]sops.KeyGroup, error) {
	var kmsKeys []keys.MasterKey
	var pgpKeys []keys.MasterKey
	kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
	if c.String("encryption-context") != "" && kmsEncryptionContext == nil {
		return nil, cli.NewExitError("Invalid KMS encryption context format", codes.ErrorInvalidKMSEncryptionContextFormat)
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
	if c.String("kms") == "" && c.String("pgp") == "" {
		var err error
		var configPath string
		if c.String("config") != "" {
		} else {
			configPath, err = config.FindConfigFile(".")
			if err != nil {
				return nil, fmt.Errorf("config file not found and no keys provided through command line options")
			}
		}
		conf, err := config.LoadForFile(configPath, file, kmsEncryptionContext)
		if err != nil {
			return nil, err
		}
		return conf.KeyGroups, err
	}
	return []sops.KeyGroup{append(kmsKeys, pgpKeys...)}, nil
}

func shamirThreshold(c *cli.Context, file string) (int, error) {
	if c.Int("shamir-secret-sharing-threshold") != 0 {
		return c.Int("shamir-secret-sharing-threshold"), nil
	}
	var err error
	var configPath string
	if c.String("config") != "" {
	} else {
		configPath, err = config.FindConfigFile(".")
		if err != nil {
			return 0, fmt.Errorf("config file not found and no keys provided through command line options")
		}
	}
	conf, err := config.LoadForFile(configPath, file, nil)
	if err != nil {
		return 0, err
	}
	return conf.ShamirThreshold, err
}

func jsonValueToTreeInsertableValue(jsonValue string) (interface{}, error) {
	var valueToInsert interface{}
	err := encodingjson.Unmarshal([]byte(jsonValue), &valueToInsert)
	if err != nil {
		return nil, cli.NewExitError("Value for --set is not valid JSON", codes.ErrorInvalidSetFormat)
	}
	// Check if decoding it as json we find a single value
	// and not a map or slice, in which case we can't marshal
	// it to a sops.TreeBranch
	kind := reflect.ValueOf(valueToInsert).Kind()
	if kind == reflect.Map || kind == reflect.Slice {
		var err error
		valueToInsert, err = (&json.Store{}).Unmarshal([]byte(jsonValue))
		if err != nil {
			return nil, cli.NewExitError("Invalid --set value format", codes.ErrorInvalidSetFormat)
		}
	}
	return valueToInsert, nil
}

func extractSetArguments(set string) (path []interface{}, valueToInsert interface{}, err error) {
	// Set is a string with the format "python-dict-index json-value"
	// Since python-dict-index has to end with ], we split at "] " to get the two parts
	pathValuePair := strings.SplitAfterN(set, "] ", 2)
	if len(pathValuePair) < 2 {
		return nil, nil, cli.NewExitError("Invalid --set format", codes.ErrorInvalidSetFormat)
	}
	fullPath := strings.TrimRight(pathValuePair[0], " ")
	jsonValue := pathValuePair[1]
	valueToInsert, err = jsonValueToTreeInsertableValue(jsonValue)

	path, err = parseTreePath(fullPath)
	if err != nil {
		return nil, nil, cli.NewExitError("Invalid --set format", codes.ErrorInvalidSetFormat)
	}
	return path, valueToInsert, nil
}
