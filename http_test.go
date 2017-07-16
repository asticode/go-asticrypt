package asticrypt_test

import (
	"testing"
	"time"

	"github.com/asticode/go-asticrypt"
	"github.com/stretchr/testify/assert"
)

func TestBodyMessage(t *testing.T) {
	// Init
	var pk1, pk2 = &asticrypt.PrivateKey{}, &asticrypt.PrivateKey{}
	pk1.SetPassphrase("test")
	err := pk1.UnmarshalText([]byte(prv1))
	assert.NoError(t, err)
	err = pk2.UnmarshalText([]byte(prv2))
	assert.NoError(t, err)

	// Assert
	b, err := asticrypt.NewBodyMessage("name", "test", pk1, pk1.Public(), pk2.Public(), time.Now())
	assert.NoError(t, err)
	m, err := b.Decrypt(pk2, pk1.Public(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, "\"test\"", string(m.Payload))
	assert.Equal(t, "name", m.Name)
}
