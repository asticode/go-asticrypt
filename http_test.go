package astimail_test

import (
	"testing"
	"time"

	"github.com/asticode/go-astimail"
	"github.com/stretchr/testify/assert"
)

func TestBodyMessage(t *testing.T) {
	// Init
	var pk1, pk2 = &astimail.PrivateKey{}, &astimail.PrivateKey{}
	pk1.SetPassphrase("test")
	err := pk1.UnmarshalText([]byte(prv1))
	assert.NoError(t, err)
	err = pk2.UnmarshalText([]byte(prv2))
	assert.NoError(t, err)

	// Assert
	b, err := astimail.NewBodyMessage("name", "test", pk1, pk1.Public(), pk2.Public(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, "name", b.Name)
	var text string
	err = b.Decrypt(&text, pk2, pk1.Public(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, "test", text)
}
