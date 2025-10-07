package ocikms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMasterKeysFromOCIDString(t *testing.T) {
	s := "ocid1.key.oc1.uk-london-1.aaaalgz5aacmg.aaaailjtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq, ocid1.key.oc1.uk-london-1.bbbblgz5aacmg.bbbbiljtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq"
	ks := MasterKeysFromOCIDString(s)
	k1 := ks[0]
	k2 := ks[1]
	expectedOcid1 := "ocid1.key.oc1.uk-london-1.aaaalgz5aacmg.aaaailjtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq"
	expectedOcid2 := "ocid1.key.oc1.uk-london-1.bbbblgz5aacmg.bbbbiljtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq"
	if k1.Ocid != expectedOcid1 {
		t.Errorf("Ocid mismatch. Expected %s, found %s", expectedOcid1, k1.Ocid)
	}
	if k2.Ocid != expectedOcid2 {
		t.Errorf("Ocid mismatch. Expected %s, found %s", expectedOcid2, k2.Ocid)
	}
}

func TestKeyToMap(t *testing.T) {
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		Ocid:         "foo",
		EncryptedKey: "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"ocid":       "foo",
		"enc":        "this is encrypted",
		"created_at": "2016-10-31T10:00:00Z",
	}, key.ToMap())
}
