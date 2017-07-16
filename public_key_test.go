package asticrypt_test

import (
	"testing"

	"github.com/asticode/go-asticrypt"
	"github.com/stretchr/testify/assert"
)

func TestPublicKey(t *testing.T) {
	var k = &asticrypt.PublicKey{}
	err := k.UnmarshalText([]byte(pub1))
	assert.NoError(t, err)
	assert.Equal(t, pub1, k.String())
	assert.Equal(t, []byte{0xc7, 0x40, 0xf5, 0x48, 0xbf, 0x53, 0x13, 0x32, 0x85, 0xf0, 0x5a, 0xec, 0xb7, 0x35, 0xd1, 0xe9, 0xe6, 0x81, 0x8, 0xe}, k.Hash())
}
