package sops

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fatih/color"
	"github.com/goware/prefixer"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

// UserError is a well-formatted error for the purpose of being displayed to
// the end user.
type UserError interface {
	error
	UserError() string
}

var statusSuccess = color.New(color.FgGreen).Sprint("SUCCESS")
var statusFailed = color.New(color.FgRed).Sprint("FAILED")

type getDataKeyError struct {
	RequiredSuccessfulKeyGroups int
	GroupResults                []error
}

func (err *getDataKeyError) successfulKeyGroups() int {
	n := 0
	for _, r := range err.GroupResults {
		if r == nil {
			n++
		}
	}
	return n
}

func (err *getDataKeyError) Error() string {
	return fmt.Sprintf("Error getting data key: %d successful groups "+
		"required, got %d", err.RequiredSuccessfulKeyGroups,
		err.successfulKeyGroups())
}

func (err *getDataKeyError) UserError() string {
	var groupErrs []string
	for i, res := range err.GroupResults {
		groupErr := decryptGroupError{
			err:       res,
			groupName: fmt.Sprintf("%d", i),
		}
		groupErrs = append(groupErrs, groupErr.UserError())
	}
	var trailer string
	if err.RequiredSuccessfulKeyGroups == 0 {
		trailer = "Recovery failed because no master key was able to decrypt " +
			"the file. In order for SOPS to recover the file, at least one key " +
			"has to be successful, but none were."
	} else {
		trailer = fmt.Sprintf("Recovery failed because the file was "+
			"encrypted with a Shamir threshold of %d, but only %d part(s) "+
			"were successfully recovered, one for each successful key group. "+
			"In order for SOPS to recover the file, at least %d groups have "+
			"to be successful. In order for a group to be successful, "+
			"decryption has to succeed with any of the keys in that key group.",
			err.RequiredSuccessfulKeyGroups, err.successfulKeyGroups(),
			err.RequiredSuccessfulKeyGroups)
	}
	trailer = wordwrap.WrapString(trailer, 75)
	return fmt.Sprintf("Failed to get the data key required to "+
		"decrypt the SOPS file.\n\n%s\n\n%s",
		strings.Join(groupErrs, "\n\n"), trailer)
}

type decryptGroupError struct {
	groupName string
	err       error
}

func (r *decryptGroupError) Error() string {
	return fmt.Sprintf("could not decrypt group %s: %s", r.groupName, r.err)
}

func (r *decryptGroupError) UserError() string {
	var status string
	if r.err == nil {
		status = statusSuccess
	} else {
		status = statusFailed
	}
	header := fmt.Sprintf(`Group %s: %s`, r.groupName, status)
	if r.err == nil {
		return header
	}
	message := r.err.Error()
	if userError, ok := r.err.(UserError); ok {
		message = userError.UserError()
	}
	reader := prefixer.New(strings.NewReader(message), "  ")
	// Safe to ignore this error, as reading from a strings.Reader can't fail
	errMsg, _ := ioutil.ReadAll(reader)
	return fmt.Sprintf("%s\n%s", header, string(errMsg))
}

type decryptKeyErrors []error

func (e decryptKeyErrors) Error() string {
	return fmt.Sprintf("error decrypting key: %s", []error(e))
}

func (e decryptKeyErrors) UserError() string {
	var errStrs []string
	for _, err := range []error(e) {
		if userErr, ok := err.(UserError); ok {
			errStrs = append(errStrs, userErr.UserError())
		} else {
			errStrs = append(errStrs, err.Error())
		}
	}
	return strings.Join(errStrs, "\n\n")
}

type decryptKeyError struct {
	keyName string
	errs    []error
}

func (e *decryptKeyError) isSuccessful() bool {
	for _, err := range e.errs {
		if err == nil {
			return true
		}
	}
	return false
}

func (e *decryptKeyError) Error() string {
	return fmt.Sprintf("error decrypting key %s: %s", e.keyName, e.errs)
}

func (e *decryptKeyError) UserError() string {
	var status string
	if e.isSuccessful() {
		status = statusSuccess
	} else {
		status = statusFailed
	}
	header := fmt.Sprintf("%s: %s", e.keyName, status)
	if e.isSuccessful() {
		return header
	}
	var errMessages []string
	for _, err := range e.errs {
		wrappedErr := wordwrap.WrapString(err.Error(), 60)
		reader := prefixer.New(strings.NewReader(wrappedErr), "  | ")
		// Safe to ignore this error, as reading from a strings.Reader can't fail
		errMsg, _ := ioutil.ReadAll(reader)
		errMsg[0] = '-'
		errMessages = append(errMessages, string(errMsg))
	}
	joinedMsgs := strings.Join(errMessages, "\n\n")
	reader := prefixer.New(strings.NewReader(joinedMsgs), "  ")
	errMsg, _ := ioutil.ReadAll(reader)
	return fmt.Sprintf("%s\n%s", header, string(errMsg))
}
