package main // import "github.com/getsops/sops/v3/cmd/sops"

import (
	"context"
	encodingjson "encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	osExec "os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	_ "github.com/getsops/sops/v3/audit"
	"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/cmd/sops/subcommand/exec"
	filestatuscmd "github.com/getsops/sops/v3/cmd/sops/subcommand/filestatus"
	"github.com/getsops/sops/v3/cmd/sops/subcommand/groups"
	keyservicecmd "github.com/getsops/sops/v3/cmd/sops/subcommand/keyservice"
	publishcmd "github.com/getsops/sops/v3/cmd/sops/subcommand/publish"
	"github.com/getsops/sops/v3/cmd/sops/subcommand/updatekeys"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/gcpkms"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/logging"
	"github.com/getsops/sops/v3/pgp"
	"github.com/getsops/sops/v3/stores/dotenv"
	"github.com/getsops/sops/v3/stores/json"
	"github.com/getsops/sops/v3/version"
)

var (
	log *logrus.Logger

	// Whether the config file warning was already shown to the user.
	// Used and set by findConfigFile().
	showedConfigFileWarning bool
)

func init() {
	log = logging.NewLogger("CMD")
}

func warnMoreThanOnePositionalArgument(c *cli.Context) {
	if c.NArg() > 1 {
		log.Warn("More than one positional argument provided. Only the first one will be used!")
		potentialFlag := ""
		for i, value := range c.Args() {
			if i > 0 && strings.HasPrefix(value, "-") {
				potentialFlag = value
			}
		}
		if potentialFlag != "" {
			log.Warn(fmt.Sprintf("Note that one of the ignored positional argument is %q, which looks like a flag. Flags must always be provided before the first positional argument!", potentialFlag))
		}
	}
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
	app.Usage = "sops - encrypted file editor with AWS KMS, GCP KMS, Azure Key Vault, age, and GPG support"
	app.ArgsUsage = "sops [options] file"
	app.Version = version.Version
	app.Authors = []cli.Author{
		{Name: "AJ Bahnken", Email: "ajvb@mozilla.com"},
		{Name: "Adrian Utrilla", Email: "adrianutrilla@gmail.com"},
		{Name: "Julien Vehent", Email: "jvehent@mozilla.com"},
	}
	app.UsageText = `sops is an editor of encrypted files that supports AWS KMS, GCP, AZKV,
	PGP, and Age

   To encrypt or decrypt a document with AWS KMS, specify the KMS ARN
   in the -k flag or in the SOPS_KMS_ARN environment variable.
   (you need valid credentials in ~/.aws/credentials or in your env)

   To encrypt or decrypt a document with GCP KMS, specify the
   GCP KMS resource ID in the --gcp-kms flag or in the SOPS_GCP_KMS_IDS
   environment variable.
   (You need to setup Google application default credentials. See
    https://developers.google.com/identity/protocols/application-default-credentials)


   To encrypt or decrypt a document with HashiCorp Vault's Transit Secret
   Engine, specify the Vault key URI name in the --hc-vault-transit flag
   or in the SOPS_VAULT_URIS environment variable (for example
   https://vault.example.org:8200/v1/transit/keys/dev, where
   'https://vault.example.org:8200' is the vault server, 'transit' the
   enginePath, and 'dev' is the name of the key).
   (You need to enable the Transit Secrets Engine in Vault. See
    https://www.vaultproject.io/docs/secrets/transit/index.html)

   To encrypt or decrypt a document with Azure Key Vault, specify the
   Azure Key Vault key URL in the --azure-kv flag or in the
   SOPS_AZURE_KEYVAULT_URL environment variable.
   (Authentication is based on environment variables, see
    https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-environment-based-authentication.
    The user/sp needs the key/encrypt and key/decrypt permissions.)

   To encrypt or decrypt using age, specify the recipient in the -a flag,
   or in the SOPS_AGE_RECIPIENTS environment variable.

   To encrypt or decrypt using PGP, specify the PGP fingerprint in the
   -p flag or in the SOPS_PGP_FP environment variable.

   To use multiple KMS or PGP keys, separate them by commas. For example:
       $ sops -p "10F2...0A, 85D...B3F21" file.yaml

   The -p, -k, --gcp-kms, --hc-vault-transit, and --azure-kv flags are only
   used to encrypt new documents. Editing or decrypting existing documents
   can be done with "sops file" or "sops decrypt file" respectively. The KMS and
   PGP keys listed in the encrypted documents are used then. To manage master
   keys in existing documents, use the "add-{kms,pgp,gcp-kms,azure-kv,hc-vault-transit}"
   and "rm-{kms,pgp,gcp-kms,azure-kv,hc-vault-transit}" flags with --rotate
   or the updatekeys command.

   To use a different GPG binary than the one in your PATH, set SOPS_GPG_EXEC.

   To select a different editor than the default (vim), set SOPS_EDITOR or
   EDITOR.

   Note that flags must always be provided before the filename to operate on.
   Otherwise, they will be ignored.

   For more information, see the README at https://github.com/getsops/sops`
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:      "exec-env",
			Usage:     "execute a command with decrypted values inserted into the environment",
			ArgsUsage: "[file to decrypt] [command to run]",
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "background",
					Usage: "background the process and don't wait for it to complete (DEPRECATED)",
				},
				cli.BoolFlag{
					Name:  "pristine",
					Usage: "insert only the decrypted values into the environment without forwarding existing environment variables",
				},
				cli.StringFlag{
					Name:  "user",
					Usage: "the user to run the command as",
				},
				cli.BoolFlag{
					Name:  "same-process",
					Usage: "run command in the current process instead of in a child process",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.NArg() != 2 {
					return common.NewExitError(fmt.Errorf("error: missing file to decrypt"), codes.ErrorGeneric)
				}

				fileName := c.Args()[0]
				command := c.Args()[1]

				inputStore, err := inputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}

				svcs := keyservices(c)

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}
				opts := decryptOpts{
					OutputStore:     &dotenv.Store{},
					InputStore:      inputStore,
					InputPath:       fileName,
					Cipher:          aes.NewCipher(),
					KeyServices:     svcs,
					DecryptionOrder: order,
					IgnoreMAC:       c.Bool("ignore-mac"),
				}

				if c.Bool("background") {
					log.Warn("exec-env's --background option is deprecated and will be removed in a future version of sops")

					if c.Bool("same-process") {
						return common.NewExitError("Error: The --same-process flag cannot be used with --background", codes.ErrorConflictingParameters)
					}
				}

				tree, err := decryptTree(opts)
				if err != nil {
					return toExitError(err)
				}

				var env []string
				for _, item := range tree.Branches[0] {
					if dotenv.IsComplexValue(item.Value) {
						return cli.NewExitError(fmt.Errorf("cannot use complex value in environment: %s", item.Value), codes.ErrorGeneric)
					}
					if _, ok := item.Key.(sops.Comment); ok {
						continue
					}
					key, ok := item.Key.(string)
					if !ok {
						return cli.NewExitError(fmt.Errorf("cannot use non-string keys in environment, got %T", item.Key), codes.ErrorGeneric)
					}
					if strings.Contains(key, "=") {
						return cli.NewExitError(fmt.Errorf("cannot use keys with '=' in environment: %s", key), codes.ErrorGeneric)
					}
					value, ok := item.Value.(string)
					if !ok {
						return cli.NewExitError(fmt.Errorf("cannot use non-string values in environment, got %T", item.Value), codes.ErrorGeneric)
					}
					env = append(env, fmt.Sprintf("%s=%s", key, value))
				}

				if err := exec.ExecWithEnv(exec.ExecOpts{
					Command:     command,
					Plaintext:   []byte{},
					Background:  c.Bool("background"),
					Pristine:    c.Bool("pristine"),
					User:        c.String("user"),
					SameProcess: c.Bool("same-process"),
					Env:         env,
				}); err != nil {
					return toExitError(err)
				}

				return nil
			},
		},
		{
			Name:      "exec-file",
			Usage:     "execute a command with the decrypted contents as a temporary file",
			ArgsUsage: "[file to decrypt] [command to run]",
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "background",
					Usage: "background the process and don't wait for it to complete (DEPRECATED)",
				},
				cli.BoolFlag{
					Name:  "no-fifo",
					Usage: "use a regular file instead of a fifo to temporarily hold the decrypted contents",
				},
				cli.StringFlag{
					Name:  "user",
					Usage: "the user to run the command as",
				},
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently ini, json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
				cli.StringFlag{
					Name:  "output-type",
					Usage: "currently ini, json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
				},
				cli.StringFlag{
					Name:  "filename",
					Usage: fmt.Sprintf("filename for the temporarily file (default: %s)", exec.FallbackFilename),
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.NArg() != 2 {
					return common.NewExitError(fmt.Errorf("error: missing file to decrypt"), codes.ErrorGeneric)
				}

				fileName := c.Args()[0]
				command := c.Args()[1]

				inputStore, err := inputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}

				svcs := keyservices(c)

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}
				opts := decryptOpts{
					OutputStore:     outputStore,
					InputStore:      inputStore,
					InputPath:       fileName,
					Cipher:          aes.NewCipher(),
					KeyServices:     svcs,
					DecryptionOrder: order,
					IgnoreMAC:       c.Bool("ignore-mac"),
				}

				output, err := decrypt(opts)
				if err != nil {
					return toExitError(err)
				}

				if c.Bool("background") {
					log.Warn("exec-file's --background option is deprecated and will be removed in a future version of sops")
				}

				if err := exec.ExecWithFile(exec.ExecOpts{
					Command:    command,
					Plaintext:  output,
					Background: c.Bool("background"),
					Fifo:       !c.Bool("no-fifo"),
					User:       c.String("user"),
					Filename:   c.String("filename"),
				}); err != nil {
					return toExitError(err)
				}

				return nil
			},
		},
		{
			Name:      "publish",
			Usage:     "Publish sops file or directory to a configured destination",
			ArgsUsage: `file`,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: `pre-approve all changes and run non-interactively`,
				},
				cli.BoolFlag{
					Name:  "omit-extensions",
					Usage: "Omit file extensions in destination path when publishing sops file to configured destinations",
				},
				cli.BoolFlag{
					Name:  "recursive",
					Usage: "If the source path is a directory, publish all its content recursively",
				},
				cli.BoolFlag{
					Name:  "verbose",
					Usage: "Enable verbose logging output",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") || c.GlobalBool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				var configPath string
				var err error
				if c.GlobalString("config") != "" {
					configPath = c.GlobalString("config")
				} else {
					configPath, err = findConfigFile()
					if err != nil {
						return common.NewExitError(err, codes.ErrorGeneric)
					}
				}
				if c.NArg() < 1 {
					return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
				}
				warnMoreThanOnePositionalArgument(c)
				path := c.Args()[0]
				info, err := os.Stat(path)
				if err != nil {
					return toExitError(err)
				}
				if info.IsDir() && !c.Bool("recursive") {
					return fmt.Errorf("can't operate on a directory without --recursive flag.")
				}
				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}
				err = filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
					if err != nil {
						return toExitError(err)
					}
					if !info.IsDir() {
						inputStore, err := inputStore(c, subPath)
						if err != nil {
							return toExitError(err)
						}
						err = publishcmd.Run(publishcmd.Opts{
							ConfigPath:      configPath,
							InputPath:       subPath,
							Cipher:          aes.NewCipher(),
							KeyServices:     keyservices(c),
							DecryptionOrder: order,
							InputStore:      inputStore,
							Interactive:     !c.Bool("yes"),
							OmitExtensions:  c.Bool("omit-extensions"),
							Recursive:       c.Bool("recursive"),
						})
						if cliErr, ok := err.(*cli.ExitError); ok && cliErr != nil {
							return cliErr
						} else if err != nil {
							return common.NewExitError(err, codes.ErrorGeneric)
						}
					}
					return nil
				})
				if err != nil {
					return toExitError(err)
				}
				return nil
			},
		},
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
				if c.Bool("verbose") || c.GlobalBool("verbose") {
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
			Name:      "filestatus",
			Usage:     "check the status of the file, returning encryption status",
			ArgsUsage: `file`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently ini, json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
				}

				fileName := c.Args()[0]
				inputStore, err := inputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				opts := filestatuscmd.Opts{
					InputStore: inputStore,
					InputPath:  fileName,
				}

				status, err := filestatuscmd.FileStatus(opts)
				if err != nil {
					return err
				}

				json, err := encodingjson.Marshal(status)
				if err != nil {
					return common.NewExitError(err, codes.ErrorGeneric)
				}

				fmt.Println(string(json))

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
						cli.StringSliceFlag{
							Name:  "hc-vault-transit",
							Usage: "the full vault path to the key used to encrypt/decrypt. Make you choose and configure a key with encryption/decryption enabled (e.g. 'https://vault.example.org:8200/v1/transit/keys/dev'). Can be specified more than once",
						},
						cli.StringSliceFlag{
							Name:  "age",
							Usage: "the age recipient the new group should contain. Can be specified more than once",
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
						vaultURIs := c.StringSlice("hc-vault-transit")
						azkvs := c.StringSlice("azure-kv")
						ageRecipients := c.StringSlice("age")
						if c.NArg() != 0 {
							return common.NewExitError(fmt.Errorf("error: no positional arguments allowed"), codes.ErrorGeneric)
						}
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
						for _, uri := range vaultURIs {
							k, err := hcvault.NewMasterKeyFromURI(uri)
							if err != nil {
								log.WithError(err).Error("Failed to add key")
								continue
							}
							group = append(group, k)
						}
						for _, url := range azkvs {
							k, err := azkv.NewMasterKeyFromURL(url)
							if err != nil {
								log.WithError(err).Error("Failed to add key")
								continue
							}
							group = append(group, k)
						}
						for _, recipient := range ageRecipients {
							keys, err := age.MasterKeysFromRecipients(recipient)
							if err != nil {
								log.WithError(err).Error("Failed to add key")
								continue
							}
							for _, key := range keys {
								group = append(group, key)
							}
						}
						inputStore, err := inputStore(c, c.String("file"))
						if err != nil {
							return toExitError(err)
						}
						outputStore, err := outputStore(c, c.String("file"))
						if err != nil {
							return toExitError(err)
						}
						return groups.Add(groups.AddOpts{
							InputPath:      c.String("file"),
							InPlace:        c.Bool("in-place"),
							InputStore:     inputStore,
							OutputStore:    outputStore,
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
						if c.NArg() != 1 {
							return common.NewExitError(fmt.Errorf("error: exactly one positional argument (index) required"), codes.ErrorGeneric)
						}
						group, err := strconv.ParseUint(c.Args().First(), 10, 32)
						if err != nil {
							return fmt.Errorf("failed to parse [index] argument: %s", err)
						}

						inputStore, err := inputStore(c, c.String("file"))
						if err != nil {
							return toExitError(err)
						}
						outputStore, err := outputStore(c, c.String("file"))
						if err != nil {
							return toExitError(err)
						}
						return groups.Delete(groups.DeleteOpts{
							InputPath:      c.String("file"),
							InPlace:        c.Bool("in-place"),
							InputStore:     inputStore,
							OutputStore:    outputStore,
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
			Usage:     "update the keys of SOPS files using the config file",
			ArgsUsage: `file`,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: `pre-approve all changes and run non-interactively`,
				},
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently ini, json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				var err error
				var configPath string
				if c.GlobalString("config") != "" {
					configPath = c.GlobalString("config")
				} else {
					configPath, err = findConfigFile()
					if err != nil {
						return common.NewExitError(err, codes.ErrorGeneric)
					}
				}
				if c.NArg() < 1 {
					return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
				}
				failedCounter := 0
				for _, path := range c.Args() {
					err := updatekeys.UpdateKeys(updatekeys.Opts{
						InputPath:       path,
						ShamirThreshold: c.Int("shamir-secret-sharing-threshold"),
						KeyServices:     keyservices(c),
						Interactive:     !c.Bool("yes"),
						ConfigPath:      configPath,
						InputType:       c.String("input-type"),
					})

					if c.NArg() == 1 {
						// a single argument was given, keep compatibility of the error
						if cliErr, ok := err.(*cli.ExitError); ok && cliErr != nil {
							return cliErr
						} else if err != nil {
							return common.NewExitError(err, codes.ErrorGeneric)
						}
					}

					// multiple arguments given (patched functionality),
					// finish updating of remaining files and fail afterwards
					if err != nil {
						failedCounter++
						log.Error(err)
					}
				}
				if failedCounter > 0 {
					return common.NewExitError(fmt.Errorf("failed updating %d key(s)", failedCounter), codes.ErrorGeneric)
				}
				return nil
			},
		},
		{
			Name:      "decrypt",
			Usage:     "decrypt a file, and output the results to stdout. If no filename is provided, stdin will be used.",
			ArgsUsage: `[file]`,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "in-place, i",
					Usage: "write output back to the same file instead of stdout",
				},
				cli.StringFlag{
					Name:  "extract",
					Usage: "extract a specific key or branch from the input document. Example: --extract '[\"somekey\"][0]'",
				},
				cli.StringFlag{
					Name:  "output",
					Usage: "Save the output after decryption to the file specified",
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
					Name:  "ignore-mac",
					Usage: "ignore Message Authentication Code during decryption",
				},
				cli.StringFlag{
					Name:  "filename-override",
					Usage: "Use this filename instead of the provided argument for loading configuration, and for determining input type and output type. Should be provided when reading from stdin.",
				},
				cli.StringFlag{
					Name:   "decryption-order",
					Usage:  "comma separated list of decryption key types",
					EnvVar: "SOPS_DECRYPTION_ORDER",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				readFromStdin := c.NArg() == 0
				if readFromStdin && c.Bool("in-place") {
					return common.NewExitError("Error: cannot use --in-place when reading from stdin", codes.ErrorConflictingParameters)
				}
				warnMoreThanOnePositionalArgument(c)
				if c.Bool("in-place") && c.String("output") != "" {
					return common.NewExitError("Error: cannot operate on both --output and --in-place", codes.ErrorConflictingParameters)
				}
				var fileName string
				var err error
				if !readFromStdin {
					fileName, err = filepath.Abs(c.Args()[0])
					if err != nil {
						return toExitError(err)
					}
					if _, err := os.Stat(fileName); os.IsNotExist(err) {
						return common.NewExitError(fmt.Sprintf("Error: cannot operate on non-existent file %q", fileName), codes.NoFileSpecified)
					}
				}
				fileNameOverride := c.String("filename-override")
				if fileNameOverride == "" {
					fileNameOverride = fileName
				} else {
					fileNameOverride, err = filepath.Abs(fileNameOverride)
					if err != nil {
						return toExitError(err)
					}
				}

				inputStore, err := inputStore(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				svcs := keyservices(c)

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}

				var extract []interface{}
				extract, err = parseTreePath(c.String("extract"))
				if err != nil {
					return common.NewExitError(fmt.Errorf("error parsing --extract path: %s", err), codes.InvalidTreePathFormat)
				}
				output, err := decrypt(decryptOpts{
					OutputStore:     outputStore,
					InputStore:      inputStore,
					InputPath:       fileName,
					ReadFromStdin:   readFromStdin,
					Cipher:          aes.NewCipher(),
					Extract:         extract,
					KeyServices:     svcs,
					DecryptionOrder: order,
					IgnoreMAC:       c.Bool("ignore-mac"),
				})
				if err != nil {
					return toExitError(err)
				}

				// We open the file *after* the operations on the tree have been
				// executed to avoid truncating it when there's errors
				if c.Bool("in-place") {
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
			},
		},
		{
			Name:      "encrypt",
			Usage:     "encrypt a file, and output the results to stdout. If no filename is provided, stdin will be used.",
			ArgsUsage: `[file]`,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "in-place, i",
					Usage: "write output back to the same file instead of stdout",
				},
				cli.StringFlag{
					Name:  "output",
					Usage: "Save the output after decryption to the file specified",
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
					Name:   "hc-vault-transit",
					Usage:  "comma separated list of vault's key URI (e.g. 'https://vault.example.org:8200/v1/transit/keys/dev')",
					EnvVar: "SOPS_VAULT_URIS",
				},
				cli.StringFlag{
					Name:   "pgp, p",
					Usage:  "comma separated list of PGP fingerprints",
					EnvVar: "SOPS_PGP_FP",
				},
				cli.StringFlag{
					Name:   "age, a",
					Usage:  "comma separated list of age recipients",
					EnvVar: "SOPS_AGE_RECIPIENTS",
				},
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
				cli.StringFlag{
					Name:  "output-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
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
					Name:  "unencrypted-regex",
					Usage: "set the unencrypted key regex. When specified, only keys matching the regex will be left unencrypted.",
				},
				cli.StringFlag{
					Name:  "encrypted-regex",
					Usage: "set the encrypted key regex. When specified, only keys matching the regex will be encrypted.",
				},
				cli.StringFlag{
					Name:  "encryption-context",
					Usage: "comma separated list of KMS encryption context key:value pairs",
				},
				cli.IntFlag{
					Name:  "shamir-secret-sharing-threshold",
					Usage: "the number of master keys required to retrieve the data key with shamir",
				},
				cli.StringFlag{
					Name:  "filename-override",
					Usage: "Use this filename instead of the provided argument for loading configuration, and for determining input type and output type. Required when reading from stdin.",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				readFromStdin := c.NArg() == 0
				if readFromStdin {
					if c.Bool("in-place") {
						return common.NewExitError("Error: cannot use --in-place when reading from stdin", codes.ErrorConflictingParameters)
					}
					if c.String("filename-override") == "" {
						return common.NewExitError("Error: must specify --filename-override when reading from stdin", codes.ErrorConflictingParameters)
					}
				}
				warnMoreThanOnePositionalArgument(c)
				if c.Bool("in-place") && c.String("output") != "" {
					return common.NewExitError("Error: cannot operate on both --output and --in-place", codes.ErrorConflictingParameters)
				}
				var fileName string
				var err error
				if !readFromStdin {
					fileName, err = filepath.Abs(c.Args()[0])
					if err != nil {
						return toExitError(err)
					}
					if _, err := os.Stat(fileName); os.IsNotExist(err) {
						return common.NewExitError(fmt.Sprintf("Error: cannot operate on non-existent file %q", fileName), codes.NoFileSpecified)
					}
				}
				fileNameOverride := c.String("filename-override")
				if fileNameOverride == "" {
					fileNameOverride = fileName
				} else {
					fileNameOverride, err = filepath.Abs(fileNameOverride)
					if err != nil {
						return toExitError(err)
					}
				}

				inputStore, err := inputStore(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				svcs := keyservices(c)

				encConfig, err := getEncryptConfig(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				output, err := encrypt(encryptOpts{
					OutputStore:   outputStore,
					InputStore:    inputStore,
					InputPath:     fileName,
					ReadFromStdin: readFromStdin,
					Cipher:        aes.NewCipher(),
					KeyServices:   svcs,
					encryptConfig: encConfig,
				})

				if err != nil {
					return toExitError(err)
				}

				// We open the file *after* the operations on the tree have been
				// executed to avoid truncating it when there's errors
				if c.Bool("in-place") {
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
			},
		},
		{
			Name:      "rotate",
			Usage:     "generate a new data encryption key and reencrypt all values with the new key",
			ArgsUsage: `file`,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "in-place, i",
					Usage: "write output back to the same file instead of stdout",
				},
				cli.StringFlag{
					Name:  "output",
					Usage: "Save the output after decryption to the file specified",
				},
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
				cli.StringFlag{
					Name:  "output-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
				},
				cli.StringFlag{
					Name:  "encryption-context",
					Usage: "comma separated list of KMS encryption context key:value pairs",
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
					Name:  "add-hc-vault-transit",
					Usage: "add the provided comma-separated list of Vault's URI key to the list of master keys on the given file ( eg. https://vault.example.org:8200/v1/transit/keys/dev)",
				},
				cli.StringFlag{
					Name:  "rm-hc-vault-transit",
					Usage: "remove the provided comma-separated list of Vault's URI key from the list of master keys on the given file ( eg. https://vault.example.org:8200/v1/transit/keys/dev)",
				},
				cli.StringFlag{
					Name:  "add-age",
					Usage: "add the provided comma-separated list of age recipients fingerprints to the list of master keys on the given file",
				},
				cli.StringFlag{
					Name:  "rm-age",
					Usage: "remove the provided comma-separated list of age recipients from the list of master keys on the given file",
				},
				cli.StringFlag{
					Name:  "add-pgp",
					Usage: "add the provided comma-separated list of PGP fingerprints to the list of master keys on the given file",
				},
				cli.StringFlag{
					Name:  "rm-pgp",
					Usage: "remove the provided comma-separated list of PGP fingerprints from the list of master keys on the given file",
				},
				cli.StringFlag{
					Name:  "filename-override",
					Usage: "Use this filename instead of the provided argument for loading configuration, and for determining input type and output type",
				},
				cli.StringFlag{
					Name:   "decryption-order",
					Usage:  "comma separated list of decryption key types",
					EnvVar: "SOPS_DECRYPTION_ORDER",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				if c.NArg() < 1 {
					return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
				}
				warnMoreThanOnePositionalArgument(c)
				if c.Bool("in-place") && c.String("output") != "" {
					return common.NewExitError("Error: cannot operate on both --output and --in-place", codes.ErrorConflictingParameters)
				}
				fileName, err := filepath.Abs(c.Args()[0])
				if err != nil {
					return toExitError(err)
				}
				if _, err := os.Stat(fileName); os.IsNotExist(err) {
					if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("add-gcp-kms") != "" || c.String("add-hc-vault-transit") != "" || c.String("add-azure-kv") != "" || c.String("add-age") != "" ||
						c.String("rm-kms") != "" || c.String("rm-pgp") != "" || c.String("rm-gcp-kms") != "" || c.String("rm-hc-vault-transit") != "" || c.String("rm-azure-kv") != "" || c.String("rm-age") != "" {
						return common.NewExitError(fmt.Sprintf("Error: cannot add or remove keys on non-existent file %q, use the `edit` subcommand instead.", fileName), codes.CannotChangeKeysFromNonExistentFile)
					}
				}
				fileNameOverride := c.String("filename-override")
				if fileNameOverride == "" {
					fileNameOverride = fileName
				} else {
					fileNameOverride, err = filepath.Abs(fileNameOverride)
					if err != nil {
						return toExitError(err)
					}
				}

				inputStore, err := inputStore(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				svcs := keyservices(c)

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}

				rotateOpts, err := getRotateOpts(c, fileName, inputStore, outputStore, svcs, order)
				if err != nil {
					return toExitError(err)
				}
				output, err := rotate(rotateOpts)
				if err != nil {
					return toExitError(err)
				}

				// We open the file *after* the operations on the tree have been
				// executed to avoid truncating it when there's errors
				if c.Bool("in-place") {
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
			},
		},
		{
			Name:      "edit",
			Usage:     "edit an encrypted file",
			ArgsUsage: `file`,
			Flags: append([]cli.Flag{
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
					Name:   "hc-vault-transit",
					Usage:  "comma separated list of vault's key URI (e.g. 'https://vault.example.org:8200/v1/transit/keys/dev')",
					EnvVar: "SOPS_VAULT_URIS",
				},
				cli.StringFlag{
					Name:   "pgp, p",
					Usage:  "comma separated list of PGP fingerprints",
					EnvVar: "SOPS_PGP_FP",
				},
				cli.StringFlag{
					Name:   "age, a",
					Usage:  "comma separated list of age recipients",
					EnvVar: "SOPS_AGE_RECIPIENTS",
				},
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
				cli.StringFlag{
					Name:  "output-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
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
					Name:  "unencrypted-regex",
					Usage: "set the unencrypted key regex. When specified, only keys matching the regex will be left unencrypted.",
				},
				cli.StringFlag{
					Name:  "encrypted-regex",
					Usage: "set the encrypted key regex. When specified, only keys matching the regex will be encrypted.",
				},
				cli.StringFlag{
					Name:  "encryption-context",
					Usage: "comma separated list of KMS encryption context key:value pairs",
				},
				cli.IntFlag{
					Name:  "shamir-secret-sharing-threshold",
					Usage: "the number of master keys required to retrieve the data key with shamir",
				},
				cli.BoolFlag{
					Name:  "show-master-keys, s",
					Usage: "display master encryption keys in the file during editing",
				},
				cli.BoolFlag{
					Name:  "ignore-mac",
					Usage: "ignore Message Authentication Code during decryption",
				},
				cli.StringFlag{
					Name:   "decryption-order",
					Usage:  "comma separated list of decryption key types",
					EnvVar: "SOPS_DECRYPTION_ORDER",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				if c.NArg() < 1 {
					return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
				}
				warnMoreThanOnePositionalArgument(c)
				fileName, err := filepath.Abs(c.Args()[0])
				if err != nil {
					return toExitError(err)
				}

				inputStore, err := inputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				svcs := keyservices(c)

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}
				var output []byte
				_, statErr := os.Stat(fileName)
				fileExists := statErr == nil
				opts := editOpts{
					OutputStore:     outputStore,
					InputStore:      inputStore,
					InputPath:       fileName,
					Cipher:          aes.NewCipher(),
					KeyServices:     svcs,
					DecryptionOrder: order,
					IgnoreMAC:       c.Bool("ignore-mac"),
					ShowMasterKeys:  c.Bool("show-master-keys"),
				}
				if fileExists {
					output, err = edit(opts)
					if err != nil {
						return toExitError(err)
					}
				} else {
					// File doesn't exist, edit the example file instead
					encConfig, err := getEncryptConfig(c, fileName)
					if err != nil {
						return toExitError(err)
					}
					output, err = editExample(editExampleOpts{
						editOpts:      opts,
						encryptConfig: encConfig,
					})
					if err != nil {
						return toExitError(err)
					}
				}

				// We open the file *after* the operations on the tree have been
				// executed to avoid truncating it when there's errors
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
			},
		},
		{
			Name:      "set",
			Usage:     `set a specific key or branch in the input document. value must be a json encoded string. eg. '/path/to/file ["somekey"][0] {"somevalue":true}'`,
			ArgsUsage: `file index value`,
			Flags: append([]cli.Flag{
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
				cli.StringFlag{
					Name:  "output-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
				},
				cli.IntFlag{
					Name:  "shamir-secret-sharing-threshold",
					Usage: "the number of master keys required to retrieve the data key with shamir",
				},
				cli.BoolFlag{
					Name:  "ignore-mac",
					Usage: "ignore Message Authentication Code during decryption",
				},
				cli.StringFlag{
					Name:   "decryption-order",
					Usage:  "comma separated list of decryption key types",
					EnvVar: "SOPS_DECRYPTION_ORDER",
				},
				cli.BoolFlag{
					Name:  "idempotent",
					Usage: "do nothing if the given index already has the given value",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				if c.NArg() != 3 {
					return common.NewExitError("Error: no file specified, or index and value are missing", codes.NoFileSpecified)
				}
				fileName, err := filepath.Abs(c.Args()[0])
				if err != nil {
					return toExitError(err)
				}

				inputStore, err := inputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				svcs := keyservices(c)

				path, err := parseTreePath(c.Args()[1])
				if err != nil {
					return common.NewExitError("Invalid set index format", codes.ErrorInvalidSetFormat)
				}

				value, err := jsonValueToTreeInsertableValue(c.Args()[2])
				if err != nil {
					return toExitError(err)
				}

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}
				output, changed, err := set(setOpts{
					OutputStore:     outputStore,
					InputStore:      inputStore,
					InputPath:       fileName,
					Cipher:          aes.NewCipher(),
					KeyServices:     svcs,
					DecryptionOrder: order,
					IgnoreMAC:       c.Bool("ignore-mac"),
					Value:           value,
					TreePath:        path,
				})
				if err != nil {
					return toExitError(err)
				}

				if !changed && c.Bool("idempotent") {
					log.Info("File not written due to no change")
					return nil
				}

				// We open the file *after* the operations on the tree have been
				// executed to avoid truncating it when there's errors
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
			},
		},
		{
			Name:      "unset",
			Usage:     `unset a specific key or branch in the input document.`,
			ArgsUsage: `file index`,
			Flags: append([]cli.Flag{
				cli.StringFlag{
					Name:  "input-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the file's extension to determine the type",
				},
				cli.StringFlag{
					Name:  "output-type",
					Usage: "currently json, yaml, dotenv and binary are supported. If not set, sops will use the input file's extension to determine the output format",
				},
				cli.IntFlag{
					Name:  "shamir-secret-sharing-threshold",
					Usage: "the number of master keys required to retrieve the data key with shamir",
				},
				cli.BoolFlag{
					Name:  "ignore-mac",
					Usage: "ignore Message Authentication Code during decryption",
				},
				cli.StringFlag{
					Name:   "decryption-order",
					Usage:  "comma separated list of decryption key types",
					EnvVar: "SOPS_DECRYPTION_ORDER",
				},
				cli.BoolFlag{
					Name:  "idempotent",
					Usage: "do nothing if the given index does not exist",
				},
			}, keyserviceFlags...),
			Action: func(c *cli.Context) error {
				if c.Bool("verbose") {
					logging.SetLevel(logrus.DebugLevel)
				}
				if c.NArg() != 2 {
					return common.NewExitError("Error: no file specified, or index is missing", codes.NoFileSpecified)
				}
				fileName, err := filepath.Abs(c.Args()[0])
				if err != nil {
					return toExitError(err)
				}

				inputStore, err := inputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				outputStore, err := outputStore(c, fileName)
				if err != nil {
					return toExitError(err)
				}
				svcs := keyservices(c)

				path, err := parseTreePath(c.Args()[1])
				if err != nil {
					return common.NewExitError("Invalid unset index format", codes.ErrorInvalidSetFormat)
				}

				order, err := decryptionOrder(c.String("decryption-order"))
				if err != nil {
					return toExitError(err)
				}
				output, err := unset(unsetOpts{
					OutputStore:     outputStore,
					InputStore:      inputStore,
					InputPath:       fileName,
					Cipher:          aes.NewCipher(),
					KeyServices:     svcs,
					DecryptionOrder: order,
					IgnoreMAC:       c.Bool("ignore-mac"),
					TreePath:        path,
				})
				if err != nil {
					if _, ok := err.(*sops.SopsKeyNotFound); ok && c.Bool("idempotent") {
						return nil
					}
					return toExitError(err)
				}

				// We open the file *after* the operations on the tree have been
				// executed to avoid truncating it when there's errors
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
		cli.BoolFlag{
			Name:   "disable-version-check",
			Usage:  "do not check whether the current version is latest during --version",
			EnvVar: "SOPS_DISABLE_VERSION_CHECK",
		},
		cli.BoolFlag{
			Name:  "check-for-updates",
			Usage: "do check whether the current version is latest during --version",
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
			Name:   "hc-vault-transit",
			Usage:  "comma separated list of vault's key URI (e.g. 'https://vault.example.org:8200/v1/transit/keys/dev')",
			EnvVar: "SOPS_VAULT_URIS",
		},
		cli.StringFlag{
			Name:   "pgp, p",
			Usage:  "comma separated list of PGP fingerprints",
			EnvVar: "SOPS_PGP_FP",
		},
		cli.StringFlag{
			Name:   "age, a",
			Usage:  "comma separated list of age recipients",
			EnvVar: "SOPS_AGE_RECIPIENTS",
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
			Name:  "add-hc-vault-transit",
			Usage: "add the provided comma-separated list of Vault's URI key to the list of master keys on the given file ( eg. https://vault.example.org:8200/v1/transit/keys/dev)",
		},
		cli.StringFlag{
			Name:  "rm-hc-vault-transit",
			Usage: "remove the provided comma-separated list of Vault's URI key from the list of master keys on the given file ( eg. https://vault.example.org:8200/v1/transit/keys/dev)",
		},
		cli.StringFlag{
			Name:  "add-age",
			Usage: "add the provided comma-separated list of age recipients fingerprints to the list of master keys on the given file",
		},
		cli.StringFlag{
			Name:  "rm-age",
			Usage: "remove the provided comma-separated list of age recipients from the list of master keys on the given file",
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
		cli.BoolFlag{
			Name:  "mac-only-encrypted",
			Usage: "compute MAC only over values which end up encrypted",
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
			Name:  "unencrypted-regex",
			Usage: "set the unencrypted key regex. When specified, only keys matching the regex will be left unencrypted.",
		},
		cli.StringFlag{
			Name:  "encrypted-regex",
			Usage: "set the encrypted key regex. When specified, only keys matching the regex will be encrypted.",
		},
		cli.StringFlag{
			Name:  "unencrypted-comment-regex",
			Usage: "set the unencrypted comment suffix. When specified, only keys that have comment matching the regex will be left unencrypted.",
		},
		cli.StringFlag{
			Name:  "encrypted-comment-regex",
			Usage: "set the encrypted comment suffix. When specified, only keys that have comment matching the regex will be encrypted.",
		},
		cli.StringFlag{
			Name:   "config",
			Usage:  "path to sops' config file. If set, sops will not search for the config file recursively.",
			EnvVar: "SOPS_CONFIG",
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
		cli.IntFlag{
			Name:  "indent",
			Usage: "the number of spaces to indent YAML or JSON encoded file",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logging output",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "Save the output after encryption or decryption to the file specified",
		},
		cli.StringFlag{
			Name:  "filename-override",
			Usage: "Use this filename instead of the provided argument for loading configuration, and for determining input type and output type",
		},
		cli.StringFlag{
			Name:   "decryption-order",
			Usage:  "comma separated list of decryption key types",
			EnvVar: "SOPS_DECRYPTION_ORDER",
		},
	}, keyserviceFlags...)

	app.Action = func(c *cli.Context) error {
		isDecryptMode := c.Bool("decrypt")
		isEncryptMode := c.Bool("encrypt")
		isRotateMode := c.Bool("rotate")
		isSetMode := c.String("set") != ""
		isEditMode := !isEncryptMode && !isDecryptMode && !isRotateMode && !isSetMode

		if c.Bool("verbose") {
			logging.SetLevel(logrus.DebugLevel)
		}
		if c.NArg() < 1 {
			return common.NewExitError("Error: no file specified", codes.NoFileSpecified)
		}
		warnMoreThanOnePositionalArgument(c)
		if c.Bool("in-place") && c.String("output") != "" {
			return common.NewExitError("Error: cannot operate on both --output and --in-place", codes.ErrorConflictingParameters)
		}
		fileName, err := filepath.Abs(c.Args()[0])
		if err != nil {
			return toExitError(err)
		}
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if c.String("add-kms") != "" || c.String("add-pgp") != "" || c.String("add-gcp-kms") != "" || c.String("add-hc-vault-transit") != "" || c.String("add-azure-kv") != "" || c.String("add-age") != "" ||
				c.String("rm-kms") != "" || c.String("rm-pgp") != "" || c.String("rm-gcp-kms") != "" || c.String("rm-hc-vault-transit") != "" || c.String("rm-azure-kv") != "" || c.String("rm-age") != "" {
				return common.NewExitError(fmt.Sprintf("Error: cannot add or remove keys on non-existent file %q, use `--kms` and `--pgp` instead.", fileName), codes.CannotChangeKeysFromNonExistentFile)
			}
			if isEncryptMode || isDecryptMode || isRotateMode {
				return common.NewExitError(fmt.Sprintf("Error: cannot operate on non-existent file %q", fileName), codes.NoFileSpecified)
			}
		}
		fileNameOverride := c.String("filename-override")
		if fileNameOverride == "" {
			fileNameOverride = fileName
		} else {
			fileNameOverride, err = filepath.Abs(fileNameOverride)
			if err != nil {
				return toExitError(err)
			}
		}

		commandCount := 0
		if isDecryptMode {
			commandCount++
		}
		if isEncryptMode {
			commandCount++
		}
		if isRotateMode {
			commandCount++
		}
		if isSetMode {
			commandCount++
		}
		if commandCount > 1 {
			log.Warn("More than one command (--encrypt, --decrypt, --rotate, --set) has been specified. Only the changes made by the last one will be visible. Note that this behavior is deprecated and will cause an error eventually.")
		}

		// Load configuration here for backwards compatibility (error out in case of bad config files),
		// but only when not just decrypting (https://github.com/getsops/sops/issues/868)
		needsCreationRule := isEncryptMode || isRotateMode || isSetMode || isEditMode
		if needsCreationRule {
			_, err = loadConfig(c, fileNameOverride, nil)
			if err != nil {
				return toExitError(err)
			}
		}

		inputStore, err := inputStore(c, fileNameOverride)
		if err != nil {
			return toExitError(err)
		}
		outputStore, err := outputStore(c, fileNameOverride)
		if err != nil {
			return toExitError(err)
		}
		svcs := keyservices(c)

		order, err := decryptionOrder(c.String("decryption-order"))
		if err != nil {
			return toExitError(err)
		}
		var output []byte
		if isEncryptMode {
			encConfig, err := getEncryptConfig(c, fileNameOverride)
			if err != nil {
				return toExitError(err)
			}
			output, err = encrypt(encryptOpts{
				OutputStore:   outputStore,
				InputStore:    inputStore,
				InputPath:     fileName,
				Cipher:        aes.NewCipher(),
				KeyServices:   svcs,
				encryptConfig: encConfig,
			})
			// While this check is also done below, the `err` in this scope shadows
			// the `err` in the outer scope.  **Only** do this in case --decrypt,
			// --rotate-, and --set are not specified, though, to keep old behavior.
			if err != nil && !isDecryptMode && !isRotateMode && !isSetMode {
				return toExitError(err)
			}
		}

		if isDecryptMode {
			var extract []interface{}
			extract, err = parseTreePath(c.String("extract"))
			if err != nil {
				return common.NewExitError(fmt.Errorf("error parsing --extract path: %s", err), codes.InvalidTreePathFormat)
			}
			output, err = decrypt(decryptOpts{
				OutputStore:     outputStore,
				InputStore:      inputStore,
				InputPath:       fileName,
				Cipher:          aes.NewCipher(),
				Extract:         extract,
				KeyServices:     svcs,
				DecryptionOrder: order,
				IgnoreMAC:       c.Bool("ignore-mac"),
			})
		}
		if isRotateMode {
			rotateOpts, err := getRotateOpts(c, fileName, inputStore, outputStore, svcs, order)
			if err != nil {
				return toExitError(err)
			}

			output, err = rotate(rotateOpts)
			// While this check is also done below, the `err` in this scope shadows
			// the `err` in the outer scope
			if err != nil {
				return toExitError(err)
			}
		}

		if isSetMode {
			var path []interface{}
			var value interface{}
			path, value, err = extractSetArguments(c.String("set"))
			if err != nil {
				return toExitError(err)
			}
			output, _, err = set(setOpts{
				OutputStore:     outputStore,
				InputStore:      inputStore,
				InputPath:       fileName,
				Cipher:          aes.NewCipher(),
				KeyServices:     svcs,
				DecryptionOrder: order,
				IgnoreMAC:       c.Bool("ignore-mac"),
				Value:           value,
				TreePath:        path,
			})
		}

		if isEditMode {
			_, statErr := os.Stat(fileName)
			fileExists := statErr == nil
			opts := editOpts{
				OutputStore:     outputStore,
				InputStore:      inputStore,
				InputPath:       fileName,
				Cipher:          aes.NewCipher(),
				KeyServices:     svcs,
				DecryptionOrder: order,
				IgnoreMAC:       c.Bool("ignore-mac"),
				ShowMasterKeys:  c.Bool("show-master-keys"),
			}
			if fileExists {
				output, err = edit(opts)
			} else {
				// File doesn't exist, edit the example file instead
				encConfig, err := getEncryptConfig(c, fileNameOverride)
				if err != nil {
					return toExitError(err)
				}
				output, err = editExample(editExampleOpts{
					editOpts:      opts,
					encryptConfig: encConfig,
				})
				// While this check is also done below, the `err` in this scope shadows
				// the `err` in the outer scope
				if err != nil {
					return toExitError(err)
				}
			}
		}

		if err != nil {
			return toExitError(err)
		}

		// We open the file *after* the operations on the tree have been
		// executed to avoid truncating it when there's errors
		if c.Bool("in-place") || isEditMode || isSetMode {
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
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getEncryptConfig(c *cli.Context, fileName string) (encryptConfig, error) {
	unencryptedSuffix := c.String("unencrypted-suffix")
	encryptedSuffix := c.String("encrypted-suffix")
	encryptedRegex := c.String("encrypted-regex")
	unencryptedRegex := c.String("unencrypted-regex")
	encryptedCommentRegex := c.String("encrypted-comment-regex")
	unencryptedCommentRegex := c.String("unencrypted-comment-regex")
	macOnlyEncrypted := c.Bool("mac-only-encrypted")
	conf, err := loadConfig(c, fileName, nil)
	if err != nil {
		return encryptConfig{}, toExitError(err)
	}
	if conf != nil {
		// command line options have precedence
		if unencryptedSuffix == "" {
			unencryptedSuffix = conf.UnencryptedSuffix
		}
		if encryptedSuffix == "" {
			encryptedSuffix = conf.EncryptedSuffix
		}
		if encryptedRegex == "" {
			encryptedRegex = conf.EncryptedRegex
		}
		if unencryptedRegex == "" {
			unencryptedRegex = conf.UnencryptedRegex
		}
		if encryptedCommentRegex == "" {
			encryptedCommentRegex = conf.EncryptedCommentRegex
		}
		if unencryptedCommentRegex == "" {
			unencryptedCommentRegex = conf.UnencryptedCommentRegex
		}
		if !macOnlyEncrypted {
			macOnlyEncrypted = conf.MACOnlyEncrypted
		}
	}

	cryptRuleCount := 0
	if unencryptedSuffix != "" {
		cryptRuleCount++
	}
	if encryptedSuffix != "" {
		cryptRuleCount++
	}
	if encryptedRegex != "" {
		cryptRuleCount++
	}
	if unencryptedRegex != "" {
		cryptRuleCount++
	}
	if encryptedCommentRegex != "" {
		cryptRuleCount++
	}
	if unencryptedCommentRegex != "" {
		cryptRuleCount++
	}

	if cryptRuleCount > 1 {
		return encryptConfig{}, common.NewExitError("Error: cannot use more than one of encrypted_suffix, unencrypted_suffix, encrypted_regex, unencrypted_regex, encrypted_comment_regex, or unencrypted_comment_regex in the same file", codes.ErrorConflictingParameters)
	}

	// only supply the default UnencryptedSuffix when EncryptedSuffix, EncryptedRegex, and others are not provided
	if cryptRuleCount == 0 {
		unencryptedSuffix = sops.DefaultUnencryptedSuffix
	}

	var groups []sops.KeyGroup
	groups, err = keyGroups(c, fileName)
	if err != nil {
		return encryptConfig{}, err
	}

	var threshold int
	threshold, err = shamirThreshold(c, fileName)
	if err != nil {
		return encryptConfig{}, err
	}

	return encryptConfig{
		UnencryptedSuffix:       unencryptedSuffix,
		EncryptedSuffix:         encryptedSuffix,
		UnencryptedRegex:        unencryptedRegex,
		EncryptedRegex:          encryptedRegex,
		UnencryptedCommentRegex: unencryptedCommentRegex,
		EncryptedCommentRegex:   encryptedCommentRegex,
		MACOnlyEncrypted:        macOnlyEncrypted,
		KeyGroups:               groups,
		GroupThreshold:          threshold,
	}, nil
}

func getMasterKeys(c *cli.Context, kmsEncryptionContext map[string]*string, kmsOptionName string, pgpOptionName string, gcpKmsOptionName string, azureKvOptionName string, hcVaultTransitOptionName string, ageOptionName string) ([]keys.MasterKey, error) {
	var masterKeys []keys.MasterKey
	for _, k := range kms.MasterKeysFromArnString(c.String(kmsOptionName), kmsEncryptionContext, c.String("aws-profile")) {
		masterKeys = append(masterKeys, k)
	}
	for _, k := range pgp.MasterKeysFromFingerprintString(c.String(pgpOptionName)) {
		masterKeys = append(masterKeys, k)
	}
	for _, k := range gcpkms.MasterKeysFromResourceIDString(c.String(gcpKmsOptionName)) {
		masterKeys = append(masterKeys, k)
	}
	azureKeys, err := azkv.MasterKeysFromURLs(c.String(azureKvOptionName))
	if err != nil {
		return nil, err
	}
	for _, k := range azureKeys {
		masterKeys = append(masterKeys, k)
	}
	hcVaultKeys, err := hcvault.NewMasterKeysFromURIs(c.String(hcVaultTransitOptionName))
	if err != nil {
		return nil, err
	}
	for _, k := range hcVaultKeys {
		masterKeys = append(masterKeys, k)
	}
	ageKeys, err := age.MasterKeysFromRecipients(c.String(ageOptionName))
	if err != nil {
		return nil, err
	}
	for _, k := range ageKeys {
		masterKeys = append(masterKeys, k)
	}
	return masterKeys, nil
}

func getRotateOpts(c *cli.Context, fileName string, inputStore common.Store, outputStore common.Store, svcs []keyservice.KeyServiceClient, decryptionOrder []string) (rotateOpts, error) {
	kmsEncryptionContext := kms.ParseKMSContext(c.String("encryption-context"))
	addMasterKeys, err := getMasterKeys(c, kmsEncryptionContext, "add-kms", "add-pgp", "add-gcp-kms", "add-azure-kv", "add-hc-vault-transit", "add-age")
	if err != nil {
		return rotateOpts{}, err
	}
	rmMasterKeys, err := getMasterKeys(c, kmsEncryptionContext, "rm-kms", "rm-pgp", "rm-gcp-kms", "rm-azure-kv", "rm-hc-vault-transit", "rm-age")
	if err != nil {
		return rotateOpts{}, err
	}
	return rotateOpts{
		OutputStore:      outputStore,
		InputStore:       inputStore,
		InputPath:        fileName,
		Cipher:           aes.NewCipher(),
		KeyServices:      svcs,
		DecryptionOrder:  decryptionOrder,
		IgnoreMAC:        c.Bool("ignore-mac"),
		AddMasterKeys:    addMasterKeys,
		RemoveMasterKeys: rmMasterKeys,
	}, nil
}

func toExitError(err error) error {
	if cliErr, ok := err.(*cli.ExitError); ok && cliErr != nil {
		return cliErr
	} else if execErr, ok := err.(*osExec.ExitError); ok && execErr != nil {
		return cli.NewExitError(err, execErr.ExitCode())
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
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(
				func(ctx context.Context, addr string) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, url.Scheme, addr)
				},
			),
		}
		log.WithField(
			"address",
			fmt.Sprintf("%s://%s", url.Scheme, addr),
		).Infof("Connecting to key service")
		conn, err := grpc.NewClient(addr, opts...)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		svcs = append(svcs, keyservice.NewKeyServiceClient(conn))
	}
	return
}

// Wrapper of config.LookupConfigFile that takes care of handling the returned warning.
func findConfigFile() (string, error) {
	result, err := config.LookupConfigFile(".")
	if len(result.Warning) > 0 && !showedConfigFileWarning {
		showedConfigFileWarning = true
		log.Warn(result.Warning)
	}
	return result.Path, err
}

func loadStoresConfig(context *cli.Context, path string) (*config.StoresConfig, error) {
	configPath := context.GlobalString("config")
	if configPath == "" {
		// Ignore config not found errors returned from findConfigFile since the config file is not mandatory
		foundPath, err := findConfigFile()
		if err != nil {
			return config.NewStoresConfig(), nil
		}
		configPath = foundPath
	}
	return config.LoadStoresConfig(configPath)
}

func inputStore(context *cli.Context, path string) (common.Store, error) {
	storesConf, err := loadStoresConfig(context, path)
	if err != nil {
		return nil, err
	}
	return common.DefaultStoreForPathOrFormat(storesConf, path, context.String("input-type")), nil
}

func outputStore(context *cli.Context, path string) (common.Store, error) {
	storesConf, err := loadStoresConfig(context, path)
	if err != nil {
		return nil, err
	}
	if context.IsSet("indent") {
		indent := context.Int("indent")
		storesConf.YAML.Indent = indent
		storesConf.JSON.Indent = indent
		storesConf.JSONBinary.Indent = indent
	}

	return common.DefaultStoreForPathOrFormat(storesConf, path, context.String("output-type")), nil
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
	var hcVaultMkKeys []keys.MasterKey
	var ageMasterKeys []keys.MasterKey
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
	if c.String("hc-vault-transit") != "" {
		hcVaultKeys, err := hcvault.NewMasterKeysFromURIs(c.String("hc-vault-transit"))
		if err != nil {
			return nil, err
		}
		for _, k := range hcVaultKeys {
			hcVaultMkKeys = append(hcVaultMkKeys, k)
		}
	}
	if c.String("pgp") != "" {
		for _, k := range pgp.MasterKeysFromFingerprintString(c.String("pgp")) {
			pgpKeys = append(pgpKeys, k)
		}
	}
	if c.String("age") != "" {
		ageKeys, err := age.MasterKeysFromRecipients(c.String("age"))
		if err != nil {
			return nil, err
		}
		for _, k := range ageKeys {
			ageMasterKeys = append(ageMasterKeys, k)
		}
	}
	if c.String("kms") == "" && c.String("pgp") == "" && c.String("gcp-kms") == "" && c.String("azure-kv") == "" && c.String("hc-vault-transit") == "" && c.String("age") == "" {
		conf, err := loadConfig(c, file, kmsEncryptionContext)
		// config file might just not be supplied, without any error
		if conf == nil {
			errMsg := "config file not found, or has no creation rules, and no keys provided through command line options"
			if err != nil {
				errMsg = fmt.Sprintf("%s: %s", errMsg, err)
			}
			return nil, fmt.Errorf("%s", errMsg)
		}
		return conf.KeyGroups, err
	}
	var group sops.KeyGroup
	group = append(group, kmsKeys...)
	group = append(group, cloudKmsKeys...)
	group = append(group, azkvKeys...)
	group = append(group, pgpKeys...)
	group = append(group, hcVaultMkKeys...)
	group = append(group, ageMasterKeys...)
	log.Debugf("Master keys available:  %+v", group)
	return []sops.KeyGroup{group}, nil
}

// loadConfig will look for an existing config file, either provided through the command line, or using findConfigFile
// Since a config file is not required, this function does not error when one is not found, and instead returns a nil config pointer
func loadConfig(c *cli.Context, file string, kmsEncryptionContext map[string]*string) (*config.Config, error) {
	var err error
	configPath := c.GlobalString("config")
	if configPath == "" {
		// Ignore config not found errors returned from findConfigFile since the config file is not mandatory
		configPath, err = findConfigFile()
		if err != nil {
			// If we can't find a config file, but we were not explicitly requested to, assume it does not exist
			return nil, nil
		}
	}
	conf, err := config.LoadCreationRuleForFile(configPath, file, kmsEncryptionContext)
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
		// 1. No config was provided, or contains no creation rules. Err will be nil and ShamirThreshold will be the default value of 0.
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
	// Fix for #461
	// Attempt conversion to TreeBranches to handle yaml multidoc. If conversion fails it's
	// most likely a string value, so just return it as-is.
	values, ok := valueToInsert.(sops.TreeBranches)
	if !ok {
		return valueToInsert, nil
	}
	return values[0], nil
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
	if err != nil {
		// All errors returned by jsonValueToTreeInsertableValue are created by common.NewExitError(),
		// so we can simply pass them on
		return nil, nil, err
	}

	path, err = parseTreePath(fullPath)
	if err != nil {
		return nil, nil, common.NewExitError("Invalid --set format", codes.ErrorInvalidSetFormat)
	}
	return path, valueToInsert, nil
}

func decryptionOrder(decryptionOrder string) ([]string, error) {
	if decryptionOrder == "" {
		return sops.DefaultDecryptionOrder, nil
	}
	orderList := strings.Split(decryptionOrder, ",")
	unique := make(map[string]struct{})
	for _, v := range orderList {
		if _, ok := unique[v]; ok {
			return nil, common.NewExitError(fmt.Sprintf("Duplicate decryption key type: %s", v), codes.DuplicateDecryptionKeyType)
		}
		unique[v] = struct{}{}
	}
	return orderList, nil
}
