package gcpkms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGCPKMSKeySourceFromString(t *testing.T) {
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

func TestKeyToMap(t *testing.T) {
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		ResourceID:   "foo",
		EncryptedKey: "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"resource_id": "foo",
		"enc":         "this is encrypted",
		"created_at":  "2016-10-31T10:00:00Z",
	}, key.ToMap())
}
