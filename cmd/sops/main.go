package main //import "go.mozilla.org/sops/cmd/sops"

import (
	encodingjson "encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/aes"
	_ "go.mozilla.org/sops/audit"
	"go.mozilla.org/sops/azkv"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/cmd/sops/subcommand/groups"
	keyservicecmd "go.mozilla.org/sops/cmd/sops/subcommand/keyservice"
	"go.mozilla.org/sops/cmd/sops/subcommand/updatekeys"
	"go.mozilla.org/sops/config"
	"go.mozilla.org/sops/gcpkms"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/logging"
	"go.mozilla.org/sops/pgp"
	"go.mozilla.org/sops/stores/dotenv"
	"go.mozilla.org/sops/stores/ini"
	"go.mozilla.org/sops/stores/json"
	yamlstores "go.mozilla.org/sops/stores/yaml"
	"go.mozilla.org/sops/version"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v1"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("CMD")
}

func main() {
	cli.VersionPrinter = version.PrintVersion
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
	app.Usage = "sops - encrypted file editor with AWS KMS, GCP KMS, Azure Key Vault and GPG support"
	app.ArgsUsage = "sops [options] file"
	app.Version = version.Version
	app.Authors = []cli.Author{
		{Name: "Julien Vehent", Email: "jvehent@mozilla.com"},
		{Name: "Adrian Utrilla", Email: "adrianutrilla@gmail.com"},
	}
	app.UsageText = `sops is an editor of encrypted files that supports AWS KMS and PGP

   To encrypt or decrypt a document with AWS KMS, specify the KMS ARN
   in the -k flag or in the SOPS_KMS_ARN environment variable.
   (you need valid credentials in ~/.aws/credentials or in your env)

   To encrypt or decrypt a document with GCP KMS, specify the
   GCP KMS resource ID in the --gcp-kms flag or in the SOPS_GCP_KMS_IDS
   environment variable.
   (you need to setup google application default credentials. See
    https://developers.google.com/identity/protocols/application-default-credentials)

   To encrypt or decrypt a document with Azure Key Vault, specify the
   Azure Key Vault key URL in the --azure-kv flag or in the SOPS_AZURE_KEYVAULT_URL
   environment variable.
   (authentication is based on environment variables, see
    https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-environment-based-authentication.
    The user/sp needs the key/encrypt and key/decrypt permissions)

   To encrypt or decrypt using PGP, specify the PGP fingerprint in the
   -p flag or in the SOPS_PGP_FP environment variable.

   To use multiple KMS or PGP keys, separate them by commas. For example:
       $ sops -p "10F2...0A, 85D...B3F21" file.yaml

   The -p, -k, --gcp-kms and --azure-kv flags are only used to encrypt new documents. Editing
   or decrypting existing documents can be done with "sops file" or
   "sops -d file" respectively. The KMS and PGP keys listed in the encrypted
   documents are used then. To manage master keys in existing documents, use
   the "add-{kms,pgp,gcp-kms,azure-kv}" and "rm-{kms,pgp,gcp-kms,azure-kv}" flags.

   To use a different GPG binary than the one in your PATH, set SOPS_GPG_EXEC.
   To use a GPG key server other than gpg.mozilla.org, set SOPS_GPG_KEYSERVER.

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
				cli.BoolFlag{
					Name:  "prompt",
					Usage: "Prompt user to confirm every incoming request",
				},
				cli.BoolFlag{
					Name:  "verbose",
					Usage: "Enable verbose logging output",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				err := keyservicecmd.Run(keyservicecmd.Opts{
					Network: c.String("network"),
					Address: c.String("address"),
					Prompt:  c.Bool("prompt"),
				})
				if err != nil {
					log.Errorf("Error running keyservice: %s", err)
					return err
				}
				return nil
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
						cli.StringFlag{
							Name:  "aws-profile",
							Usage: "The AWS profile to use for requests to AWS",
						},
						cli.StringSliceFlag{
							Name:  "gcp-kms",
							Usage: "the GCP KMS Resource ID the new group should contain. Can be specified more than once",
						},
						cli.StringSliceFlag{
							Name:  "azure-kv",
							Usage: "the Azure Key Vault key URL the new group should contain. Can be specified more than once",
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
						gcpKmses := c.StringSlice("gcp-kms")
						azkvs := c.StringSlice("azure-kv")
						var group sops.KeyGroup
						for _, fp := range pgpFps {
							group = append(group, pgp.NewMasterKeyFromFingerprint(fp))
						}
						for _, arn := range kmsArns {
							group = append(group, kms.NewMasterKeyFromArn(arn, kms.ParseKMSContext(c.String("encryption-context")), c.String("aws-profile")))
						}
						for _, kms := range gcpKmses {
							group = append(group, gcpkms.NewMasterKeyFromResourceID(kms))
						}
						for _, url := range azkvs {
							k, err := azkv.NewMasterKeyFromURL(url)
							if err != nil {
								log.WithError(err).Error("Failed to add key")
								continue
							}
							group = append(group, k)
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
		{
			Name:      "updatekeys",
			Usage:     "update the keys of a SOPS file using the config file",
			ArgsUsage: `file`,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: `pre-approve all changes and run non-interactively`,
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				configPath, err := config.FindConfigFile(".")
				if err != nil {
					return common.NewExitError(err, codes.ErrorGeneric)
				}
				if c.NArg() < 1 {
					return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
				}
				err = updatekeys.UpdateKeys(updatekeys.Opts{
					InputPath:   c.Args()[0],
					GroupQuorum: c.Int("shamir-secret-sharing-quorum"),
					KeyServices: keyservices(c),
					Interactive: !c.Bool("yes"),
					ConfigPath:  configPath,
				})
				if cliErr, ok := err.(*cli.ExitError); ok && cliErr != nil {
					return cliErr
				} else if err != nil {
					return common.NewExitError(err, codes.ErrorGeneric)
				}
				return nil
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
			Name:  "aws-profile",
			Usage: "The AWS profile to use for requests to AWS",
		},
		cli.StringFlag{
			Name:   "gcp-kms",
			Usage:  "comma separated list of GCP KMS resource IDs",
			EnvVar: "SOPS_GCP_KMS_IDS",
		},
		cli.StringFlag{
			Name:   "azure-kv",
			Usage:  "comma separated list of Azure Key Vault URLs",
			EnvVar: "SOPS_AZURE_KEYVAULT_URLS",
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
			Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
		},
		cli.StringFlag{
			Name:  "output-type",
			Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
		},
		cli.BoolFlag{
			Name:  "show-master-keys, s",
			Usage: "display master encryption keys in the file during editing",
		},
		cli.StringFlag{
			Name:  "add-gcp-kms",
			Usage: "add the provided comma-separated list of GCP KMS key resource IDs to the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "rm-gcp-kms",
			Usage: "remove the provided comma-separated list of GCP KMS key resource IDs from the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "add-azure-kv",
			Usage: "add the provided comma-separated list of Azure Key Vault key URLs to the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "rm-azure-kv",
			Usage: "remove the provided comma-separated list of Azure Key Vault key URLs from the list of master keys on the given file",
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
			Usage: "override the unencrypted key suffix.",
		},
		cli.StringFlag{
			Name:  "encrypted-suffix",
			Usage: "override the encrypted key suffix. When empty, all keys will be encrypted, unless otherwise marked with unencrypted-suffix.",
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
			Usage: `set a specific key or branch in the input document. value must be a json encoded string. (edit mode only). eg. --set '["somekey"][0] {"somevalue":true}'`,
		},
		cli.IntFlag{
			Name:  "shamir-secret-sharing-threshold",
			Usage: "the number of master keys required to retrieve the data key with shamir",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logging output",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "Save the output after encryption or decryption to the file specified",
		},
	}, keyserviceFlags...)

	app.Action = func(c *cli.Context) error {
		if c.Bool("verbose") {
			logging.SetLevel(logrus.DebugLevel)
		}
		if c.NArg() < 1 {
			return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
		}
		if c.Bool("in-place") && c.String("output") != "" {
			return common.NewExitError("Error: cannot operate on both --output and --in-place", codes.ErrorConflictingParameters)
		}
		fileName := c.Args()[0]
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("add-gcp-kms") != "" || c.String("add-azure-kv") != "" ||
				c.String("rm-kms") != "" || c.String("rm-pgp") != "" || c.String("rm-gcp-kms") != "" || c.String("rm-azure-kv") != "" {
				return common.NewExitError("Error: cannot add or remove keys on non-existent files, use `--kms` and `--pgp` instead.", codes.CannotChangeKeysFromNonExistentFile)
			}
			if c.Bool("encrypt") || c.Bool("decrypt") || c.Bool("rotate") {
				return common.NewExitError("Error: cannot operate on non-existent file", codes.NoFileSpecified)
			}
		}

		unencryptedSuffix := c.String("unencrypted-suffix")
		encryptedSuffix := c.String("encrypted-suffix")
		conf, err := loadConfig(c, fileName, nil)
		if err != nil {
			return toExitError(err)
		}
		if conf != nil {
			// command line options have precedence
			if unencryptedSuffix == "" {
				unencryptedSuffix = conf.UnencryptedSuffix
			}
			if encryptedSuffix == "" {
				encryptedSuffix = conf.EncryptedSuffix
			}
		}
		if unencryptedSuffix != "" && encryptedSuffix != "" {
			return common.NewExitError("Error: cannot use both encrypted_suffix and unencrypted_suffix in the same file", codes.ErrorConflictingParameters)
		}
		// only supply the default UnencryptedSuffix when EncryptedSuffix is not provided
		if unencryptedSuffix == "" && encryptedSuffix == "" {
			unencryptedSuffix = sops.DefaultUnencryptedSuffix
		}

		inputStore := inputStore(c, fileName)
		outputStore := outputStore(c, fileName)
		svcs := keyservices(c)

		var output []byte
		if c.Bool("encrypt") {
			var groups []sops.KeyGroup
			groups, err = keyGroups(c, fileName)
			if err != nil {
				return toExitError(err)
			}
			var threshold int
			threshold, err = shamirThreshold(c, fileName)
			if err != nil {
				return toExitError(err)
			}
			output, err = encrypt(encryptOpts{
				OutputStore:       outputStore,
				InputStore:        inputStore,
				InputPath:         fileName,
				Cipher:            aes.NewCipher(),
				UnencryptedSuffix: unencryptedSuffix,
				EncryptedSuffix:   encryptedSuffix,
				KeyServices:       svcs,
				KeyGroups:         groups,
				GroupThreshold:    threshold,
			})
		}

		if c.Bool("decrypt") {
			var extract []interface{}
			extract, err = parseTreePath(c.String("extract"))
			if err != nil {
				return common.NewExitError(fmt.Errorf("error parsing --extract path: %s", err), codes.InvalidTreePathFormat)
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
		}
		if c.Bool("rotate") {
			var addMasterKeys []keys.MasterKey
			kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
			for _, k := range kms.MasterKeysFromArnString(c.String("add-kms"), kmsEncryptionContext, c.String("aws-profile")) {
				addMasterKeys = append(addMasterKeys, k)
			}
			for _, k := range pgp.MasterKeysFromFingerprintString(c.String("add-pgp")) {
				addMasterKeys = append(addMasterKeys, k)
			}
			for _, k := range gcpkms.MasterKeysFromResourceIDString(c.String("add-gcp-kms")) {
				addMasterKeys = append(addMasterKeys, k)
			}
			azureKeys, err := azkv.MasterKeysFromURLs(c.String("add-azure-kv"))
			if err != nil {
				return err
			}
			for _, k := range azureKeys {
				addMasterKeys = append(addMasterKeys, k)
			}

			var rmMasterKeys []keys.MasterKey
			for _, k := range kms.MasterKeysFromArnString(c.String("rm-kms"), kmsEncryptionContext, c.String("aws-profile")) {
				rmMasterKeys = append(rmMasterKeys, k)
			}
			for _, k := range pgp.MasterKeysFromFingerprintString(c.String("rm-pgp")) {
				rmMasterKeys = append(rmMasterKeys, k)
			}
			for _, k := range gcpkms.MasterKeysFromResourceIDString(c.String("rm-gcp-kms")) {
				rmMasterKeys = append(rmMasterKeys, k)
			}
			azureKeys, err = azkv.MasterKeysFromURLs(c.String("rm-azure-kv"))
			if err != nil {
				return err
			}
			for _, k := range azureKeys {
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
		}

		if c.String("set") != "" {
			var path []interface{}
			var value interface{}
			path, value, err = extractSetArguments(c.String("set"))
			if err != nil {
				return toExitError(err)
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
				var groups []sops.KeyGroup
				groups, err = keyGroups(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				var threshold int
				threshold, err = shamirThreshold(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				output, err = editExample(editExampleOpts{
					editOpts:          opts,
					UnencryptedSuffix: unencryptedSuffix,
					EncryptedSuffix:   encryptedSuffix,
					KeyGroups:         groups,
					GroupThreshold:    threshold,
				})
			}
		}

		if err != nil {
			return toExitError(err)
		}

		// We open the file *after* the operations on the tree have been
		// executed to avoid truncating it when there's errors
		if c.Bool("in-place") || isEditMode || c.String("set") != "" {
			file, err := os.Create(fileName)
			if err != nil {
				return common.NewExitError(fmt.Sprintf("Could not open in-place file for writing: %s", err), codes.CouldNotWriteOutputFile)
			}
			defer file.Close()
			_, err = file.Write(output)
			if err != nil {
				return toExitError(err)
			}
			log.Info("File written successfully")
			return nil
		}

		outputFile := os.Stdout
		if c.String("output") != "" {
			file, err := os.Create(c.String("output"))
			if err != nil {
				return common.NewExitError(fmt.Sprintf("Could not open output file for writing: %s", err), codes.CouldNotWriteOutputFile)
			}
			defer file.Close()
			outputFile = file
		}
		_, err = outputFile.Write(output)
		return toExitError(err)
	}
	app.Run(os.Args)
}

func toExitError(err error) error {
	if cliErr, ok := err.(*cli.ExitError); ok && cliErr != nil {
		return cliErr
	} else if err != nil {
		return cli.NewExitError(err, codes.ErrorGeneric)
	}
	return nil
}

func keyservices(c *cli.Context) (svcs []keyservice.KeyServiceClient) {
	if c.Bool("enable-local-keyservice") {
		svcs = append(svcs, keyservice.NewLocalClient())
	}
	uris := c.StringSlice("keyservice")
	for _, uri := range uris {
		url, err := url.Parse(uri)
		if err != nil {
			log.WithField("uri", uri).
				Warnf("Error parsing URI for keyservice, skipping")
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
		log.WithField(
			"address",
			fmt.Sprintf("%s://%s", url.Scheme, addr),
		).Infof("Connecting to key service")
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		svcs = append(svcs, keyservice.NewKeyServiceClient(conn))
	}
	return
}

func inputStore(context *cli.Context, path string) common.Store {
	switch context.String("input-type") {
	case "yaml":
		return &yamlstores.Store{}
	case "json":
		return &json.Store{}
	case "dotenv":
		return &dotenv.Store{}
	case "ini":
		return &ini.Store{}
	case "binary":
		return &json.BinaryStore{}
	default:
		return common.DefaultStoreForPath(path)
	}
}

func outputStore(context *cli.Context, path string) common.Store {
	switch context.String("output-type") {
	case "yaml":
		return &yamlstores.Store{}
	case "json":
		return &json.Store{}
	case "dotenv":
		return &dotenv.Store{}
	case "ini":
		return &ini.Store{}
	case "binary":
		return &json.BinaryStore{}
	default:
		return common.DefaultStoreForPath(path)
	}
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
	var cloudKmsKeys []keys.MasterKey
	var azkvKeys []keys.MasterKey
	kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
	if c.String("encryption-context") != "" && kmsEncryptionContext == nil {
		return nil, common.NewExitError("Invalid KMS encryption context format", codes.ErrorInvalidKMSEncryptionContextFormat)
	}
	if c.String("kms") != "" {
		for _, k := range kms.MasterKeysFromArnString(c.String("kms"), kmsEncryptionContext, c.String("aws-profile")) {
			kmsKeys = append(kmsKeys, k)
		}
	}
	if c.String("gcp-kms") != "" {
		for _, k := range gcpkms.MasterKeysFromResourceIDString(c.String("gcp-kms")) {
			cloudKmsKeys = append(cloudKmsKeys, k)
		}
	}
	if c.String("azure-kv") != "" {
		azureKeys, err := azkv.MasterKeysFromURLs(c.String("azure-kv"))
		if err != nil {
			return nil, err
		}
		for _, k := range azureKeys {
			azkvKeys = append(azkvKeys, k)
		}
	}
	if c.String("pgp") != "" {
		for _, k := range pgp.MasterKeysFromFingerprintString(c.String("pgp")) {
			pgpKeys = append(pgpKeys, k)
		}
	}
	if c.String("kms") == "" && c.String("pgp") == "" && c.String("gcp-kms") == "" && c.String("azure-kv") == "" {
		conf, err := loadConfig(c, file, kmsEncryptionContext)
		// config file might just not be supplied, without any error
		if conf == nil {
			errMsg := "config file not found and no keys provided through command line options"
			if err != nil {
				errMsg = fmt.Sprintf("%s: %s", errMsg, err)
			}
			return nil, fmt.Errorf(errMsg)
		}
		return conf.KeyGroups, err
	}
	var group sops.KeyGroup
	group = append(group, kmsKeys...)
	group = append(group, cloudKmsKeys...)
	group = append(group, azkvKeys...)
	group = append(group, pgpKeys...)
	return []sops.KeyGroup{group}, nil
}

// loadConfig will look for an existing config file, either provided through the command line, or using config.FindConfigFile.
// Since a config file is not required, this function does not error when one is not found, and instead returns a nil config pointer
func loadConfig(c *cli.Context, file string, kmsEncryptionContext map[string]*string) (*config.Config, error) {
	var err error
	var configPath string
	if c.String("config") != "" {
		configPath = c.String("config")
	} else {
		// Ignore config not found errors returned from FindConfigFile since the config file is not mandatory
		configPath, err = config.FindConfigFile(".")
		if err != nil {
			// If we can't find a config file, but we were not explicitly requested to, assume it does not exist
			return nil, nil
		}
	}
	conf, err := config.LoadForFile(configPath, file, kmsEncryptionContext)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func shamirThreshold(c *cli.Context, file string) (int, error) {
	if c.Int("shamir-secret-sharing-threshold") != 0 {
		return c.Int("shamir-secret-sharing-threshold"), nil
	}
	conf, err := loadConfig(c, file, nil)
	if conf == nil {
		// This takes care of the following two case:
		// 1. No config was provided. Err will be nil and ShamirThreshold will be the default value of 0.
		// 2. We did find a config file, but failed to load it. In that case the calling function will print the error and exit.
		return 0, err
	}
	return conf.ShamirThreshold, nil
}

func jsonValueToTreeInsertableValue(jsonValue string) (interface{}, error) {
	var valueToInsert interface{}
	err := encodingjson.Unmarshal([]byte(jsonValue), &valueToInsert)
	if err != nil {
		return nil, common.NewExitError("Value for --set is not valid JSON", codes.ErrorInvalidSetFormat)
	}
	// Check if decoding it as json we find a single value
	// and not a map or slice, in which case we can't marshal
	// it to a sops.TreeBranch
	kind := reflect.ValueOf(valueToInsert).Kind()
	if kind == reflect.Map || kind == reflect.Slice {
		var err error
		valueToInsert, err = (&json.Store{}).LoadPlainFile([]byte(jsonValue))
		if err != nil {
			return nil, common.NewExitError("Invalid --set value format", codes.ErrorInvalidSetFormat)
		}
	}
	return valueToInsert.(sops.TreeBranches)[0], nil
}

func extractSetArguments(set string) (path []interface{}, valueToInsert interface{}, err error) {
	// Set is a string with the format "python-dict-index json-value"
	// Since python-dict-index has to end with ], we split at "] " to get the two parts
	pathValuePair := strings.SplitAfterN(set, "] ", 2)
	if len(pathValuePair) < 2 {
		return nil, nil, common.NewExitError("Invalid --set format", codes.ErrorInvalidSetFormat)
	}
	fullPath := strings.TrimRight(pathValuePair[0], " ")
	jsonValue := pathValuePair[1]
	valueToInsert, err = jsonValueToTreeInsertableValue(jsonValue)

	path, err = parseTreePath(fullPath)
	if err != nil {
		return nil, nil, common.NewExitError("Invalid --set format", codes.ErrorInvalidSetFormat)
	}
	return path, valueToInsert, nil
}
