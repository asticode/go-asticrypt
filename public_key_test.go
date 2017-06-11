package astimail_test

import (
	"testing"

	"github.com/asticode/go-astimail"
	"github.com/stretchr/testify/assert"
)

func TestPublicKey(t *testing.T) {
	var k = &astimail.PublicKey{}
	err := k.UnmarshalText([]byte(pub1))
	assert.NoError(t, err)
	assert.Equal(t, pub1, k.String())
	assert.Equal(t, "\xc7@\xf5H\xbfS\x132\x85\xf0Z\xec\xb75\xd1\xe9\xe6\x81\b\x0e", k.Hash())
}
