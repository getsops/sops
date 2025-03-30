package gcpkms

import (
	"encoding/base64"
	"fmt"
	"net"
	"testing"
	"time"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	testResourceID = "projects/test-sops/locations/global/keyRings/test-sops/cryptoKeys/sops"
	decryptedData  = "decrypted data"
	encryptedData  = "encrypted data"
)

var (
	mockKeyManagement mockKeyManagementServer
)

func TestMasterKeysFromResourceIDString(t *testing.T) {
	s := "projects/sops-testing1/locations/global/keyRings/creds/cryptoKeys/key1, projects/sops-testing2/locations/global/keyRings/creds/cryptoKeys/key2"
	ks := MasterKeysFromResourceIDString(s)
	k1 := ks[0]
	k2 := ks[1]
	expectedResourceID1 := "projects/sops-testing1/locations/global/keyRings/creds/cryptoKeys/key1"
	expectedResourceID2 := "projects/sops-testing2/locations/global/keyRings/creds/cryptoKeys/key2"
	if k1.ResourceID != expectedResourceID1 {
		t.Errorf("ResourceID mismatch. Expected %s, found %s", expectedResourceID1, k1.ResourceID)
	}
	if k2.ResourceID != expectedResourceID2 {
		t.Errorf("ResourceID mismatch. Expected %s, found %s", expectedResourceID2, k2.ResourceID)
	}
}

func TestTokenSource_ApplyToMasterKey(t *testing.T) {
	src := NewTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "some-token"}))
	key := &MasterKey{}
	src.ApplyToMasterKey(key)
	assert.Equal(t, src.source, key.tokenSource)
}

func TestCredentialJSON_ApplyToMasterKey(t *testing.T) {
	key := &MasterKey{}
	credential := CredentialJSON("mock")
	credential.ApplyToMasterKey(key)
	assert.EqualValues(t, credential, key.credentialJSON)
}

func TestMasterKey_Encrypt(t *testing.T) {
	mockKeyManagement.err = nil
	mockKeyManagement.reqs = nil
	mockKeyManagement.resps = append(mockKeyManagement.resps[:0], &kmspb.EncryptResponse{
		Ciphertext: []byte(encryptedData),
	})

	key := MasterKey{
		grpcConn:       newGRPCServer("0"),
		ResourceID:     testResourceID,
		credentialJSON: []byte("arbitrary credentials"),
	}
	err := key.Encrypt([]byte("encrypt"))
	assert.NoError(t, err)
	assert.EqualValues(t, base64.StdEncoding.EncodeToString([]byte(encryptedData)), key.EncryptedDataKey())
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key := MasterKey{EncryptedKey: encryptedData}
	assert.EqualValues(t, encryptedData, key.EncryptedDataKey())
	assert.NoError(t, key.EncryptIfNeeded([]byte("sops data key")))
	assert.EqualValues(t, encryptedData, key.EncryptedDataKey())
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := MasterKey{EncryptedKey: encryptedData}
	assert.EqualValues(t, encryptedData, key.EncryptedDataKey())
}

func TestMasterKey_Decrypt(t *testing.T) {
	mockKeyManagement.err = nil
	mockKeyManagement.reqs = nil
	mockKeyManagement.resps = append(mockKeyManagement.resps[:0], &kmspb.DecryptResponse{
		Plaintext: []byte(decryptedData),
	})
	key := MasterKey{
		grpcConn:       newGRPCServer("0"),
		ResourceID:     testResourceID,
		EncryptedKey:   "encryptedKey",
		credentialJSON: []byte("arbitrary credentials"),
	}
	data, err := key.Decrypt()
	assert.NoError(t, err)
	assert.EqualValues(t, decryptedData, data)
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	enc := "encrypted key"
	key := &MasterKey{}
	key.SetEncryptedDataKey([]byte(enc))
	assert.EqualValues(t, enc, key.EncryptedDataKey())
}

func TestMasterKey_ToString(t *testing.T) {
	rsrcId := testResourceID
	key := NewMasterKeyFromResourceID(rsrcId)
	assert.Equal(t, rsrcId, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := MasterKey{
		credentialJSON: []byte("sensitive creds"),
		CreationDate:   time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		ResourceID:     testResourceID,
		EncryptedKey:   "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"resource_id": testResourceID,
		"enc":         "this is encrypted",
		"created_at":  "2016-10-31T10:00:00Z",
	}, key.ToMap())
}

func TestMasterKey_createCloudKMSService_withCredentialsFile(t *testing.T) {
	tests := []struct {
		key       MasterKey
		errString string
	}{
		{
			key: MasterKey{
				ResourceID:     "/projects",
				credentialJSON: []byte("some secret"),
			},
			errString: "no valid resource ID",
		},
		{
			key: MasterKey{
				ResourceID: testResourceID,
				credentialJSON: []byte(`{ "client_id": "<client-id>.apps.googleusercontent.com",
 		"client_secret": "<secret>",
		"type": "authorized_user"}`),
			},
		},
		{
			key: MasterKey{
				ResourceID: testResourceID,
			},
			errString: `credentials: failed to obtain credentials from "SOPS_GOOGLE_CREDENTIALS"`,
		},
	}

	for _, tt := range tests {
		_, err := tt.key.newKMSClient()
		if tt.errString != "" {
			assert.Error(t, err)
			assert.ErrorContains(t, err, tt.errString)
			return
		}
		assert.NoError(t, err)
	}
}

func TestMasterKey_createCloudKMSService_withOauthToken(t *testing.T) {
	t.Setenv(SopsGoogleCredentialsOAuthTokenEnv, "token")

	masterKey := MasterKey{
		ResourceID: testResourceID,
	}

	_, err := masterKey.newKMSClient()

	assert.NoError(t, err)
}

func TestMasterKey_createCloudKMSService_withoutCredentials(t *testing.T) {
	masterKey := MasterKey{
		ResourceID: testResourceID,
	}

	_, err := masterKey.newKMSClient()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "credentials: could not find default credentials")
}

func newGRPCServer(port string) *grpc.ClientConn {
	serv := grpc.NewServer()
	kmspb.RegisterKeyManagementServiceServer(serv, &mockKeyManagement)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatal(err)
	}
	go serv.Serve(lis)

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
