package yckms

import (
	"context"
	"encoding/base64"
	yckms "github.com/yandex-cloud/go-genproto/yandex/cloud/kms/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	dummyKey        = "920aff2e-c5f1-4040-943a-047fa387b27e"
	anotherDummyKey = "920aff2e-c5f1-4040-943a-047fa587b27e"
	dummyKeys       = dummyKey + ", " + anotherDummyKey
	decodedKey      = "I want to be a DJ"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	yckms.RegisterSymmetricCryptoServiceServer(s, mockSymmetricCryptoServiceServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type mockSymmetricCryptoServiceServer struct {
}

func (mockSymmetricCryptoServiceServer) Encrypt(ctx context.Context, req *yckms.SymmetricEncryptRequest) (*yckms.SymmetricEncryptResponse, error) {
	return &yckms.SymmetricEncryptResponse{
		Ciphertext: []byte(base64.StdEncoding.EncodeToString(req.Plaintext)),
	}, nil
}

func (mockSymmetricCryptoServiceServer) Decrypt(ctx context.Context, req *yckms.SymmetricDecryptRequest) (*yckms.SymmetricDecryptResponse, error) {
	plain, err := base64.StdEncoding.DecodeString(string(req.Ciphertext))
	if err != nil {
		return nil, err
	}
	return &yckms.SymmetricDecryptResponse{
		Plaintext: plain,
	}, nil
}
func (mockSymmetricCryptoServiceServer) ReEncrypt(context.Context, *yckms.SymmetricReEncryptRequest) (*yckms.SymmetricReEncryptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReEncrypt not implemented")
}
func (mockSymmetricCryptoServiceServer) GenerateDataKey(context.Context, *yckms.GenerateDataKeyRequest) (*yckms.GenerateDataKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateDataKey not implemented")
}

func TestNewMasterKeyFromKeyID(t *testing.T) {
	key := NewMasterKeyFromKeyID(dummyKey)
	assert.Equal(t, dummyKey, key.KeyID)
	assert.NotNil(t, key.CreationDate)
}

func TestNewMasterKeyFromKeyIDString(t *testing.T) {
	keys := NewMasterKeyFromKeyIDString(dummyKeys)
	assert.Len(t, keys, 2)

	k1 := keys[0]
	k2 := keys[1]

	assert.Equal(t, dummyKey, k1.KeyID)
	assert.Equal(t, anotherDummyKey, k2.KeyID)
}

func TestMasterKey_Encrypt(t *testing.T) {
	t.Run("encrypt", func(t *testing.T) {
		grpcConn, err := createMockGRPCClient()
		assert.NoError(t, err)

		key := &MasterKey{
			grpcConn: grpcConn,
		}

		dataKey := []byte(decodedKey)
		err = key.Encrypt(dataKey)
		assert.NoError(t, err)

		// Double base64 is used because encrypted data stored as base64
		// and our mock uses base64 instead of actual encryption
		assert.EqualValues(t, base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(decodedKey)))), key.EncryptedDataKey())
	})
}

func TestMasterKey_Decrypt(t *testing.T) {
	t.Run("decrypt", func(t *testing.T) {
		grpcConn, err := createMockGRPCClient()
		assert.NoError(t, err)

		// Double base64 is used because encrypted data stored as base64
		// and our mock uses base64 instead of actual encryption
		key := &MasterKey{
			EncryptedKey: base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(decodedKey)))),
			grpcConn:     grpcConn,
		}

		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.Equal(t, []byte(decodedKey), got)
	})
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

func TestMasterKey_NeedsRotation(t *testing.T) {
	t.Run("false", func(t *testing.T) {
		k := &MasterKey{}
		k.CreationDate = time.Now().UTC()

		assert.False(t, k.NeedsRotation())
	})

	t.Run("true", func(t *testing.T) {
		k := &MasterKey{}
		k.CreationDate = time.Now().UTC().Add(-kmsTTL - 1)

		assert.True(t, k.NeedsRotation())
	})
}

func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKeyFromKeyID(dummyKey)

	assert.Equal(t, dummyKey, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := NewMasterKeyFromKeyID(dummyKey)

	data := []byte("some data")
	key.SetEncryptedDataKey(data)

	res := key.ToMap()
	assert.Equal(t, dummyKey, res["key_id"])
	assert.Equal(t, key.CreationDate.UTC().Format(time.RFC3339), res["created_at"])
	assert.Equal(t, "some data", res["enc"])
}

func createMockGRPCClient() (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
