package dotenv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops/v3"
)

var ORIGINAL_PLAIN = []byte(strings.TrimLeft(`
#Comment
#    Trimmed comment 
UNQUOTED=value
UNQUOTED_ESCAPED_NEWLINE=escaped\nnewline
UNQUOTED_WHITESPACE= trimmed whitespace 
SINGLEQUOTED='value'
SINGLEQUOTED_NEWLINE='real
newline'
SINGLEQUOTED_ESCAPED_NEWLINE='escaped\nnewline'
SINGLEQUOTED_ESCAPED_QUOTE='escaped\'quote'
SINGLEQUOTED_WHITESPACE=' untrimmed whitespace '
DOUBLEQUOTED="value"
DOUBLEQUOTED_NEWLINE="real
newline"
DOUBLEQUOTED_ESCAPED_NEWLINE="real\nnewline"
DOUBLEQUOTED_ESCAPED_QUOTE="escaped\"quote"
DOUBLEQUOTED_WHITESPACE=" untrimmed whitespace "
`, "\n"))

var EMITTED_PLAIN = []byte(strings.TrimLeft(`
# Comment
# Trimmed comment
UNQUOTED='value'
UNQUOTED_ESCAPED_NEWLINE='escaped\nnewline'
UNQUOTED_WHITESPACE='trimmed whitespace'
SINGLEQUOTED='value'
SINGLEQUOTED_NEWLINE='real
newline'
SINGLEQUOTED_ESCAPED_NEWLINE='escaped\nnewline'
SINGLEQUOTED_ESCAPED_QUOTE='escaped\'quote'
SINGLEQUOTED_WHITESPACE=' untrimmed whitespace '
DOUBLEQUOTED='value'
DOUBLEQUOTED_NEWLINE='real
newline'
DOUBLEQUOTED_ESCAPED_NEWLINE='real
newline'
DOUBLEQUOTED_ESCAPED_QUOTE='escaped"quote'
DOUBLEQUOTED_WHITESPACE=' untrimmed whitespace '
`, "\n"))

var BRANCH = sops.TreeBranch{
	sops.TreeItem{Key: sops.Comment{"Comment"}, Value: nil},
	sops.TreeItem{Key: sops.Comment{"Trimmed comment"}, Value: nil},
	sops.TreeItem{Key: "UNQUOTED", Value: "value"},
	sops.TreeItem{Key: "UNQUOTED_ESCAPED_NEWLINE", Value: "escaped\\nnewline"},
	sops.TreeItem{Key: "UNQUOTED_WHITESPACE", Value: "trimmed whitespace"},
	sops.TreeItem{Key: "SINGLEQUOTED", Value: "value"},
	sops.TreeItem{Key: "SINGLEQUOTED_NEWLINE", Value: "real\nnewline"},
	sops.TreeItem{Key: "SINGLEQUOTED_ESCAPED_NEWLINE", Value: "escaped\\nnewline"},
	sops.TreeItem{Key: "SINGLEQUOTED_ESCAPED_QUOTE", Value: "escaped'quote"},
	sops.TreeItem{Key: "SINGLEQUOTED_WHITESPACE", Value: " untrimmed whitespace "},
	sops.TreeItem{Key: "DOUBLEQUOTED", Value: "value"},
	sops.TreeItem{Key: "DOUBLEQUOTED_NEWLINE", Value: "real\nnewline"},
	sops.TreeItem{Key: "DOUBLEQUOTED_ESCAPED_NEWLINE", Value: "real\nnewline"},
	sops.TreeItem{Key: "DOUBLEQUOTED_ESCAPED_QUOTE", Value: "escaped\"quote"},
	sops.TreeItem{Key: "DOUBLEQUOTED_WHITESPACE", Value: " untrimmed whitespace "},
}

func TestLoadPlainFile(t *testing.T) {
	branches, err := (&Store{}).LoadPlainFile(ORIGINAL_PLAIN)
	assert.Nil(t, err)
	assert.Equal(t, BRANCH, branches[0])
}

func TestInvalidKeyError(t *testing.T) {
	_, err := (&Store{}).LoadPlainFile([]byte("INVALID KEY=irrelevant value"))
	assert.Equal(t, err.Error(), "invalid dotenv key: \"INVALID KEY\"")
}

func TestEmitPlainFile(t *testing.T) {
	branches := sops.TreeBranches{
		BRANCH,
	}
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, EMITTED_PLAIN, bytes)
}

func TestEmitValueString(t *testing.T) {
	bytes, err := (&Store{}).EmitValue("hello")
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), bytes)
}

func TestEmitValueNonstring(t *testing.T) {
	_, err := (&Store{}).EmitValue(BRANCH)
	assert.NotNil(t, err)
}