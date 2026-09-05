package keyservice

import (
	"github.com/getsops/sops/v3/azkv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKmsKeyToMasterKey(t *testing.T) {

	cases := []struct {
		description        string
		expectedArn        string
		expectedRole       string
		expectedCtx        map[string]string
		expectedAwsProfile string
	}{
		{
			description:        "empty context",
			expectedArn:        "arn:aws:kms:eu-west-1:123456789012:key/d5c90a06-f824-4628-922b-12424571ed4d",
			expectedRole:       "ExampleRole",
			expectedCtx:        map[string]string{},
			expectedAwsProfile: "",
		},
		{
			description:  "context with one key-value pair",
			expectedArn:  "arn:aws:kms:eu-west-1:123456789012:key/d5c90a06-f824-4628-922b-12424571ed4d",
			expectedRole: "",
			expectedCtx: map[string]string{
				"firstKey": "first value",
			},
			expectedAwsProfile: "ExampleProfile",
		},
		{
			description:  "context with three key-value pairs",
			expectedArn:  "arn:aws:kms:eu-west-1:123456789012:key/d5c90a06-f824-4628-922b-12424571ed4d",
			expectedRole: "",
			expectedCtx: map[string]string{
				"firstKey":  "first value",
				"secondKey": "second value",
				"thirdKey":  "third value",
			},
			expectedAwsProfile: "",
		},
	}

	for _, c := range cases {

		t.Run(c.description, func(t *testing.T) {

			inputCtx := make(map[string]string)
			for k, v := range c.expectedCtx {
				inputCtx[k] = v
			}

			key := &KmsKey{
				Arn:        c.expectedArn,
				Role:       c.expectedRole,
				Context:    inputCtx,
				AwsProfile: c.expectedAwsProfile,
			}

			masterKey := kmsKeyToMasterKey(key)
			foundCtx := masterKey.EncryptionContext

			for k := range c.expectedCtx {
				require.Containsf(t, foundCtx, k, "Context does not contain expected key '%s'", k)
			}
			for k := range foundCtx {
				require.Containsf(t, c.expectedCtx, k, "Context contains an unexpected key '%s' which cannot be found from expected map", k)
			}
			for k, expected := range c.expectedCtx {
				foundVal := *foundCtx[k]
				assert.Equalf(t, expected, foundVal, "Context key '%s' value '%s' does not match expected value '%s'", k, foundVal, expected)
			}
			assert.Equalf(t, c.expectedArn, masterKey.Arn, "Expected ARN to be '%s', but found '%s'", c.expectedArn, masterKey.Arn)
			assert.Equalf(t, c.expectedRole, masterKey.Role, "Expected Role to be '%s', but found '%s'", c.expectedRole, masterKey.Role)
			assert.Equalf(t, c.expectedAwsProfile, masterKey.AwsProfile, "Expected AWS profile to be '%s', but found '%s'", c.expectedAwsProfile, masterKey.AwsProfile)
		})
	}
}

// Azure KV tests for env-driven skip URI validation affecting client options.
func TestAzureKeyVaultClientOptionsAppliedOnEncryptDecrypt(t *testing.T) {
	// ensure we don't perform network calls
	testHookSkipAzureNetwork = true
	t.Cleanup(func() {
		testHookSkipAzureNetwork = false
		testHookCaptureAzureKey = nil
	})

	for _, tt := range []struct {
		name            string
		envValue        string
		expectOptionSet bool
	}{
		{name: "encrypt sets option when env true", envValue: "true", expectOptionSet: true},
		{name: "encrypt leaves option nil when env false", envValue: "false", expectOptionSet: false},
		{name: "encrypt leaves option nil when env invalid", envValue: "not-a-bool", expectOptionSet: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(azkv.SopsAzureKeyvaultSkipUriVerificationEnv, tt.envValue)

			captured := []*azkv.MasterKey{}
			testHookCaptureAzureKey = func(mk *azkv.MasterKey) { captured = append(captured, mk) }
			server := &Server{}
			key := &AzureKeyVaultKey{VaultUrl: "https://vault.example", Name: "keyname", Version: "v1"}
			_, err := server.encryptWithAzureKeyVault(key, []byte("secret"))
			require.NoError(t, err)
			require.Len(t, captured, 1)

			co := captured[0].EffectiveClientOptions()
			if tt.expectOptionSet {
				require.NotNil(t, co)
				assert.True(t, co.DisableChallengeResourceVerification)
			} else {
				assert.Nil(t, co)
			}
		})
	}

	for _, tt := range []struct {
		name            string
		envValue        string
		expectOptionSet bool
	}{
		{name: "decrypt sets option when env true", envValue: "1", expectOptionSet: true},
		{name: "decrypt leaves option nil when env false", envValue: "0", expectOptionSet: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(azkv.SopsAzureKeyvaultSkipUriVerificationEnv, tt.envValue)

			captured := []*azkv.MasterKey{}
			testHookCaptureAzureKey = func(mk *azkv.MasterKey) { captured = append(captured, mk) }
			server := &Server{}
			key := &AzureKeyVaultKey{VaultUrl: "https://vault.example", Name: "keyname", Version: "v1"}
			_, err := server.decryptWithAzureKeyVault(key, []byte("c2VjcmV0"))
			require.NoError(t, err)
			require.Len(t, captured, 1)

			co := captured[0].EffectiveClientOptions()
			if tt.expectOptionSet {
				require.NotNil(t, co)
				assert.True(t, co.DisableChallengeResourceVerification)
			} else {
				assert.Nil(t, co)
			}
		})
	}
}
