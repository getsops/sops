/*
Package sops manages JSON, YAML and BINARY documents to be encrypted or decrypted.

This package should not be used directly. Instead, Sops users should install the
command line client via `go get -u github.com/getsops/sops/v3/cmd/sops`, or use the
decryption helper provided at `github.com/getsops/sops/v3/decrypt`.

We do not guarantee API stability for any package other than `github.com/getsops/sops/v3/decrypt`.

A Sops document is a Tree composed of a data branch with arbitrary key/value pairs
and a metadata branch with encryption and integrity information.

In JSON and YAML formats, the structure of the cleartext tree is preserved, keys are
stored in cleartext and only values are encrypted. Keeping the keys in cleartext
provides better readability when storing Sops documents in version controls, and allows
for merging competing changes on documents. This is a major difference between Sops
and other encryption tools that store documents as encrypted blobs.

In BINARY format, the cleartext data is treated as a single blob and the encrypted
document is in JSON format with a single `data` key and a single encrypted value.

Sops allows operators to encrypt their documents with multiple master keys. Each of
the master key defined in the document is able to decrypt it, allowing users to
share documents amongst themselves without sharing keys, or using a PGP key as a
backup for KMS.

In practice, this is achieved by generating a data key for each document that is used
to encrypt all values, and encrypting the data with each master key defined. Being
able to decrypt the data key gives access to the document.

The integrity of each document is guaranteed by calculating a Message Authentication Code
(MAC) that is stored encrypted by the data key. When decrypting a document, the MAC should
be recalculated and compared with the MAC stored in the document to verify that no
fraudulent changes have been applied. The MAC covers keys and values as well as their
ordering.
*/
package sops // import "github.com/getsops/sops/v3"

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/audit"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/logging"
	"github.com/getsops/sops/v3/pgp"
	"github.com/getsops/sops/v3/shamir"
)

// DefaultUnencryptedSuffix is the default suffix a TreeItem key has to end with for sops to leave its Value unencrypted
const DefaultUnencryptedSuffix = "_unencrypted"

var DefaultDecryptionOrder = []string{age.KeyTypeIdentifier, pgp.KeyTypeIdentifier}

type sopsError string

func (e sopsError) Error() string {
	return string(e)
}

// MacMismatch occurs when the computed MAC does not match the expected ones
const MacMismatch = sopsError("MAC mismatch")

// MetadataNotFound occurs when the input file is malformed and doesn't have sops metadata in it
const MetadataNotFound = sopsError("sops metadata not found")

type SopsKeyNotFound struct {
	Key interface{}
	Msg string
}

func (e *SopsKeyNotFound) Error() string {
	return fmt.Sprintf(e.Msg, e.Key)
}

// MACOnlyEncryptedInitialization is a constant and known sequence of 32 bytes used to initialize
// MAC which is computed only over values which end up encrypted. That assures that a MAC with the
// setting enabled is always different from a MAC with this setting disabled.
// The following numbers are taken from the output of `echo -n sops | sha256sum` (shell) or `hashlib.sha256(b'sops').hexdigest()` (Python).
var MACOnlyEncryptedInitialization = []byte{0x8a, 0x3f, 0xd2, 0xad, 0x54, 0xce, 0x66, 0x52, 0x7b, 0x10, 0x34, 0xf3, 0xd1, 0x47, 0xbe, 0xb, 0xb, 0x97, 0x5b, 0x3b, 0xf4, 0x4f, 0x72, 0xc6, 0xfd, 0xad, 0xec, 0x81, 0x76, 0xf2, 0x7d, 0x69}

var log *logrus.Logger

func init() {
	log = logging.NewLogger("SOPS")
}

// Cipher provides a way to encrypt and decrypt the data key used to encrypt and decrypt sops files, so that the
// data key can be stored alongside the encrypted content. A Cipher must be able to decrypt the values it encrypts.
type Cipher interface {
	// Encrypt takes a plaintext, a key and additional data and returns the plaintext encrypted with the key, using the
	// additional data for authentication
	Encrypt(plaintext interface{}, key []byte, additionalData string) (ciphertext string, err error)
	// Encrypt takes a ciphertext, a key and additional data and returns the ciphertext encrypted with the key, using
	// the additional data for authentication
	Decrypt(ciphertext string, key []byte, additionalData string) (plaintext interface{}, err error)
}

// Comment represents a comment in the sops tree for the file formats that actually support them.
type Comment struct {
	Value string
}

// TreeItem is an item inside sops's tree
type TreeItem struct {
	Key   interface{}
	Value interface{}
}

// TreeBranch is a branch inside sops's tree. It is a slice of TreeItems and is therefore ordered
type TreeBranch []TreeItem

// TreeBranches is a collection of TreeBranch
// Trees usually have more than one branch
type TreeBranches []TreeBranch

func equals(oneBranch interface{}, otherBranch interface{}) bool {
	switch oneBranch := oneBranch.(type) {
	case TreeBranch:
		otherBranch, ok := otherBranch.(TreeBranch)
		if !ok || len(oneBranch) != len(otherBranch) {
			return false
		}
		for i, item := range oneBranch {
			otherItem := otherBranch[i]
			if !equals(item.Key, otherItem.Key) || !equals(item.Value, otherItem.Value) {
				return false
			}
		}
		return true
	case []interface{}:
		otherBranch, ok := otherBranch.([]interface{})
		if !ok || len(oneBranch) != len(otherBranch) {
			return false
		}
		for i, item := range oneBranch {
			if !equals(item, otherBranch[i]) {
				return false
			}
		}
		return true
	case Comment:
		otherBranch, ok := otherBranch.(Comment)
		if !ok {
			return false
		}
		return oneBranch.Value == otherBranch.Value
	default:
		// Unexpected type
		return oneBranch == otherBranch
	}
}

// Compare a branch with another one
func (branch TreeBranch) Equals(other TreeBranch) bool {
	return equals(branch, other)
}

func valueFromPathAndLeaf(path []interface{}, leaf interface{}) interface{} {
	switch component := path[0].(type) {
	case int:
		if len(path) == 1 {
			return []interface{}{
				leaf,
			}
		}
		return []interface{}{
			valueFromPathAndLeaf(path[1:], leaf),
		}
	default:
		if len(path) == 1 {
			return TreeBranch{
				TreeItem{
					Key:   component,
					Value: leaf,
				},
			}
		}
		return TreeBranch{
			TreeItem{
				Key:   component,
				Value: valueFromPathAndLeaf(path[1:], leaf),
			},
		}
	}
}

func set(branch interface{}, path []interface{}, value interface{}) (interface{}, bool) {
	switch branch := branch.(type) {
	case TreeBranch:
		for i, item := range branch {
			if item.Key == path[0] {
				var changed bool
				if len(path) == 1 {
					changed = !equals(branch[i].Value, value)
					branch[i].Value = value
				} else {
					branch[i].Value, changed = set(item.Value, path[1:], value)
				}
				return branch, changed
			}
		}
		// Not found, need to add the next path entry to the branch
		value := valueFromPathAndLeaf(path, value)
		if newBranch, ok := value.(TreeBranch); ok && len(newBranch) > 0 {
			return append(branch, newBranch[0]), true
		}
		return branch, true
	case []interface{}:
		position := path[0].(int)
		var changed bool
		if len(path) == 1 {
			if position >= len(branch) {
				return append(branch, value), true
			}
			changed = !equals(branch[position], value)
			branch[position] = value
		} else {
			if position >= len(branch) {
				branch = append(branch, valueFromPathAndLeaf(path[1:], value))
				changed = true
			} else {
				branch[position], changed = set(branch[position], path[1:], value)
			}
		}
		return branch, changed
	default:
		newValue := valueFromPathAndLeaf(path, value)
		return newValue, !equals(branch, newValue)
	}
}

// Set sets a value on a given tree for the specified path
func (branch TreeBranch) Set(path []interface{}, value interface{}) (TreeBranch, bool) {
	v, changed := set(branch, path, value)
	return v.(TreeBranch), changed
}

func unset(branch interface{}, path []interface{}) (interface{}, error) {
	switch branch := branch.(type) {
	case TreeBranch:
		for i, item := range branch {
			if item.Key == path[0] {
				if len(path) == 1 {
					branch = slices.Delete(branch, i, i+1)
				} else {
					v, err := unset(item.Value, path[1:])
					if err != nil {
						return nil, err
					}
					branch[i].Value = v
				}
				return branch, nil
			}
		}
		return nil, &SopsKeyNotFound{Msg: "Key not found: %s", Key: path[0]}
	case []interface{}:
		position := path[0].(int)
		if position >= len(branch) {
			return nil, &SopsKeyNotFound{Msg: "Index %d out of bounds", Key: path[0]}
		}
		if len(path) == 1 {
			branch = slices.Delete(branch, position, position+1)
		} else {
			v, err := unset(branch[position], path[1:])
			if err != nil {
				return nil, err
			}
			branch[position] = v
		}
		return branch, nil
	default:
		return nil, fmt.Errorf("Unsupported type: %T for item '%s'", branch, path[0])
	}
}

// Unset removes a value on a given tree from the specified path
func (branch TreeBranch) Unset(path []interface{}) (TreeBranch, error) {
	v, err := unset(branch, path)
	if err != nil {
		return nil, err
	}
	return v.(TreeBranch), nil
}

// Tree is the data structure used by sops to represent documents internally
type Tree struct {
	Metadata Metadata
	Branches TreeBranches
	// FilePath is the path of the file this struct represents
	FilePath string
}

// Truncate truncates the tree to the path specified
func (branch TreeBranch) Truncate(path []interface{}) (interface{}, error) {
	log.WithField("path", path).Info("Truncating tree")
	var current interface{} = branch
	for _, component := range path {
		switch component := component.(type) {
		case string:
			found := false
			for _, item := range current.(TreeBranch) {
				if item.Key == component {
					current = item.Value
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("component ['%s'] not found", component)
			}
		case int:
			if reflect.ValueOf(current).Kind() != reflect.Slice {
				return nil, fmt.Errorf("component [%d] is integer, but tree part is not a slice", component)
			}
			if reflect.ValueOf(current).Len() <= component {
				return nil, fmt.Errorf("component [%d] accesses out of bounds", component)
			}
			current = reflect.ValueOf(current).Index(component).Interface()
		}
	}
	return current, nil
}

func (branch TreeBranch) walkValue(in interface{}, path []string, commentsStack [][]string, onLeaves func(in interface{}, path []string, commentsStack [][]string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, path, commentsStack)
	case []byte:
		return onLeaves(string(in), path, commentsStack)
	case int:
		return onLeaves(in, path, commentsStack)
	case bool:
		return onLeaves(in, path, commentsStack)
	case float64:
		return onLeaves(in, path, commentsStack)
	case time.Time:
		return onLeaves(in, path, commentsStack)
	case Comment:
		return onLeaves(in, path, commentsStack)
	case TreeBranch:
		return branch.walkBranch(in, path, commentsStack, onLeaves)
	case []interface{}:
		return branch.walkSlice(in, path, commentsStack, onLeaves)
	case nil:
		// the value returned remains the same since it doesn't make
		// sense to encrypt or decrypt a nil value
		return nil, nil
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (branch TreeBranch) walkSlice(in []interface{}, path []string, commentsStack [][]string, onLeaves func(in interface{}, path []string, commentsStack [][]string) (interface{}, error)) ([]interface{}, error) {
	// Because append returns a new slice, the original stack is not changed.
	commentsStack = append(commentsStack, []string{})
	for i, v := range in {
		c, vIsComment := v.(Comment)
		if vIsComment {
			// If v is a comment, we add it to the slice of active comments.
			// This allows us to also encrypt comments themselves by enabling encryption in a prior comment.
			commentsStack[len(commentsStack)-1] = append(commentsStack[len(commentsStack)-1], c.Value)
		}
		newV, err := branch.walkValue(v, path, commentsStack, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
		if !vIsComment {
			// If v is not a comment, we clear the slice of active comments.
			commentsStack[len(commentsStack)-1] = []string{}
		}
	}
	return in, nil
}

func (branch TreeBranch) walkBranch(in TreeBranch, path []string, commentsStack [][]string, onLeaves func(in interface{}, path []string, commentsStack [][]string) (interface{}, error)) (TreeBranch, error) {
	// Because append returns a new slice, the original stack is not changed.
	commentsStack = append(commentsStack, []string{})
	for i, item := range in {
		if c, ok := item.Key.(Comment); ok {
			// If key is a comment, we add it to the slice of active comments.
			// This allows us to also encrypt comments themselves by enabling encryption in a prior comment.
			commentsStack[len(commentsStack)-1] = append(commentsStack[len(commentsStack)-1], c.Value)
			enc, err := branch.walkValue(item.Key, path, commentsStack, onLeaves)
			if err != nil {
				return nil, err
			}
			if encComment, ok := enc.(Comment); ok {
				in[i].Key = encComment
				continue
			} else if comment, ok := enc.(string); ok {
				in[i].Key = Comment{Value: comment}
				continue
			} else {
				return nil, fmt.Errorf("walkValue of Comment should be either Comment or string, was %T", enc)
			}
		}
		c, valueIsComment := item.Value.(Comment)
		if valueIsComment {
			// If value is a comment, we add it to the slice of active comments.
			// This allows us to also encrypt comments themselves by enabling encryption in a prior comment.
			commentsStack[len(commentsStack)-1] = append(commentsStack[len(commentsStack)-1], c.Value)
		}
		key, ok := item.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Tree contains a non-string key (type %T): %s. Only string keys are"+
				"supported", item.Key, item.Key)
		}
		newV, err := branch.walkValue(item.Value, append(path, key), commentsStack, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
		if !valueIsComment {
			// If value is not a comment, we clear the slice of active comments.
			commentsStack[len(commentsStack)-1] = []string{}
		}
	}
	return in, nil
}

func (tree Tree) shouldBeEncrypted(path []string, commentsStack [][]string, isComment bool) bool {
	encrypted := true
	if tree.Metadata.UnencryptedSuffix != "" {
		for _, v := range path {
			if strings.HasSuffix(v, tree.Metadata.UnencryptedSuffix) {
				encrypted = false
				break
			}
		}
	}
	if tree.Metadata.EncryptedSuffix != "" {
		encrypted = false
		for _, v := range path {
			if strings.HasSuffix(v, tree.Metadata.EncryptedSuffix) {
				encrypted = true
				break
			}
		}
	}
	if tree.Metadata.UnencryptedRegex != "" {
		for _, p := range path {
			matched, _ := regexp.Match(tree.Metadata.UnencryptedRegex, []byte(p))
			if matched {
				encrypted = false
				break
			}
		}
	}
	if tree.Metadata.EncryptedRegex != "" {
		encrypted = false
		for _, p := range path {
			matched, _ := regexp.Match(tree.Metadata.EncryptedRegex, []byte(p))
			if matched {
				encrypted = true
				break
			}
		}
	}
	if tree.Metadata.UnencryptedCommentRegex != "" {
	unencryptedComments:
		for _, cs := range commentsStack {
			for _, c := range cs {
				matched, _ := regexp.Match(tree.Metadata.UnencryptedCommentRegex, []byte(c))
				if matched {
					encrypted = false
					break unencryptedComments
				}
			}
		}
	}
	if tree.Metadata.EncryptedCommentRegex != "" {
		lenCommentsStack := len(commentsStack)
		lenLastCommentsStack := len(commentsStack[lenCommentsStack-1])
		encrypted = false
	encryptedComments:
		for i, cs := range commentsStack {
			for j, c := range cs {
				// A special case. We do not encrypt the comment line itself which matches the regex.
				// So we skip the last line of the last set of comments. Only if the matches any previous
				// line, we encrypt this comment. Otherwise we do not.
				if isComment && i == lenCommentsStack-1 && j == lenLastCommentsStack-1 {
					continue
				}
				matched, _ := regexp.Match(tree.Metadata.EncryptedCommentRegex, []byte(c))
				if matched {
					encrypted = true
					break encryptedComments
				}
			}
		}
	}
	return encrypted
}

// Encrypt walks over the tree and encrypts all values with the provided cipher,
// except those whose key ends with the UnencryptedSuffix specified on the
// Metadata struct, those not ending with EncryptedSuffix, if EncryptedSuffix
// is provided (by default it is not), those not matching EncryptedRegex,
// if EncryptedRegex is provided (by default it is not), those matching UnencryptedRegex,
// if UnencryptedRegex is provided (by default it is not), those with their comment
// not matching EncryptedCommentRegex, if EncryptedCommentRegex is provided (by default
// it is not), or those with their comment matching UnencryptedCommentRegex, if
// UnencryptedCommentRegex is provided (by default it is not).
// If encryption is successful, it returns the MAC for the encrypted tree
// (all values if MACOnlyEncrypted is false, or only over values which end
// up encrypted if MACOnlyEncrypted is true).
func (tree Tree) Encrypt(key []byte, cipher Cipher) (string, error) {
	audit.SubmitEvent(audit.EncryptEvent{
		File: tree.FilePath,
	})
	hash := sha512.New()
	if tree.Metadata.MACOnlyEncrypted {
		// We initialize with known set of bytes so that a MAC with this setting
		// enabled is always different from a MAC with this setting disabled.
		hash.Write(MACOnlyEncryptedInitialization)
	}
	walk := func(branch TreeBranch) error {
		_, err := branch.walkBranch(branch, make([]string, 0), make([][]string, 0), func(in interface{}, path []string, commentsStack [][]string) (interface{}, error) {
			_, ok := in.(Comment)
			encrypted := tree.shouldBeEncrypted(path, commentsStack, ok)
			if !tree.Metadata.MACOnlyEncrypted || encrypted {
				// Only add to MAC if not a comment
				if !ok {
					bytes, err := ToBytes(in)
					if err != nil {
						return nil, fmt.Errorf("Could not convert %s to bytes: %s", in, err)
					}
					hash.Write(bytes)
				}
			}
			if encrypted {
				var err error
				pathString := strings.Join(path, ":") + ":"
				in, err = cipher.Encrypt(in, key, pathString)
				if err != nil {
					return nil, fmt.Errorf("Could not encrypt value: %s", err)
				}
				if ok && tree.Metadata.UnencryptedCommentRegex != "" {
					// If an encrypted comment matches tree.Metadata.UnencryptedCommentRegex, decryption will fail
					// as the MAC does not match, and the commented value will not be decrypted.
					matched, _ := regexp.Match(tree.Metadata.UnencryptedCommentRegex, []byte(in.(string)))
					if matched {
						return nil, fmt.Errorf("Encrypted comment %q matches UnencryptedCommentRegex! Make sure that UnencryptedCommentRegex cannot match an encrypted comment.", in)
					}
				}
			}
			return in, nil
		})
		return err
	}

	for _, branch := range tree.Branches {
		err := walk(branch)
		if err != nil {
			return "", fmt.Errorf("Error walking tree: %s", err)
		}
	}
	return fmt.Sprintf("%X", hash.Sum(nil)), nil
}

// Decrypt walks over the tree and decrypts all values with the provided cipher,
// except those whose key ends with the UnencryptedSuffix specified on the Metadata struct,
// those not ending with EncryptedSuffix, if EncryptedSuffix is provided (by default it is not),
// those not matching EncryptedRegex, if EncryptedRegex is provided (by default it is not),
// or those matching UnencryptedRegex, if UnencryptedRegex is provided (by default it is not).
// If decryption is successful, it returns the MAC for the decrypted tree
// (all values if MACOnlyEncrypted is false, or only over values which end
// up decrypted if MACOnlyEncrypted is true).
func (tree Tree) Decrypt(key []byte, cipher Cipher) (string, error) {
	log.Debug("Decrypting tree")
	audit.SubmitEvent(audit.DecryptEvent{
		File: tree.FilePath,
	})
	hash := sha512.New()
	if tree.Metadata.MACOnlyEncrypted {
		// We initialize with known set of bytes so that a MAC with this setting
		// enabled is always different from a MAC with this setting disabled.
		hash.Write(MACOnlyEncryptedInitialization)
	}
	walk := func(branch TreeBranch) error {
		_, err := branch.walkBranch(branch, make([]string, 0), make([][]string, 0), func(in interface{}, path []string, commentsStack [][]string) (interface{}, error) {
			c, ok := in.(Comment)
			encrypted := tree.shouldBeEncrypted(path, commentsStack, ok)
			var v interface{}
			if encrypted {
				var err error
				pathString := strings.Join(path, ":") + ":"
				if ok {
					v, err = cipher.Decrypt(c.Value, key, pathString)
					if err != nil {
						// Assume the comment was not encrypted in the first place
						log.WithField("comment", c.Value).
							Warn("Found possibly unencrypted comment in file. " +
								"This is to be expected if the file being " +
								"decrypted was created with an older version of " +
								"SOPS.")
						v = c
					}
				} else {
					v, err = cipher.Decrypt(in.(string), key, pathString)
					if err != nil {
						return nil, fmt.Errorf("Could not decrypt value: %s", err)
					}
				}
			} else {
				v = in
			}
			if !tree.Metadata.MACOnlyEncrypted || encrypted {
				// Only add to MAC if not a comment
				if _, ok := v.(Comment); !ok {
					bytes, err := ToBytes(v)
					if err != nil {
						return nil, fmt.Errorf("Could not convert %s to bytes: %s", in, err)
					}
					hash.Write(bytes)
				}
			}
			return v, nil
		})
		return err
	}
	for _, branch := range tree.Branches {
		err := walk(branch)
		if err != nil {
			return "", fmt.Errorf("Error walking tree: %s", err)
		}
	}
	return fmt.Sprintf("%X", hash.Sum(nil)), nil
}

// GenerateDataKey generates a new random data key and encrypts it with all MasterKeys.
func (tree Tree) GenerateDataKey() ([]byte, []error) {
	newKey := make([]byte, 32)
	_, err := rand.Read(newKey)
	if err != nil {
		return nil, []error{fmt.Errorf("Could not generate random key: %s", err)}
	}
	return newKey, tree.Metadata.UpdateMasterKeys(newKey)
}

// GenerateDataKeyWithKeyServices generates a new random data key and encrypts it with all MasterKeys.
func (tree *Tree) GenerateDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient) ([]byte, []error) {
	newKey := make([]byte, 32)
	_, err := rand.Read(newKey)
	if err != nil {
		return nil, []error{fmt.Errorf("Could not generate random key: %s", err)}
	}
	return newKey, tree.Metadata.UpdateMasterKeysWithKeyServices(newKey, svcs)
}

// Metadata holds information about a file encrypted by sops
type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	EncryptedSuffix           string
	UnencryptedRegex          string
	EncryptedRegex            string
	UnencryptedCommentRegex   string
	EncryptedCommentRegex     string
	MessageAuthenticationCode string
	MACOnlyEncrypted          bool
	Version                   string
	KeyGroups                 []KeyGroup
	// ShamirThreshold is the number of key groups required to recover the
	// original data key
	ShamirThreshold int
	// DataKey caches the decrypted data key so it doesn't have to be decrypted with a master key every time it's needed
	DataKey []byte
}

// KeyGroup is a slice of SOPS MasterKeys that all encrypt the same part of the data key
type KeyGroup []keys.MasterKey

// EncryptedFileLoader is the interface for loading of encrypted files. It provides a
// way to load encrypted SOPS files into the internal SOPS representation. Because it
// loads encrypted files, the returned data structure already contains all SOPS
// metadata.
type EncryptedFileLoader interface {
	LoadEncryptedFile(in []byte) (Tree, error)
}

// PlainFileLoader is the interface for loading of plain text files. It provides a
// way to load unencrypted files into SOPS. Because the files it loads are
// unencrypted, the returned data structure does not contain any metadata.
type PlainFileLoader interface {
	LoadPlainFile(in []byte) (TreeBranches, error)
}

// EncryptedFileEmitter is the interface for emitting encrypting files. It provides a
// way to emit encrypted files from the internal SOPS representation.
type EncryptedFileEmitter interface {
	EmitEncryptedFile(Tree) ([]byte, error)
}

// PlainFileEmitter is the interface for emitting plain text files. It provides a way
// to emit plain text files from the internal SOPS representation so that they can be
// shown
type PlainFileEmitter interface {
	EmitPlainFile(TreeBranches) ([]byte, error)
}

// ValueEmitter is the interface for emitting a value. It provides a way to emit
// values from the internal SOPS representation so that they can be shown
type ValueEmitter interface {
	EmitValue(interface{}) ([]byte, error)
}

// CheckEncrypted is the interface for testing whether a branch contains sops
// metadata. This is used to check whether a file is already encrypted or not.
type CheckEncrypted interface {
	HasSopsTopLevelKey(TreeBranch) bool
}

// Store is used to interact with files, both encrypted and unencrypted.
type Store interface {
	EncryptedFileLoader
	PlainFileLoader
	EncryptedFileEmitter
	PlainFileEmitter
	ValueEmitter
	CheckEncrypted
}

// MasterKeyCount returns the number of master keys available
func (m *Metadata) MasterKeyCount() int {
	count := 0
	for _, group := range m.KeyGroups {
		count += len(group)
	}
	return count
}

// UpdateMasterKeysWithKeyServices encrypts the data key with all master keys using the provided key services
func (m *Metadata) UpdateMasterKeysWithKeyServices(dataKey []byte, svcs []keyservice.KeyServiceClient) (errs []error) {
	if len(svcs) == 0 {
		return []error{
			fmt.Errorf("no key services provided, cannot update master keys"),
		}
	}
	if len(m.KeyGroups) == 0 {
		return []error{
			fmt.Errorf("no key groups provided"),
		}
	}
	var parts [][]byte
	if len(m.KeyGroups) == 1 {
		// If there's only one key group, we can't do Shamir. All keys
		// in the group encrypt the whole data key.
		parts = append(parts, dataKey)
	} else {
		var err error
		if m.ShamirThreshold == 0 {
			m.ShamirThreshold = len(m.KeyGroups)
		}
		log.WithFields(logrus.Fields{
			"quorum": m.ShamirThreshold,
			"parts":  len(m.KeyGroups),
		}).Info("Splitting data key with Shamir Secret Sharing")
		parts, err = shamir.Split(dataKey, len(m.KeyGroups), int(m.ShamirThreshold))
		if err != nil {
			errs = append(errs, fmt.Errorf("could not split data key into parts for Shamir: %s", err))
			return
		}
		if len(parts) != len(m.KeyGroups) {
			errs = append(errs, fmt.Errorf("not enough parts obtained from Shamir: need %d, got %d", len(m.KeyGroups), len(parts)))
			return
		}
	}
	for i, group := range m.KeyGroups {
		part := parts[i]
		if len(group) == 0 {
			return []error{
				fmt.Errorf("empty key group provided"),
			}
		}
		for _, key := range group {
			svcKey := keyservice.KeyFromMasterKey(key)
			var keyErrs []error
			encrypted := false
			for _, svc := range svcs {
				rsp, err := svc.Encrypt(context.Background(), &keyservice.EncryptRequest{
					Key:       &svcKey,
					Plaintext: part,
				})
				if err != nil {
					keyErrs = append(keyErrs, fmt.Errorf("failed to encrypt new data key with master key %q: %w", key.ToString(), err))
					continue
				}
				key.SetEncryptedDataKey(rsp.Ciphertext)
				encrypted = true
				// Only need to encrypt the key successfully with one service
				break
			}
			if !encrypted {
				errs = append(errs, keyErrs...)
			}
		}
	}
	m.DataKey = dataKey
	return
}

// UpdateMasterKeys encrypts the data key with all master keys
func (m *Metadata) UpdateMasterKeys(dataKey []byte) (errs []error) {
	return m.UpdateMasterKeysWithKeyServices(dataKey, []keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	})
}

// GetDataKeyWithKeyServices retrieves the data key, asking KeyServices to decrypt it with each
// MasterKey in the Metadata's KeySources until one of them succeeds.
func (m *Metadata) GetDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient, decryptionOrder []string) ([]byte, error) {
	if m.DataKey != nil {
		return m.DataKey, nil
	}
	getDataKeyErr := getDataKeyError{
		RequiredSuccessfulKeyGroups: m.ShamirThreshold,
		GroupResults:                make([]error, len(m.KeyGroups)),
	}
	var parts [][]byte
	for i, group := range m.KeyGroups {
		part, err := decryptKeyGroup(group, svcs, decryptionOrder)
		if err == nil {
			parts = append(parts, part)
		}
		getDataKeyErr.GroupResults[i] = err
	}
	var dataKey []byte
	if len(m.KeyGroups) > 1 {
		if len(parts) < m.ShamirThreshold {
			return nil, &getDataKeyErr
		}
		var err error
		dataKey, err = shamir.Combine(parts)
		if err != nil {
			return nil, fmt.Errorf("could not get data key from shamir parts: %s", err)
		}
	} else {
		if len(parts) != 1 {
			return nil, &getDataKeyErr
		}
		dataKey = parts[0]
	}
	log.Info("Data key recovered successfully")
	m.DataKey = dataKey
	return dataKey, nil
}

// decryptKeyGroup tries to decrypt the contents of the provided KeyGroup with
// any of the MasterKeys in the KeyGroup with any of the provided key services,
// returning as soon as one key service succeeds.
func decryptKeyGroup(group KeyGroup, svcs []keyservice.KeyServiceClient, decryptionOrder []string) ([]byte, error) {
	var keyErrs []error
	// Sort MasterKeys in the group so we try them in specific order
	// Use sorted indices to avoid group slice modification
	indices := sortKeyGroupIndices(group, decryptionOrder)
	for _, indexVal := range indices {
		key := group[indexVal]
		part, err := decryptKey(key, svcs)
		if err != nil {
			keyErrs = append(keyErrs, err)
		} else {
			return part, nil
		}
	}
	return nil, decryptKeyErrors(keyErrs)
}

// sortKeyGroupIndices returns indices that would sort the KeyGroup
// according to decryptionOrder
func sortKeyGroupIndices(group KeyGroup, decryptionOrder []string) []int {
	priorities := make(map[string]int)
	// give ordered weights
	for i, v := range decryptionOrder {
		priorities[v] = i
	}
	maxPriority := len(decryptionOrder)
	// initialize indices
	n := len(group)
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	sort.SliceStable(indices, func(i, j int) bool {
		keyTypeI := group[indices[i]].TypeToIdentifier()
		keyTypeJ := group[indices[j]].TypeToIdentifier()
		priorityI, ok := priorities[keyTypeI]
		if !ok {
			priorityI = maxPriority
		}
		priorityJ, ok := priorities[keyTypeJ]
		if !ok {
			priorityJ = maxPriority
		}
		return priorityI < priorityJ
	})
	return indices
}

// decryptKey tries to decrypt the contents of the provided MasterKey with any
// of the key services, returning as soon as one key service succeeds.
func decryptKey(key keys.MasterKey, svcs []keyservice.KeyServiceClient) ([]byte, error) {
	svcKey := keyservice.KeyFromMasterKey(key)
	var part []byte
	decryptErr := decryptKeyError{
		keyName: key.ToString(),
	}
	for _, svc := range svcs {
		// All keys in a key group encrypt the same part, so as soon
		// as we decrypt it successfully with one key, we need to
		// proceed with the next group
		var err error
		if part == nil {
			var rsp *keyservice.DecryptResponse
			rsp, err = svc.Decrypt(
				context.Background(),
				&keyservice.DecryptRequest{
					Ciphertext: key.EncryptedDataKey(),
					Key:        &svcKey,
				})
			if err == nil {
				part = rsp.Plaintext
			}
		}
		decryptErr.errs = append(decryptErr.errs, err)
	}
	if part != nil {
		return part, nil
	}
	return nil, &decryptErr
}

// GetDataKey retrieves the data key from the first MasterKey in the Metadata's KeySources that's able to return it,
// using the local KeyService
func (m Metadata) GetDataKey() ([]byte, error) {
	return m.GetDataKeyWithKeyServices([]keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	}, nil)
}

// ToBytes converts a string, int, float or bool to a byte representation.
func ToBytes(in interface{}) ([]byte, error) {
	switch in := in.(type) {
	case string:
		return []byte(in), nil
	case int:
		return []byte(strconv.Itoa(in)), nil
	case float64:
		return []byte(strconv.FormatFloat(in, 'f', -1, 64)), nil
	case bool:
		boolB := []byte("True")
		if !in {
			boolB = []byte("False")
		}
		return boolB, nil
	case []byte:
		return in, nil
	case time.Time:
		return in.MarshalText()
	case Comment:
		return ToBytes(in.Value)
	default:
		return nil, fmt.Errorf("Could not convert unknown type %T to bytes", in)
	}
}

// EmitAsMap will emit the tree branches as a map. This is used by the publish
// command for writing decrypted trees to various destinations. Should only be
// used for outputting to data structures in code.
func EmitAsMap(in TreeBranches) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	for _, branch := range in {
		for _, item := range branch {
			if _, ok := item.Key.(Comment); ok {
				continue
			}
			val, err := encodeValueForMap(item.Value)
			if err != nil {
				return nil, err
			}
			data[item.Key.(string)] = val
		}
	}

	return data, nil
}

func encodeValueForMap(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case TreeBranch:
		return EmitAsMap([]TreeBranch{v})
	default:
		return v, nil
	}
}
