package kms

import (
	"context"
	"encoding/base64"
	"fmt"
	logger "log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	// testKMSServerURL is the URL of the AWS KMS server running in Docker.
	// It is loaded by TestMain.
	testKMSServerURL string
	// testKMSARN is the ARN on the test AWS KMS server. It is loaded
	// by TestMain.
	testKMSARN string
)

const (
	// dummyARN is a dummy AWS ARN which passes validation.
	dummyARN = "arn:aws:kms:us-west-2:107501996527:key/612d5f0p-p1l3-45e6-aca6-a5b005693a48"
	// testLocalKMSImage is a container image repository reference to a mock
	// version of AWS' Key Management Service.
	// Ref: https://github.com/nsmithuk/local-kms
	testLocalKMSImage = "docker.io/nsmithuk/local-kms"
	// testLocalKMSImage is the container image tag to use.
	testLocalKMSTag = "3.11.1"
)

// TestMain initializes an AWS KMS server using Docker, writes the HTTP address
// to testKMSServerURL, tries to generate a key for encryption-decryption using a
// backoff retry approach, and then sets testKMSARN to the ID of the generated key.
// It continues to run all the tests, which can make use of the various `test*`
// variables.
func TestMain(m *testing.M) {
	// Uses a sensible default on Windows (TCP/HTTP) and Linux/MacOS (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatalf("could not connect to docker: %s", err)
	}

	// Pull the image, create a container based on it, and run it
	// resource, err := pool.Run("nsmithuk/local-kms", testLocalKMSVersion, []string{})
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   testLocalKMSImage,
		Tag:          testLocalKMSTag,
		ExposedPorts: []string{"8080"},
	})
	if err != nil {
		logger.Fatalf("could not start resource: %s", err)
	}

	purgeResource := func() {
		if err := pool.Purge(resource); err != nil {
			logger.Printf("could not purge resource: %s", err)
		}
	}

	testKMSServerURL = fmt.Sprintf("http://127.0.0.1:%v", resource.GetPort("8080/tcp"))
	masterKey := createTestMasterKey(dummyARN)

	kmsClient, err := createTestKMSClient(masterKey)
	if err != nil {
		purgeResource()
		logger.Fatalf("could not create session: %s", err)
	}

	var key *kms.CreateKeyOutput
	if err := pool.Retry(func() error {
		key, err = kmsClient.CreateKey(context.TODO(), &kms.CreateKeyInput{})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		purgeResource()
		logger.Fatalf("could not create key: %s", err)
	}

	if key.KeyMetadata.Arn != nil {
		testKMSARN = *key.KeyMetadata.Arn
	} else {
		purgeResource()
		logger.Fatalf("could not set arn")
	}

	// Run the tests, but only if we succeeded in setting up the AWS KMS server.
	var code int
	if err == nil {
		code = m.Run()
	}

	// This can't be deferred, as os.Exit simply does not care
	if err := pool.Purge(resource); err != nil {
		logger.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestNewMasterKey(t *testing.T) {
	var (
		dummyRole              = "a-role"
		dummyEncryptionContext = map[string]*string{
			"foo": aws.String("bar"),
		}
	)
	key := NewMasterKey(dummyARN, dummyRole, dummyEncryptionContext)
	assert.Equal(t, dummyARN, key.Arn)
	assert.Equal(t, dummyRole, key.Role)
	assert.Equal(t, dummyEncryptionContext, key.EncryptionContext)
	assert.NotNil(t, key.CreationDate)
}

func TestNewMasterKeyWithProfile(t *testing.T) {
	var (
		dummyRole              = "a-role"
		dummyEncryptionContext = map[string]*string{
			"foo": aws.String("bar"),
		}
		dummyProfile = "a-profile"
	)
	key := NewMasterKeyWithProfile(dummyARN, dummyRole, dummyEncryptionContext, dummyProfile)
	assert.Equal(t, dummyARN, key.Arn)
	assert.Equal(t, dummyRole, key.Role)
	assert.Equal(t, dummyEncryptionContext, key.EncryptionContext)
	assert.Equal(t, dummyProfile, key.AwsProfile)
	assert.NotNil(t, key.CreationDate)
}

func TestNewMasterKeyFromArn(t *testing.T) {
	t.Run("arn", func(t *testing.T) {
		var (
			dummyEncryptionContext = map[string]*string{
				"foo": aws.String("bar"),
			}
			dummyProfile = "a-profile"
		)
		key := NewMasterKeyFromArn(dummyARN, dummyEncryptionContext, dummyProfile)
		assert.Equal(t, dummyARN, key.Arn)
		assert.Equal(t, dummyEncryptionContext, key.EncryptionContext)
		assert.Equal(t, dummyProfile, key.AwsProfile)
		assert.Empty(t, key.Role)
		assert.NotNil(t, key.CreationDate)
	})

	t.Run("arn with spaces", func(t *testing.T) {
		key := NewMasterKeyFromArn(" arn:aws:kms:us-west-2 :107501996527:key/612d5f 0p-p1l3-45e6-aca6-a5b00569 3a48 ", nil, "")
		assert.Equal(t, "arn:aws:kms:us-west-2:107501996527:key/612d5f0p-p1l3-45e6-aca6-a5b005693a48", key.Arn)
	})

	t.Run("arn with role", func(t *testing.T) {
		key := NewMasterKeyFromArn("arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500+arn:aws:iam::927034868273:role/sops-dev-xyz", nil, "")
		assert.Equal(t, "arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500", key.Arn)
		assert.Equal(t, "arn:aws:iam::927034868273:role/sops-dev-xyz", key.Role)
	})
}

func TestMasterKeysFromArnString(t *testing.T) {
	s := "arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e+arn:aws:iam::927034868273:role/sops-dev, arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
	ks := MasterKeysFromArnString(s, nil, "foo")
	k1 := ks[0]
	k2 := ks[1]

	expectedArn1 := "arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e"
	expectedRole1 := "arn:aws:iam::927034868273:role/sops-dev"
	assert.Equal(t, expectedArn1, k1.Arn)
	assert.Equal(t, expectedRole1, k1.Role)

	expectedArn2 := "arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
	assert.Equal(t, expectedArn2, k2.Arn)
	assert.Empty(t, k2.Role)
}

func TestParseKMSContext(t *testing.T) {
	value1 := "value1"
	value2 := "value2"
	// map from YAML
	var yamlmap = map[interface{}]interface{}{
		"key1": value1,
		"key2": value2,
	}
	assert.Equal(t, ParseKMSContext(yamlmap), map[string]*string{
		"key1": &value1,
		"key2": &value2,
	})
	assert.Nil(t, ParseKMSContext(map[interface{}]interface{}{}))
	assert.Nil(t, ParseKMSContext(map[interface{}]interface{}{
		"key1": 1,
	}))
	assert.Nil(t, ParseKMSContext(map[interface{}]interface{}{
		1: "value",
	}))
	// map from JSON
	var jsonmap = map[string]interface{}{
		"key1": value1,
		"key2": value2,
	}
	assert.Equal(t, ParseKMSContext(jsonmap), map[string]*string{
		"key1": &value1,
		"key2": &value2,
	})
	assert.Nil(t, ParseKMSContext(map[string]interface{}{}))
	assert.Nil(t, ParseKMSContext(map[string]interface{}{
		"key1": 1,
	}))
	// sops 2.0.x formatted encryption context as a comma-separated list of key:value pairs
	assert.Equal(t, ParseKMSContext("key1:value1,key2:value2"), map[string]*string{
		"key1": &value1,
		"key2": &value2,
	})
	assert.Equal(t, ParseKMSContext("key1:value1"), map[string]*string{
		"key1": &value1,
	})
	assert.Nil(t, ParseKMSContext("key1,key2:value2"))
	assert.Nil(t, ParseKMSContext("key1"))
}

func TestCreds_ApplyToMasterKey(t *testing.T) {
	creds := NewCredentialsProvider(credentials.NewStaticCredentialsProvider("", "", ""))
	key := &MasterKey{}
	creds.ApplyToMasterKey(key)
	assert.Equal(t, creds.provider, key.credentialsProvider)
}

func TestMasterKey_Encrypt(t *testing.T) {
	t.Run("encrypt", func(t *testing.T) {
		key := createTestMasterKey(testKMSARN)
		dataKey := []byte("UFO sightings")
		assert.NoError(t, key.Encrypt(dataKey))
		assert.NotEmpty(t, key.EncryptedKey)

		kmsClient, err := createTestKMSClient(key)
		assert.NoError(t, err)

		k, err := base64.StdEncoding.DecodeString(key.EncryptedKey)
		assert.NoError(t, err)

		input := &kms.DecryptInput{
			CiphertextBlob:    k,
			EncryptionContext: stringPointerToStringMap(key.EncryptionContext),
		}
		decrypted, err := kmsClient.Decrypt(context.TODO(), input)
		assert.NoError(t, err)
		assert.Equal(t, dataKey, decrypted.Plaintext)
	})

	t.Run("encrypt error", func(t *testing.T) {
		// Valid ARN but invalid for test server.
		key := createTestMasterKey(dummyARN)
		err := key.Encrypt([]byte("UFO sightings"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to encrypt sops data key with AWS KMS")
		assert.Empty(t, key.EncryptedKey)
	})

	t.Run("config error", func(t *testing.T) {
		key := createTestMasterKey("arn:gcp:kms:antartica-north-2::key/45e6-aca6-a5b005693a48")
		err := key.Encrypt([]byte(""))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no valid ARN found")
		assert.Empty(t, key.EncryptedKey)
	})
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key := createTestMasterKey(testKMSARN)
	assert.NoError(t, key.EncryptIfNeeded([]byte("data")))

	encryptedKey := key.EncryptedKey
	assert.NotEmpty(t, encryptedKey)

	assert.NoError(t, key.EncryptIfNeeded([]byte("some other data")))
	assert.Equal(t, encryptedKey, key.EncryptedKey)
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: "some key"}
	assert.EqualValues(t, key.EncryptedKey, key.EncryptedDataKey())
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	key := &MasterKey{}
	data := []byte("some data")
	key.SetEncryptedDataKey(data)
	assert.EqualValues(t, data, key.EncryptedKey)
}

func TestMasterKey_Decrypt(t *testing.T) {
	t.Run("decrypt", func(t *testing.T) {
		key := createTestMasterKey(testKMSARN)
		kmsClient, err := createTestKMSClient(key)
		assert.NoError(t, err)

		dataKey := []byte("it's always DNS")
		out, err := kmsClient.Encrypt(context.TODO(), &kms.EncryptInput{
			Plaintext: dataKey, KeyId: &key.Arn, EncryptionContext: stringPointerToStringMap(key.EncryptionContext),
		})
		assert.NoError(t, err)

		key.EncryptedKey = base64.StdEncoding.EncodeToString(out.CiphertextBlob)
		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.Equal(t, dataKey, got)
	})

	t.Run("data key error", func(t *testing.T) {
		key := createTestMasterKey(testKMSARN)
		key.EncryptedKey = "invalid"
		got, err := key.Decrypt()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error base64-decoding encrypted data key")
		assert.Nil(t, got)
	})

	t.Run("decrypt error", func(t *testing.T) {
		// Valid ARN but invalid for test server.
		key := createTestMasterKey(dummyARN)
		key.EncryptedKey = base64.StdEncoding.EncodeToString([]byte("invalid"))
		got, err := key.Decrypt()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to decrypt sops data key with AWS KMS")
		assert.Nil(t, got)
	})

	t.Run("config error", func(t *testing.T) {
		key := createTestMasterKey("arn:gcp:kms:antartica-north-2::key/45e6-aca6-a5b005693a48")
		got, err := key.Decrypt()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no valid ARN found")
		assert.Nil(t, got)
	})
}

func TestMasterKey_EncryptDecrypt_RoundTrip(t *testing.T) {
	dataKey := []byte("the wheels on the bus go round and round")

	encryptKey := createTestMasterKey(testKMSARN)
	assert.NoError(t, encryptKey.Encrypt(dataKey))
	assert.NotEmpty(t, encryptKey.EncryptedKey)

	decryptKey := createTestMasterKey(testKMSARN)
	decryptKey.EncryptedKey = encryptKey.EncryptedKey

	decryptedData, err := decryptKey.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, dataKey, decryptedData)
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := NewMasterKeyFromArn(dummyARN, nil, "")
	assert.False(t, key.NeedsRotation())

	key.CreationDate = key.CreationDate.Add(-(kmsTTL + time.Second))
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKeyFromArn(dummyARN, nil, "")
	assert.Equal(t, dummyARN, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	value1 := "value1"
	value2 := "value2"
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		Arn:          "foo",
		Role:         "bar",
		EncryptedKey: "this is encrypted",
		EncryptionContext: map[string]*string{
			"key1": &value1,
			"key2": &value2,
		},
	}
	assert.Equal(t, map[string]interface{}{
		"arn":        "foo",
		"role":       "bar",
		"enc":        "this is encrypted",
		"created_at": "2016-10-31T10:00:00Z",
		"context": map[string]string{
			"key1": value1,
			"key2": value2,
		},
	}, key.ToMap())
}

func TestMasterKey_createKMSConfig(t *testing.T) {
	tests := []struct {
		name       string
		key        MasterKey
		envFunc    func(t *testing.T)
		assertFunc func(t *testing.T, cfg *aws.Config, err error)
		fallback   bool
	}{
		{
			name: "valid config with credentials provider",
			key: MasterKey{
				credentialsProvider: credentials.NewStaticCredentialsProvider("test-id", "test-secret", "test-token"),
				Arn:                 "arn:aws:kms:us-west-2:107501996527:key/612d5f0p-p1l3-45e6-aca6-a5b005693a48",
			},
			assertFunc: func(t *testing.T, cfg *aws.Config, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "us-west-2", cfg.Region)

				creds, err := cfg.Credentials.Retrieve(context.TODO())
				assert.NoError(t, err)
				assert.Equal(t, "test-id", creds.AccessKeyID)
				assert.Equal(t, "test-secret", creds.SecretAccessKey)
				assert.Equal(t, "test-token", creds.SessionToken)
			},
		},
		{
			name: "valid config with profile",
			key: MasterKey{
				AwsProfile: "test-profile",
				Arn:        "arn:aws:kms:us-west-2:107501996527:key/612d5f0p-p1l3-45e6-aca6-a5b005693a48",
			},
			envFunc: func(t *testing.T) {
				credentialsFile := filepath.Join(t.TempDir(), ".aws", "credentials")
				assert.NoError(t, os.MkdirAll(filepath.Dir(credentialsFile), 0o700))
				assert.NoError(t, os.WriteFile(credentialsFile, []byte(`[test-profile]
aws_access_key_id = test-id
aws_secret_access_key = test-secret`), 0600))

				t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsFile)
			},
			assertFunc: func(t *testing.T, cfg *aws.Config, err error) {
				assert.NoError(t, err)

				creds, err := cfg.Credentials.Retrieve(context.TODO())
				assert.NoError(t, err)
				assert.Equal(t, "test-id", creds.AccessKeyID)
				assert.Equal(t, "test-secret", creds.SecretAccessKey)

				// ConfigSources is a slice of config.Config, which in turn is an interface.
				// Since we use a LoadOptions object, we assert the type of cfgSrc and then
				// check if the expected profile is present.
				for _, cfgSrc := range cfg.ConfigSources {
					if src, ok := cfgSrc.(config.LoadOptions); ok {
						assert.Equal(t, "test-profile", src.SharedConfigProfile)
					}
				}
			},
		},
		{
			name: "invalid arn",
			key: MasterKey{
				Arn: "arn:gcp:kms:antartica-north-2::key/45e6-aca6-a5b005693a48",
			},
			assertFunc: func(t *testing.T, cfg *aws.Config, err error) {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "no valid ARN found")
				assert.Nil(t, cfg)
			},
		},
		{
			name: "STS config attempt",
			key: MasterKey{
				Arn:  dummyARN,
				Role: "role",
			},
			assertFunc: func(t *testing.T, cfg *aws.Config, err error) {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "failed to assume role 'role'")
				assert.Nil(t, cfg)
			},
		},
		{
			name: "client default fallback",
			key: MasterKey{
				Arn: "arn:aws:kms:us-west-2:107501996527:key/612d5f0p-p1l3-45e6-aca6-a5b005693a48",
			},
			envFunc: func(t *testing.T) {
				t.Setenv("AWS_ACCESS_KEY_ID", "id")
				t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
				t.Setenv("AWS_SESSION_TOKEN", "token")
			},
			assertFunc: func(t *testing.T, cfg *aws.Config, err error) {
				assert.NoError(t, err)

				creds, err := cfg.Credentials.Retrieve(context.TODO())
				assert.Nil(t, err)
				assert.Equal(t, "id", creds.AccessKeyID)
				assert.Equal(t, "secret", creds.SecretAccessKey)
				assert.Equal(t, "token", creds.SessionToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if tt.envFunc != nil {
				tt.envFunc(t)
			}
			cfg, err := tt.key.createKMSConfig()
			tt.assertFunc(t, cfg, err)
		})
	}
}

func TestMasterKey_createSTSConfig(t *testing.T) {
	t.Run("session name error", func(t *testing.T) {
		defer func() { osHostname = os.Hostname }()
		osHostname = func() (name string, err error) {
			err = fmt.Errorf("an error")
			return
		}
		key := NewMasterKeyFromArn(dummyARN, nil, "")
		cfg, err := key.createSTSConfig(nil)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to construct STS session name")
		assert.Nil(t, cfg)
	})

	t.Run("role assumption error", func(t *testing.T) {
		key := NewMasterKeyFromArn(dummyARN, nil, "")
		key.Role = "role"
		got, err := key.createSTSConfig(&aws.Config{})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to assume role 'role'")
		assert.Nil(t, got)
	})
}

func Test_stsSessionName(t *testing.T) {
	t.Run("STS session name", func(t *testing.T) {
		defer func() { osHostname = os.Hostname }()
		const mockHostname = "hostname"
		osHostname = func() (name string, err error) {
			name = mockHostname
			return
		}
		got, err := stsSessionName()
		assert.NoError(t, err)
		assert.Equal(t, "sops@"+mockHostname, got)
	})

	t.Run("hostname error", func(t *testing.T) {
		defer func() { osHostname = os.Hostname }()
		osHostname = func() (name string, err error) {
			err = fmt.Errorf("an error")
			return
		}
		got, err := stsSessionName()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to construct STS session name")
		assert.Empty(t, got)
	})

	t.Run("replaces with stsSessionRegex", func(t *testing.T) {
		const mockHostname = "some-hostname"
		defer func() { osHostname = os.Hostname }()
		osHostname = func() (name string, err error) {
			name = mockHostname
			return
		}
		got, err := stsSessionName()
		assert.NoError(t, err)
		assert.Equal(t, "sops@somehostname", got)
	})

	t.Run("hostname exceeding roleSessionNameLengthLimit", func(t *testing.T) {
		const mockHostname = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		defer func() { osHostname = os.Hostname }()
		osHostname = func() (name string, err error) {
			name = mockHostname
			return
		}
		got, err := stsSessionName()
		assert.NoError(t, err)
		assert.NotEqual(t, "sops@"+mockHostname, got)
		assert.Len(t, got, roleSessionNameLengthLimit)
	})
}

// createTestMasterKey creates a MasterKey with the provided ARN and a dummy
// credentials.StaticCredentialsProvider.
func createTestMasterKey(arn string) MasterKey {
	return MasterKey{
		Arn:                 arn,
		credentialsProvider: credentials.NewStaticCredentialsProvider("id", "secret", ""),
		baseEndpoint:        testKMSServerURL,
	}
}

// createTestKMSClient creates a new client with the
// aws.EndpointResolverWithOptions set to epResolver.
func createTestKMSClient(key MasterKey) (*kms.Client, error) {
	cfg, err := key.createKMSConfig()
	if err != nil {
		return nil, err
	}
	return kms.NewFromConfig(*cfg, func(options *kms.Options) {
		options.BaseEndpoint = aws.String(testKMSServerURL)
	}), nil
}
