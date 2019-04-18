package keyservice

import (
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

			for k, _ := range c.expectedCtx {
				require.Containsf(t, foundCtx, k, "Context does not contain expected key '%s'", k)
			}
			for k, _ := range foundCtx {
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
