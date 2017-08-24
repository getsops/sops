package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"crypto/md5"
	"io"
	"os/exec"
	"strings"

	"time"

	"github.com/google/shlex"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/stores/json"
	"gopkg.in/urfave/cli.v1"
)

type EditOpts struct {
	Cipher         sops.DataKeyCipher
	InputStore     sops.Store
	OutputStore    sops.Store
	InputPath      string
	IgnoreMAC      bool
	KeyServices    []keyservice.KeyServiceClient
	ShowMasterKeys bool
}

type EditExampleOpts struct {
	Cipher            sops.DataKeyCipher
	InputStore        sops.Store
	OutputStore       sops.Store
	InputPath         string
	IgnoreMAC         bool
	KeyServices       []keyservice.KeyServiceClient
	ShowMasterKeys    bool
	UnencryptedSuffix string
	KeyGroups         []sops.KeyGroup
	GroupQuorum       uint
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

func EditExample(opts EditExampleOpts) ([]byte, error) {
	// Load the example file
	var fileBytes []byte
	if _, ok := opts.InputStore.(*json.BinaryStore); ok {
		// Get the value under the first key
		fileBytes = []byte(exampleTree[0].Value.(string))
	} else {
		var err error
		fileBytes, err = opts.InputStore.Marshal(exampleTree)
		if err != nil {
			return nil, err
		}
	}
	var tree sops.Tree
	branch, err := opts.InputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), exitCouldNotReadInputFile)
	}
	tree.Branch = branch
	tree.Metadata = sops.Metadata{
		KeyGroups:         opts.KeyGroups,
		UnencryptedSuffix: opts.UnencryptedSuffix,
		Version:           version,
		ShamirQuorum:      int(opts.GroupQuorum),
	}

	// Generate a data key
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		return nil, cli.NewExitError(fmt.Sprintf("Error encrypting the data key with one or more master keys: %s", errs), exitCouldNotRetrieveKey)
	}

	// Create temporary file for editing
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not create temporary directory: %s", err), exitCouldNotWriteOutputFile)
	}
	defer os.RemoveAll(tmpdir)
	tmpfile, err := os.Create(path.Join(tmpdir, path.Base(opts.InputPath)))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not create temporary file: %s", err), exitCouldNotWriteOutputFile)
	}

	// Write to temporary file
	var out []byte
	if opts.ShowMasterKeys {
		out, err = opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	} else {
		out, err = opts.OutputStore.Marshal(tree.Branch)
	}
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	_, err = tmpfile.Write(out)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not write output file: %s", err), exitCouldNotWriteOutputFile)
	}

	// Compute file hash to detect if the file has been edited
	origHash, err := hashFile(tmpfile.Name())
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not hash file: %s", err), exitCouldNotReadInputFile)
	}

	// Let the user edit the file
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
		newBranch, err := opts.InputStore.Unmarshal(edited)
		if err != nil {
			fmt.Printf("Could not load tree: %s\nProbably invalid syntax. Press a key to return to the editor, or Ctrl+C to exit.", err)
			bufio.NewReader(os.Stdin).ReadByte()
			continue
		}
		if opts.ShowMasterKeys {
			metadata, err := opts.InputStore.UnmarshalMetadata(edited)
			if err != nil {
				fmt.Printf("sops branch is invalid: %s.\nPress a key to return to the editor, or Ctrl+C to exit.", err)
				bufio.NewReader(os.Stdin).ReadByte()
				continue
			}
			tree.Metadata = *metadata
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

	// Encrypt the file
	mac, err := tree.Encrypt(dataKey, opts.Cipher, make(map[string][]interface{}))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	tree.Metadata.LastModified = time.Now().UTC()
	mac, err = opts.Cipher.Encrypt(mac, dataKey, tree.Metadata.LastModified.Format(time.RFC3339), nil)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), exitErrorEncryptingMac)
	}
	tree.Metadata.MessageAuthenticationCode = mac

	// Output the file
	encryptedFile, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return encryptedFile, nil
}

func Edit(opts EditOpts) ([]byte, error) {
	// Load the file
	fileBytes, err := ioutil.ReadFile(opts.InputPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error reading file: %s", err), exitCouldNotReadInputFile)
	}
	metadata, err := opts.InputStore.UnmarshalMetadata(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file metadata: %s", err), exitCouldNotReadInputFile)
	}
	branch, err := opts.InputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree := sops.Tree{
		Branch:   branch,
		Metadata: *metadata,
	}
	// Decrypt the file
	dataKey, err := tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return nil, cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	computedMac, err := tree.Decrypt(dataKey, opts.Cipher, make(map[string][]interface{}))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), exitErrorDecryptingTree)
	}
	fileMac, _, err := opts.Cipher.Decrypt(tree.Metadata.MessageAuthenticationCode, dataKey, tree.Metadata.LastModified.Format(time.RFC3339))
	if !opts.IgnoreMAC {
		if fileMac != computedMac {
			// If the file has an empty MAC, display "no MAC" instead of not displaying anything
			if fileMac == "" {
				fileMac = "no MAC"
			}
			return nil, cli.NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", fileMac, computedMac), exitMacMismatch)
		}
	}
	// Create temporary file for editing
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not create temporary directory: %s", err), exitCouldNotWriteOutputFile)
	}
	defer os.RemoveAll(tmpdir)
	tmpfile, err := os.Create(path.Join(tmpdir, path.Base(opts.InputPath)))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not create temporary file: %s", err), exitCouldNotWriteOutputFile)
	}

	// Write to temporary file
	var out []byte
	if opts.ShowMasterKeys {
		out, err = opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	} else {
		out, err = opts.OutputStore.Marshal(tree.Branch)
	}
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	_, err = tmpfile.Write(out)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not write output file: %s", err), exitCouldNotWriteOutputFile)
	}

	// Compute file hash to detect if the file has been edited
	origHash, err := hashFile(tmpfile.Name())
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not hash file: %s", err), exitCouldNotReadInputFile)
	}

	// Let the user edit the file
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
		newBranch, err := opts.InputStore.Unmarshal(edited)
		if err != nil {
			fmt.Printf("Could not load tree: %s\nProbably invalid syntax. Press a key to return to the editor, or Ctrl+C to exit.", err)
			bufio.NewReader(os.Stdin).ReadByte()
			continue
		}
		if opts.ShowMasterKeys {
			metadata, err := opts.InputStore.UnmarshalMetadata(edited)
			if err != nil {
				fmt.Printf("sops branch is invalid: %s.\nPress a key to return to the editor, or Ctrl+C to exit.", err)
				bufio.NewReader(os.Stdin).ReadByte()
				continue
			}
			tree.Metadata = *metadata
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

	// Encrypt the file
	mac, err := tree.Encrypt(dataKey, opts.Cipher, make(map[string][]interface{}))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	tree.Metadata.LastModified = time.Now().UTC()
	mac, err = opts.Cipher.Encrypt(mac, dataKey, tree.Metadata.LastModified.Format(time.RFC3339), nil)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), exitErrorEncryptingMac)
	}
	tree.Metadata.MessageAuthenticationCode = mac

	// Output the file
	encryptedFile, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return encryptedFile, nil
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
		parts, err := shlex.Split(editor)
		if err != nil {
			return fmt.Errorf("Invalid $EDITOR: %s", editor)
		}
		parts = append(parts, path)
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
