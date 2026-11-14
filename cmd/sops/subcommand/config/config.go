// Package config implements the `sops config <file>` subcommand, which
// prints the .sops.yaml rules that would apply to a given file path.
package config

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/config"
)

// SchemaVersion is the version of the JSON schema emitted by this command.
// See docs/superpowers/specs/2026-05-22-sops-config-subcommand-design.md
// for the schema evolution policy.
const SchemaVersion = 1

// Output is the top-level JSON shape printed by `sops config <file>`.
// CreationRules and DestinationRules are guaranteed non-nil — even when
// empty they marshal to "[]", never "null". Implementation MUST use
// make([]T, 0) to construct them.
type Output struct {
	SchemaVersion    int                   `json:"schema_version"`
	ConfigPath       string                `json:"config_path"`
	FilePath         string                `json:"file_path"`
	CreationRules    []CreationRuleView    `json:"creation_rules"`
	DestinationRules []DestinationRuleView `json:"destination_rules"`
}

// CreationRuleView is the public JSON representation of a matched creation_rule.
//
// Recipient fields (KMS, GCPKMS, AzureKeyVault, HCKms, HCVaultTransitURI, Age,
// PGP, KeyGroups) use omitempty because sops's parser uses either KeyGroups OR
// the flat recipient fields, never both (see config/config.go:371-381). Echoing
// only what was written in .sops.yaml keeps the output unambiguous.
//
// Non-recipient fields are always emitted for schema stability.
type CreationRuleView struct {
	RuleIndex int    `json:"rule_index"`
	PathRegex string `json:"path_regex"`

	KMS               []KmsKeyView     `json:"kms,omitempty"`
	GCPKMS            []GcpKmsKeyView  `json:"gcp_kms,omitempty"`
	AzureKeyVault     []AzureKVKeyView `json:"azure_keyvault,omitempty"`
	HCKms             []HCKmsKeyView   `json:"hckms,omitempty"`
	HCVaultTransitURI []string         `json:"hc_vault_transit_uri,omitempty"`
	Age               []string         `json:"age,omitempty"`
	PGP               []string         `json:"pgp,omitempty"`
	KeyGroups         []KeyGroupView   `json:"key_groups,omitempty"`

	ShamirThreshold         int    `json:"shamir_threshold"`
	UnencryptedSuffix       string `json:"unencrypted_suffix"`
	EncryptedSuffix         string `json:"encrypted_suffix"`
	UnencryptedRegex        string `json:"unencrypted_regex"`
	EncryptedRegex          string `json:"encrypted_regex"`
	UnencryptedCommentRegex string `json:"unencrypted_comment_regex"`
	EncryptedCommentRegex   string `json:"encrypted_comment_regex"`
	MACOnlyEncrypted        bool   `json:"mac_only_encrypted"`
}

// DestinationRuleView is the public JSON representation of a matched destination_rule.
//
// Destination and RecreationRule are pointer types with omitempty because Go's
// encoding/json cannot treat value-typed structs as empty regardless of their
// field values — pointer + omitempty is the only way to get true optionality.
//
// RecreationRule reuses CreationRuleView. In the nested recreation_rule
// context, the inherited RuleIndex and PathRegex fields have no meaning and
// will marshal as their zero values ("rule_index": 0, "path_regex": "");
// consumers should ignore those two fields when reading recreation_rule.
type DestinationRuleView struct {
	RuleIndex      int               `json:"rule_index"`
	PathRegex      string            `json:"path_regex"`
	Destination    *DestinationView  `json:"destination,omitempty"`
	OmitExtensions bool              `json:"omit_extensions"`
	RecreationRule *CreationRuleView `json:"recreation_rule,omitempty"`
}

// DestinationView is the discriminated-union JSON form of a publish target.
// Type is one of "s3", "gcs", "vault"; other fields are populated according
// to Type.
type DestinationView struct {
	Type        string `json:"type"`
	Bucket      string `json:"bucket,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Address     string `json:"address,omitempty"`
	Path        string `json:"path,omitempty"`
	KVMountName string `json:"kv_mount_name,omitempty"`
	KVVersion   int    `json:"kv_version,omitempty"`
}

// KmsKeyView mirrors the internal config kmsKey struct with JSON tags.
type KmsKeyView struct {
	Arn        string             `json:"arn"`
	Role       string             `json:"role,omitempty"`
	Context    map[string]*string `json:"context,omitempty"`
	AwsProfile string             `json:"aws_profile,omitempty"`
}

// GcpKmsKeyView mirrors the internal config gcpKmsKey struct.
type GcpKmsKeyView struct {
	ResourceID string `json:"resource_id"`
}

// AzureKVKeyView mirrors the internal config azureKVKey struct.
type AzureKVKeyView struct {
	VaultURL string `json:"vaultUrl"`
	Key      string `json:"key"`
	Version  string `json:"version,omitempty"`
}

// HCKmsKeyView mirrors the internal config hckmsKey struct.
type HCKmsKeyView struct {
	KeyID string `json:"key_id"`
}

// KeyGroupView mirrors the internal config keyGroup struct for use inside
// CreationRuleView.KeyGroups.
type KeyGroupView struct {
	Merge          []KeyGroupView   `json:"merge,omitempty"`
	KMS            []KmsKeyView     `json:"kms,omitempty"`
	GCPKMS         []GcpKmsKeyView  `json:"gcp_kms,omitempty"`
	HCKms          []HCKmsKeyView   `json:"hckms,omitempty"`
	AzureKeyVault  []AzureKVKeyView `json:"azure_keyvault,omitempty"`
	HCVaultTransit []string         `json:"hc_vault,omitempty"`
	Age            []string         `json:"age,omitempty"`
	PGP            []string         `json:"pgp,omitempty"`
}

// buildCreationRuleView converts a config.CreationRuleMatch into the
// public JSON view. All recipient-list errors propagate (e.g., malformed
// interface{} types).
//
// Sops's parser uses key_groups XOR flat recipient fields (never both).
// When key_groups is populated the flat fields are silently ignored by the
// encryption pipeline, so the view only populates the branch that is active.
func buildCreationRuleView(m *config.CreationRuleMatch) (*CreationRuleView, error) {
	view := &CreationRuleView{
		RuleIndex:               m.RuleIndex,
		PathRegex:               m.PathRegex(),
		ShamirThreshold:         m.ShamirThreshold(),
		UnencryptedSuffix:       m.UnencryptedSuffix(),
		EncryptedSuffix:         m.EncryptedSuffix(),
		UnencryptedRegex:        m.UnencryptedRegex(),
		EncryptedRegex:          m.EncryptedRegex(),
		UnencryptedCommentRegex: m.UnencryptedCommentRegex(),
		EncryptedCommentRegex:   m.EncryptedCommentRegex(),
		MACOnlyEncrypted:        m.MACOnlyEncrypted(),
	}

	if len(m.KeyGroups()) > 0 {
		// key_groups is active; flat fields are dead. Only populate KeyGroups.
		for _, g := range m.KeyGroups() {
			view.KeyGroups = append(view.KeyGroups, convertKeyGroupView(g))
		}
	} else {
		kmsEntries, err := m.KMSEntries()
		if err != nil {
			return nil, fmt.Errorf("invalid kms: %w", err)
		}
		for _, e := range kmsEntries {
			view.KMS = append(view.KMS, KmsKeyView{
				Arn:        e.Arn,
				Role:       e.Role,
				Context:    e.Context,
				AwsProfile: e.AwsProfile,
			})
		}

		age, err := m.AgeRecipients()
		if err != nil {
			return nil, fmt.Errorf("invalid age: %w", err)
		}
		view.Age = age

		pgp, err := m.PGPFingerprints()
		if err != nil {
			return nil, fmt.Errorf("invalid pgp: %w", err)
		}
		view.PGP = pgp

		gcp, err := m.GCPKMSResourceIDs()
		if err != nil {
			return nil, fmt.Errorf("invalid gcp_kms: %w", err)
		}
		for _, id := range gcp {
			view.GCPKMS = append(view.GCPKMS, GcpKmsKeyView{ResourceID: id})
		}

		azkv, err := m.AzureKeyVaults()
		if err != nil {
			return nil, fmt.Errorf("invalid azure_keyvault: %w", err)
		}
		for _, u := range azkv {
			view.AzureKeyVault = append(view.AzureKeyVault, parseAzureKeyVaultURL(u))
		}

		vaults, err := m.HCVaultTransitURIs()
		if err != nil {
			return nil, fmt.Errorf("invalid hc_vault_transit_uri: %w", err)
		}
		view.HCVaultTransitURI = vaults

		for _, kid := range m.HCKmsKeyIDs() {
			view.HCKms = append(view.HCKms, HCKmsKeyView{KeyID: kid})
		}
	}

	return view, nil
}

// buildDestinationRuleView converts a config.DestinationRuleMatch into the
// public JSON view. The destination type is determined by which set of
// destinationRule fields is populated (S3, GCS, or Vault). The nested
// recreation_rule is converted via buildCreationRuleView.
func buildDestinationRuleView(m *config.DestinationRuleMatch) (*DestinationRuleView, error) {
	view := &DestinationRuleView{
		RuleIndex:      m.RuleIndex,
		PathRegex:      m.PathRegex(),
		OmitExtensions: m.OmitExtensions(),
	}

	if bucket, prefix, ok := m.S3(); ok {
		view.Destination = &DestinationView{Type: "s3", Bucket: bucket, Prefix: prefix}
	} else if bucket, prefix, ok := m.GCS(); ok {
		view.Destination = &DestinationView{Type: "gcs", Bucket: bucket, Prefix: prefix}
	} else if addr, path, kvMount, kvVer, ok := m.Vault(); ok {
		view.Destination = &DestinationView{
			Type:        "vault",
			Address:     addr,
			Path:        path,
			KVMountName: kvMount,
			KVVersion:   kvVer,
		}
	}

	if rr := m.RecreationRule(); rr != nil {
		recView, err := buildCreationRuleView(rr)
		if err != nil {
			return nil, fmt.Errorf("invalid recreation_rule: %w", err)
		}
		view.RecreationRule = recView
	}

	return view, nil
}

// Opts describes the inputs to Run.
type Opts struct {
	ConfigPath   string // absolute path to .sops.yaml
	FilePath     string // absolute path to the query file
	RequireMatch bool   // if true, exit non-zero when no rule matches
}

// Run executes the `sops config <file>` subcommand and returns the Output
// to serialize, the exit code, and any error.
//
// Output is non-nil in exactly two cases:
//   - Success (err == nil, exitCode == 0).
//   - --require-match with no matches (err != nil, exitCode == codes.NoRulesMatched).
//
// In the require-match case the Output is still valid JSON and SHOULD be
// printed to stdout — consumers may want the structured data even when the
// overall run is treated as a failure.
//
// Output is nil on all other errors (config load failure, view-conversion
// failure, etc.). The caller should NOT print anything in those cases.
func Run(opts Opts) (*Output, int, error) {
	mr, err := config.MatchRulesForFile(opts.ConfigPath, opts.FilePath)
	if err != nil {
		return nil, codes.ErrorReadingConfig, err
	}

	output := &Output{
		SchemaVersion: SchemaVersion,
		ConfigPath:    mr.ConfigPath,
		FilePath:      mr.FilePath,
		// make([]T, 0) so the slices marshal to "[]" rather than "null".
		CreationRules:    make([]CreationRuleView, 0),
		DestinationRules: make([]DestinationRuleView, 0),
	}

	if mr.CreationRule != nil {
		view, err := buildCreationRuleView(mr.CreationRule)
		if err != nil {
			return nil, codes.ErrorReadingConfig, fmt.Errorf("convert creation_rule: %w", err)
		}
		output.CreationRules = append(output.CreationRules, *view)
	}

	if mr.DestinationRule != nil {
		view, err := buildDestinationRuleView(mr.DestinationRule)
		if err != nil {
			return nil, codes.ErrorReadingConfig, fmt.Errorf("convert destination_rule: %w", err)
		}
		output.DestinationRules = append(output.DestinationRules, *view)
	}

	if opts.RequireMatch && len(output.CreationRules) == 0 && len(output.DestinationRules) == 0 {
		return output, codes.NoRulesMatched, fmt.Errorf("no matching rules found for %q in %s", opts.FilePath, opts.ConfigPath)
	}

	return output, 0, nil
}

// Marshal returns the JSON representation of the Output suitable for stdout.
func Marshal(output *Output) ([]byte, error) {
	return json.MarshalIndent(output, "", "  ")
}

// parseAzureKeyVaultURL splits an Azure Key Vault key URL into its
// constituent parts. Mirrors the parser at azkv/keysource.go:96. If the
// URL doesn't match the expected shape, the whole URL is returned in
// VaultURL with Key and Version empty — preserves the input for debugging.
func parseAzureKeyVaultURL(url string) AzureKVKeyView {
	re := regexp.MustCompile(`^(https://[^/]+)/keys/([^/]+)(/[^/]*)?$`)
	parts := re.FindStringSubmatch(url)
	if len(parts) < 3 {
		return AzureKVKeyView{VaultURL: url}
	}
	view := AzureKVKeyView{VaultURL: parts[1], Key: parts[2]}
	if len(parts[3]) > 1 {
		view.Version = parts[3][1:] // strip the leading "/"
	}
	return view
}

// convertKeyGroupView translates a config.KeyGroupEntry to its JSON view form.
// Nested merge groups recurse.
func convertKeyGroupView(g config.KeyGroupEntry) KeyGroupView {
	out := KeyGroupView{
		HCVaultTransit: g.HCVaultTransit,
		Age:            g.Age,
		PGP:            g.PGP,
	}
	for _, k := range g.KMS {
		out.KMS = append(out.KMS, KmsKeyView{
			Arn:        k.Arn,
			Role:       k.Role,
			Context:    k.Context,
			AwsProfile: k.AwsProfile,
		})
	}
	for _, id := range g.GCPKMS {
		out.GCPKMS = append(out.GCPKMS, GcpKmsKeyView{ResourceID: id})
	}
	for _, id := range g.HCKms {
		out.HCKms = append(out.HCKms, HCKmsKeyView{KeyID: id})
	}
	for _, kv := range g.AzureKeyVault {
		out.AzureKeyVault = append(out.AzureKeyVault, AzureKVKeyView{
			VaultURL: kv.VaultURL,
			Key:      kv.Key,
			Version:  kv.Version,
		})
	}
	for _, sub := range g.Merge {
		out.Merge = append(out.Merge, convertKeyGroupView(sub))
	}
	return out
}
