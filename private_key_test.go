package astimail_test

import (
	"testing"

	"github.com/asticode/go-astimail"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKey(t *testing.T) {
	var k = &astimail.PrivateKey{}
	k.SetPassphrase("test")
	err := k.UnmarshalText([]byte(prv1))
	assert.NoError(t, err)
	assert.Equal(t, prv1, k.String())
	assert.Equal(t, pub1, k.Public().String())
}
