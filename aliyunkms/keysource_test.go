package aliyunkms


import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyToMap(t *testing.T) {
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		Role:         "bar",
		EncryptedKey: "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"role":        "bar",
		"enc":         "this is encrypted",
		"created_at":  "2016-10-31T10:00:00Z",
	}, key.ToMap())
}
