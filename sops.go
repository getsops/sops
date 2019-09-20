/*
Package sops manages JSON, YAML and BINARY documents to be encrypted or decrypted.

This package should not be used directly. Instead, Sops users should install the
command line client via `go get -u go.mozilla.org/sops/cmd/sops`, or use the
decryption helper provided at `go.mozilla.org/sops/decrypt`.

We do not guarantee API stability for any package other than `go.mozilla.org/sops/decrypt`.

A Sops document is a Tree composed of a data branch with arbitrary key/value pairs
and a metadata branch with encryption and integrity information.

In JSON and YAML formats, the structure of the cleartext tree is preserved, keys are
stored in cleartext and only values are encrypted. Keeping the values in cleartext
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
package sops //import "go.mozilla.org/sops"

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/audit"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/logging"
	"go.mozilla.org/sops/shamir"
	"golang.org/x/net/context"
)

// DefaultUnencryptedSuffix is the default suffix a TreeItem key has to end with for sops to leave its Value unencrypted
const DefaultUnencryptedSuffix = "_unencrypted"

type sopsError string

func (e sopsError) Error() string {
	return string(e)
}

// MacMismatch occurs when the computed MAC does not match the expected ones
const MacMismatch = sopsError("MAC mismatch")

// MetadataNotFound occurs when the input file is malformed and doesn't have sops metadata in it
const MetadataNotFound = sopsError("sops metadata not found")

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

func set(branch interface{}, path []interface{}, value interface{}) interface{} {
	switch branch := branch.(type) {
	case TreeBranch:
		for i, item := range branch {
			if item.Key == path[0] {
				if len(path) == 1 {
					branch[i].Value = value
				} else {
					branch[i].Value = set(item.Value, path[1:], value)
				}
				return branch
			}
		}
		// Not found, need to add the next path entry to the branch
		if len(path) == 1 {
			return append(branch, TreeItem{Key: path[0], Value: value})
		}
		return valueFromPathAndLeaf(path, value)
	case []interface{}:
		position := path[0].(int)
		if len(path) == 1 {
			if position >= len(branch) {
				return append(branch, value)
			}
			branch[position] = value
		} else {
			if position >= len(branch) {
				branch = append(branch, valueFromPathAndLeaf(path[1:], value))
			}
			branch[position] = set(branch[position], path[1:], value)
		}
		return branch
	default:
		return valueFromPathAndLeaf(path, value)
	}
}

// Set sets a value on a given tree for the specified path
func (branch TreeBranch) Set(path []interface{}, value interface{}) TreeBranch {
	return set(branch, path, value).(TreeBranch)
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

func (branch TreeBranch) walkValue(in interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, path)
	case []byte:
		return onLeaves(string(in), path)
	case int:
		return onLeaves(in, path)
	case bool:
		return onLeaves(in, path)
	case float64:
		return onLeaves(in, path)
	case Comment:
		return onLeaves(in, path)
	case TreeBranch:
		return branch.walkBranch(in, path, onLeaves)
	case []interface{}:
		return branch.walkSlice(in, path, onLeaves)
	case nil:
		// the value returned remains the same since it doesn't make
		// sense to encrypt or decrypt a nil value
		return nil, nil
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (branch TreeBranch) walkSlice(in []interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) ([]interface{}, error) {
	for i, v := range in {
		newV, err := branch.walkValue(v, path, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
	}
	return in, nil
}

func (branch TreeBranch) walkBranch(in TreeBranch, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (TreeBranch, error) {
	for i, item := range in {
		if _, ok := item.Key.(Comment); ok {
			enc, err := branch.walkValue(item.Key, path, onLeaves)
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
		key, ok := item.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Tree contains a non-string key (type %T): %s. Only string keys are"+
				"supported", item.Key, item.Key)
		}
		newV, err := branch.walkValue(item.Value, append(path, key), onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
	}
	return in, nil
}

// Encrypt walks over the tree and encrypts all values with the provided cipher,
// except those whose key ends with the UnencryptedSuffix specified on the
// Metadata struct, those not ending with EncryptedSuffix, if EncryptedSuffix
// is provided (by default it is not), or those not matching EncryptedRegex,
// if EncryptedRegex is provided (by default it is not).  If encryption is
// successful, it returns the MAC for the encrypted tree.
func (tree Tree) Encrypt(key []byte, cipher Cipher) (string, error) {
	audit.SubmitEvent(audit.EncryptEvent{
		File: tree.FilePath,
	})
	hash := sha512.New()
	walk := func(branch TreeBranch) error {
		_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
			// Only add to MAC if not a comment
			if _, ok := in.(Comment); !ok {
				bytes, err := ToBytes(in)
				if err != nil {
					return nil, fmt.Errorf("Could not convert %s to bytes: %s", in, err)
				}
				hash.Write(bytes)
			}
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
			if encrypted {
				var err error
				pathString := strings.Join(path, ":") + ":"
				in, err = cipher.Encrypt(in, key, pathString)
				if err != nil {
					return nil, fmt.Errorf("Could not encrypt value: %s", err)
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
// or those not matching EncryptedRegex, if EncryptedRegex is provided (by default it is not).
// If decryption is successful, it returns the MAC for the decrypted tree.
func (tree Tree) Decrypt(key []byte, cipher Cipher) (string, error) {
	log.Debug("Decrypting tree")
	audit.SubmitEvent(audit.DecryptEvent{
		File: tree.FilePath,
	})
	hash := sha512.New()
	walk := func(branch TreeBranch) error {
		_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
			encrypted := true
			if tree.Metadata.UnencryptedSuffix != "" {
				for _, p := range path {
					if strings.HasSuffix(p, tree.Metadata.UnencryptedSuffix) {
						encrypted = false
						break
					}
				}
			}
			if tree.Metadata.EncryptedSuffix != "" {
				encrypted = false
				for _, p := range path {
					if strings.HasSuffix(p, tree.Metadata.EncryptedSuffix) {
						encrypted = true
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
			var v interface{}
			if encrypted {
				var err error
				pathString := strings.Join(path, ":") + ":"
				if c, ok := in.(Comment); ok {
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
			// Only add to MAC if not a comment
			if _, ok := v.(Comment); !ok {
				bytes, err := ToBytes(v)
				if err != nil {
					return nil, fmt.Errorf("Could not convert %s to bytes: %s", in, err)
				}
				hash.Write(bytes)
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
	EncryptedRegex            string
	MessageAuthenticationCode string
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

// Store is used to interact with files, both encrypted and unencrypted.
type Store interface {
	EncryptedFileLoader
	PlainFileLoader
	EncryptedFileEmitter
	PlainFileEmitter
	ValueEmitter
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
					keyErrs = append(keyErrs, fmt.Errorf("failed to encrypt new data key with master key %q: %v", key.ToString(), err))
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
func (m Metadata) GetDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient) ([]byte, error) {
	if m.DataKey != nil {
		return m.DataKey, nil
	}
	getDataKeyErr := getDataKeyError{
		RequiredSuccessfulKeyGroups: m.ShamirThreshold,
		GroupResults:                make([]error, len(m.KeyGroups)),
	}
	var parts [][]byte
	for i, group := range m.KeyGroups {
		part, err := decryptKeyGroup(group, svcs)
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
func decryptKeyGroup(group KeyGroup, svcs []keyservice.KeyServiceClient) ([]byte, error) {
	var keyErrs []error
	for _, key := range group {
		part, err := decryptKey(key, svcs)
		if err != nil {
			keyErrs = append(keyErrs, err)
		} else {
			return part, nil
		}
	}
	return nil, decryptKeyErrors(keyErrs)
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
	})
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
		return []byte(strings.Title(strconv.FormatBool(in))), nil
	case []byte:
		return in, nil
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
