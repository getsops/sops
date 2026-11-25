package hcl

import (
	"strings"
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Inputs are in the format emitted by the store, so they must round-trip unchanged.
var roundTripCases = map[string]string{
	"attributes": `region   = "eu-west-1"
count    = 3
ratio    = 0.5
enabled  = true
disabled = false
nothing  = null
`,
	"lists and maps": `ports = [
  80,
  443,
]
tags = {
  env          = "prod"
  "with space" = "yes"
  nested = {
    a = 1
  }
}
empty_list = []
empty_map  = {}
mixed = [
  {
    a = 1
  },
  [
    true,
  ],
]
`,
	"blocks": `variable "name" {
  default = "x"
}

resource "aws_instance" "web" {
  ami = "ami-123"

  ebs_block_device {
    size = 8
  }

  ebs_block_device {
    size = 16
  }
}

empty {
}
`,
	"opaque expressions": `a = var.password
b = "prefix-${var.name}"
c = file("${path.module}/key.pem")
d = var.enabled ? "x" : "y"
e = lookup(var.m, "k", null)
f = [
  var.a,
  "literal",
]
g = {
  k = local.v
}
h = merge(
  var.a,
  {
    k = "v"
  },
)
`,
	"comments": `# leading
a = 1 # inline
# between
b = 2

# before block
block "x" {
  # first in block
  c = 3 # inline in block
  # last in block
}
# trailing
`,
	"escapes": `s = "line1\nline2 \"quoted\" $${not_interp} %%{not_directive}"
`,
}

func TestRoundTrip(t *testing.T) {
	for name, in := range roundTripCases {
		t.Run(name, func(t *testing.T) {
			store := &Store{}
			branches, err := store.LoadPlainFile([]byte(in))
			require.NoError(t, err)
			out, err := store.EmitPlainFile(branches)
			require.NoError(t, err)
			assert.Equal(t, in, string(out))
		})
	}
}

func TestTreeShape(t *testing.T) {
	in := `# top
password = "hunter2" # why
port     = 8080
ratio    = 0.1
ref      = var.x
obj      = { a = "b" }

provider "aws" {
  region = "us-east-1"
}
block {
}
`
	branches, err := (&Store{}).LoadPlainFile([]byte(in))
	require.NoError(t, err)
	assert.Equal(t, sops.TreeBranches{{
		{Key: sops.Comment{Value: "top"}},
		{Key: "password", Value: "hunter2"},
		{Key: sops.Comment{Value: "why", Inline: true}},
		{Key: "port", Value: 8080},
		{Key: "ratio", Value: 0.1},
		{Key: "ref", Value: "${var.x}"},
		{Key: "obj", Value: sops.TreeBranch{{Key: "a", Value: "b"}}},
		{Key: `provider "aws" {`, Value: sops.TreeBranch{{Key: "region", Value: "us-east-1"}}},
		{Key: "block {", Value: sops.TreeBranch(nil)},
	}}, branches)
}

func TestBlockAndObjectAreDistinct(t *testing.T) {
	in := "a {\n  x = 1\n}\n\nb = {\n  x = 1\n}\n"
	store := &Store{}
	branches, err := store.LoadPlainFile([]byte(in))
	require.NoError(t, err)
	out, err := store.EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestCommentInsideLiteralStaysOpaque(t *testing.T) {
	in := "l = [\n  \"a\", # keep me\n  \"b\",\n]\n"
	store := &Store{}
	branches, err := store.LoadPlainFile([]byte(in))
	require.NoError(t, err)
	assert.Equal(t, "${[\n  \"a\", # keep me\n  \"b\",\n]}", branches[0][0].Value)
	out, err := store.EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestBlockCommentsAreNormalized(t *testing.T) {
	branches, err := (&Store{}).LoadPlainFile([]byte("/* one\ntwo */\n// three\na = 1\n"))
	require.NoError(t, err)
	out, err := (&Store{}).EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, "# one\n# two\n# three\na = 1\n", string(out))
}

func TestLiteralLookingLikeExpressionIsRejected(t *testing.T) {
	_, err := (&Store{}).LoadPlainFile([]byte(`a = "$${x}"`))
	assert.ErrorContains(t, err, "cannot be told apart")
}

func TestInvalidInput(t *testing.T) {
	_, err := (&Store{}).LoadPlainFile([]byte(`a = = 1`))
	assert.Error(t, err)
}

func TestEmitRejectsInvalidExpression(t *testing.T) {
	_, err := (&Store{}).EmitPlainFile(sops.TreeBranches{{{Key: "a", Value: "${1 +}"}}})
	assert.Error(t, err)
}

func TestEmitFromOtherFormats(t *testing.T) {
	out, err := (&Store{}).EmitPlainFile(sops.TreeBranches{{
		{Key: "m", Value: sops.TreeBranch{{Key: "not an ident", Value: 1}}},
		{Key: "list", Value: []interface{}{"a", sops.TreeBranch{{Key: "k", Value: 2.5}}}},
	}})
	require.NoError(t, err)
	assert.Equal(t, "m = {\n  \"not an ident\" = 1\n}\nlist = [\n  \"a\",\n  {\n    k = 2.5\n  },\n]\n", string(out))
}

func TestNoMetadata(t *testing.T) {
	_, err := (&Store{}).LoadEncryptedFile([]byte(`hello = 2`))
	assert.Equal(t, sops.MetadataNotFound, err)
}

func TestHasSopsTopLevelKey(t *testing.T) {
	store := &Store{}
	assert.False(t, store.HasSopsTopLevelKey(sops.TreeBranch{{Key: "a", Value: "b"}}))
	assert.True(t, store.HasSopsTopLevelKey(sops.TreeBranch{{Key: "sops", Value: "b"}}))
}

// TestEncryptDecrypt runs a file through the real cipher: every literal and every
// opaque expression must be encrypted, metadata must survive the file, and
// decrypting must give back the input.
func TestEncryptDecrypt(t *testing.T) {
	in := roundTripCases["opaque expressions"] +
		"# leading\nport = 8080\nratio = 0.1\npassword = \"hunter2\"\nplain_unencrypted = \"visible\"\n" +
		"resource \"a\" \"b\" {\n  key = \"secret\"\n}\n"
	in = string(mustFormat(t, in))
	store := &Store{}
	branches, err := store.LoadPlainFile([]byte(in))
	require.NoError(t, err)

	dataKey := []byte("0123456789abcdef0123456789abcdef")
	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			UnencryptedSuffix: "_unencrypted",
			LastModified:      time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
			Version:           "3.0.0",
			KeyGroups: []sops.KeyGroup{{&age.MasterKey{
				Recipient:    "age1yt3tfqlfrwdwx0z0ynwl3jvqmc5eg2zn2c4w5zcdlp0r92x2cwuq9n0qs8",
				EncryptedKey: "-----BEGIN AGE ENCRYPTED FILE-----\nabc\n-----END AGE ENCRYPTED FILE-----\n",
			}}},
		},
	}
	cipher := aes.NewCipher()
	mac, err := tree.Encrypt(dataKey, cipher)
	require.NoError(t, err)
	tree.Metadata.MessageAuthenticationCode, err = cipher.Encrypt(mac, dataKey, tree.Metadata.LastModified.Format(time.RFC3339))
	require.NoError(t, err)

	enc, err := store.EmitEncryptedFile(tree)
	require.NoError(t, err)
	for _, secret := range []string{"hunter2", "secret", "var.password", "var.name", "0.1", "leading"} {
		assert.NotContains(t, strings.ReplaceAll(string(enc), "ENC[", ""), secret)
	}
	assert.Contains(t, string(enc), `plain_unencrypted = "visible"`)
	assert.Contains(t, string(enc), `resource "a" "b" {`)

	loaded, err := store.LoadEncryptedFile(enc)
	require.NoError(t, err)
	assert.Equal(t, tree.Metadata.LastModified, loaded.Metadata.LastModified)
	assert.Equal(t, tree.Metadata.UnencryptedSuffix, loaded.Metadata.UnencryptedSuffix)
	require.Len(t, loaded.Metadata.KeyGroups, 1)
	key := loaded.Metadata.KeyGroups[0][0].(*age.MasterKey)
	assert.Equal(t, "-----BEGIN AGE ENCRYPTED FILE-----\nabc\n-----END AGE ENCRYPTED FILE-----\n", key.EncryptedKey)

	computed, err := loaded.Decrypt(dataKey, cipher)
	require.NoError(t, err)
	assert.Equal(t, mac, computed)
	out, err := store.EmitPlainFile(loaded.Branches)
	require.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func mustFormat(t *testing.T, in string) []byte {
	branches, err := (&Store{}).LoadPlainFile([]byte(in))
	require.NoError(t, err)
	out, err := (&Store{}).EmitPlainFile(branches)
	require.NoError(t, err)
	return out
}

func TestEmitExample(t *testing.T) {
	assert.NotEmpty(t, (&Store{}).EmitExample())
}

func TestNegativeNumbersKeepTheirType(t *testing.T) {
	branches, err := (&Store{}).LoadPlainFile([]byte("i = -5\nf = -0.5\nx = -var.a\n"))
	require.NoError(t, err)
	assert.Equal(t, sops.TreeBranch{
		{Key: "i", Value: -5},
		{Key: "f", Value: -0.5},
		{Key: "x", Value: "${-var.a}"},
	}, branches[0])
	out, err := (&Store{}).EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, "i = -5\nf = -0.5\nx = -var.a\n", string(out))
}

func TestObjectKeyLookingLikeBlock(t *testing.T) {
	in := "a = {\n  \"foo {\" = \"x\"\n}\n"
	store := &Store{}
	branches, err := store.LoadPlainFile([]byte(in))
	require.NoError(t, err)
	out, err := store.EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, in, string(out))
}

// HCL normalizes decoded strings to NFC; anything that is not NFC must be kept as source.
func TestNonNFCStringsAreKeptByteExact(t *testing.T) {
	in := "s = \"café\"\nl = [\n  \"café\",\n]\nm = {\n  \"café\" = 1\n}\nn = \"café\"\n"
	store := &Store{}
	branches, err := store.LoadPlainFile([]byte(in))
	require.NoError(t, err)
	assert.Equal(t, "${\"café\"}", branches[0][0].Value)
	assert.Equal(t, "café", branches[0][3].Value)
	out, err := store.EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, []byte(in), out)

	// The same text written with ASCII escapes must not be decoded either.
	esc := "esc = \"e\\u0301\"\nobj = {\n  \"e\\U00000301\" = 1\n}\n"
	branches, err = store.LoadPlainFile([]byte(esc))
	require.NoError(t, err)
	assert.Equal(t, "${\"e\\u0301\"}", branches[0][0].Value)
	out, err = store.EmitPlainFile(branches)
	require.NoError(t, err)
	assert.Equal(t, esc, string(out))

	_, err = store.LoadPlainFile([]byte("b \"café\" {\n}\n"))
	assert.ErrorContains(t, err, "NFC")
}
