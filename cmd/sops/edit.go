package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/version"
	"github.com/google/shlex"
	exec "golang.org/x/sys/execabs"
)

type editOpts struct {
	Cipher          sops.Cipher
	InputStore      common.Store
	OutputStore     common.Store
	InputPath       string
	IgnoreMAC       bool
	KeyServices     []keyservice.KeyServiceClient
	DecryptionOrder []string
	ShowMasterKeys  bool
}

type editExampleOpts struct {
	editOpts
	encryptConfig
}

type runEditorUntilOkOpts struct {
	TmpFileName    string
	OriginalHash   []byte
	InputStore     sops.Store
	ShowMasterKeys bool
	Tree           *sops.Tree
}

func editExample(opts editExampleOpts) ([]byte, error) {
	fileBytes := opts.InputStore.EmitExample()
	branches, err := opts.InputStore.LoadPlainFile(fileBytes)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	path, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return nil, err
	}
	tree := sops.Tree{
		Branches: branches,
		Metadata: metadataFromEncryptionConfig(opts.encryptConfig),
		FilePath: path,
	}

	// Generate a data key
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		return nil, common.NewExitError(fmt.Sprintf("Error encrypting the data key with one or more master keys: %s", errs), codes.CouldNotRetrieveKey)
	}

	return editTree(opts.editOpts, &tree, dataKey)
}

func edit(opts editOpts) ([]byte, error) {
	// Load the file
	tree, err := common.LoadEncryptedFileWithBugFixes(common.GenericDecryptOpts{
		Cipher:      opts.Cipher,
		InputStore:  opts.InputStore,
		InputPath:   opts.InputPath,
		IgnoreMAC:   opts.IgnoreMAC,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}
	// Decrypt the file
	dataKey, err := common.DecryptTree(common.DecryptTreeOpts{
		Cipher:          opts.Cipher,
		IgnoreMac:       opts.IgnoreMAC,
		Tree:            tree,
		KeyServices:     opts.KeyServices,
		DecryptionOrder: opts.DecryptionOrder,
	})
	if err != nil {
		return nil, err
	}

	return editTree(opts, tree, dataKey)
}

func editTree(opts editOpts, tree *sops.Tree, dataKey []byte) ([]byte, error) {
	// Create temporary file for editing
	tmpdir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not create temporary directory: %s", err), codes.CouldNotWriteOutputFile)
	}
	defer os.RemoveAll(tmpdir)

	tmpfile, err := os.Create(filepath.Join(tmpdir, filepath.Base(opts.InputPath)))
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not create temporary file: %s", err), codes.CouldNotWriteOutputFile)
	}
	// Ensure that in any case, the temporary file is always closed.
	defer tmpfile.Close()

	tmpfileName := tmpfile.Name()

	// Write to temporary file
	var out []byte
	if opts.ShowMasterKeys {
		out, err = opts.OutputStore.EmitEncryptedFile(*tree)
	} else {
		out, err = opts.OutputStore.EmitPlainFile(tree.Branches)
	}
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	_, err = tmpfile.Write(out)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not write output file: %s", err), codes.CouldNotWriteOutputFile)
	}

	// Compute file hash to detect if the file has been edited
	origHash, err := hashFile(tmpfileName)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not hash file: %s", err), codes.CouldNotReadInputFile)
	}

	// Close the temporary file, so that an editor can open it.
	// We need to do this because some editors (e.g. VSCode) will refuse to
	// open a file on Windows due to the Go standard library not opening
	// files with shared delete access.
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	// Let the user edit the file
	err = runEditorUntilOk(runEditorUntilOkOpts{
		InputStore: opts.InputStore, OriginalHash: origHash, TmpFileName: tmpfileName,
		ShowMasterKeys: opts.ShowMasterKeys, Tree: tree})
	if err != nil {
		return nil, err
	}

	// Encrypt the file
	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey, Tree: tree, Cipher: opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	// Output the file
	encryptedFile, err := opts.OutputStore.EmitEncryptedFile(*tree)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return encryptedFile, nil
}

func runEditorUntilOk(opts runEditorUntilOkOpts) error {
	for {
		err := runEditor(opts.TmpFileName)
		if err != nil {
			return common.NewExitError(fmt.Sprintf("Could not run editor: %s", err), codes.NoEditorFound)
		}
		newHash, err := hashFile(opts.TmpFileName)
		if err != nil {
			return common.NewExitError(fmt.Sprintf("Could not hash file: %s", err), codes.CouldNotReadInputFile)
		}
		if bytes.Equal(newHash, opts.OriginalHash) {
			return common.NewExitError("File has not changed, exiting.", codes.FileHasNotBeenModified)
		}
		edited, err := os.ReadFile(opts.TmpFileName)
		if err != nil {
			return common.NewExitError(fmt.Sprintf("Could not read edited file: %s", err), codes.CouldNotReadInputFile)
		}
		newBranches, err := opts.InputStore.LoadPlainFile(edited)
		if err != nil {
			log.WithField(
				"error",
				err,
			).Errorf("Could not load tree, probably due to invalid " +
				"syntax. Press a key to return to the editor, or Ctrl+C to " +
				"exit.")
			bufio.NewReader(os.Stdin).ReadByte()
			continue
		}
		if opts.ShowMasterKeys {
			// The file is not actually encrypted, but it contains SOPS
			// metadata
			t, err := opts.InputStore.LoadEncryptedFile(edited)
			if err != nil {
				log.WithField(
					"error",
					err,
				).Errorf("SOPS metadata is invalid. Press a key to " +
					"return to the editor, or Ctrl+C to exit.")
				bufio.NewReader(os.Stdin).ReadByte()
				continue
			}
			// Replace the whole tree, because otherwise newBranches would
			// contain the SOPS metadata
			opts.Tree = &t
		}
		opts.Tree.Branches = newBranches
		needVersionUpdated, err := version.AIsNewerThanB(version.Version, opts.Tree.Metadata.Version)
		if err != nil {
			return common.NewExitError(fmt.Sprintf("Failed to compare document version %q with program version %q: %v", opts.Tree.Metadata.Version, version.Version, err), codes.FailedToCompareVersions)
		}
		if needVersionUpdated {
			opts.Tree.Metadata.Version = version.Version
		}
		if opts.Tree.Metadata.MasterKeyCount() == 0 {
			log.Error("No master keys were provided, so sops can't " +
				"encrypt the file. Press a key to return to the editor, or " +
				"Ctrl+C to exit.")
			bufio.NewReader(os.Stdin).ReadByte()
			continue
		}
		break
	}
	return nil
}

func hashFile(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}

func runEditor(path string) error {
	envVar := "SOPS_EDITOR"
	editor := os.Getenv(envVar)
	if editor == "" {
		envVar = "EDITOR"
		editor = os.Getenv(envVar)
	}
	var cmd *exec.Cmd
	if editor == "" {
		editor, err := lookupAnyEditor("vim", "nano", "vi")
		if err != nil {
			return err
		}
		cmd = exec.Command(editor, path)
	} else {
		parts, err := shlex.Split(editor)
		if err != nil {
			return fmt.Errorf("invalid $%s: %s", envVar, editor)
		}
		parts = append(parts, path)
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func lookupAnyEditor(editorNames ...string) (editorPath string, err error) {
	for _, editorName := range editorNames {
		editorPath, err = exec.LookPath(editorName)
		if err == nil {
			return editorPath, nil
		}
	}
	return "", fmt.Errorf("no editor available: sops attempts to use the editor defined in the SOPS_EDITOR or EDITOR environment variables, and if that's not set defaults to any of %s, but none of them could be found", strings.Join(editorNames, ", "))
}
