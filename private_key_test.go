package asticrypt_test

import (
	"testing"

	"github.com/asticode/go-asticrypt"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKey(t *testing.T) {
	var k = &asticrypt.PrivateKey{}
	k.SetPassphrase("test")
	err := k.UnmarshalText([]byte(prv1))
	assert.NoError(t, err)
	assert.Equal(t, prv1, k.String())
	assert.Equal(t, pub1, k.Public().String())
}
