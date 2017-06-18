package astimail

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
)

// PrivateKey represents a marshalable/unmarshalable private key
type PrivateKey struct {
	key        *rsa.PrivateKey
	passphrase string
	public     *PublicKey
	string     string
}

// GeneratePrivateKey generates a new private key
func GeneratePrivateKey(passphrase string) (p *PrivateKey, err error) {
	// Generate key
	var k *rsa.PrivateKey
	if k, err = rsa.GenerateKey(rand.Reader, privateKeyBits); err != nil {
		err = errors.Wrap(err, "generating private key failed")
		return
	}

	// Build private key
	if p, err = newPrivateKey(k, passphrase); err != nil {
		err = errors.Wrap(err, "building new private key failed")
		return
	}
	return
}

// newPrivateKey builds a new private key with an optional passphrase
func newPrivateKey(i *rsa.PrivateKey, passphrase string) (k *PrivateKey, err error) {
	// Init private key
	k = &PrivateKey{
		key:        i,
		passphrase: passphrase,
	}

	// Set public field
	if k.public, err = newPublicKey(i.Public()); err != nil {
		err = errors.Wrap(err, "creating public key failed")
		return
	}

	// Set string field
	var b []byte
	if b, err = k.MarshalText(); err != nil {
		err = errors.Wrap(err, "marshaling private key failed")
		return
	}
	k.string = string(b)
	return
}

// SetPassphrase sets the passphrase
func (p *PrivateKey) SetPassphrase(passphrase string) {
	p.passphrase = passphrase
}

// Key returns the *rsa.PrivateKey
func (p PrivateKey) Key() *rsa.PrivateKey {
	return p.key
}

// Public returns the *rsa.PublicKey
func (p PrivateKey) Public() *PublicKey {
	return p.public
}

// String allows PrivateKey to implement the Stringer interface
func (p PrivateKey) String() string {
	return p.string
}

// MarshalText allows PrivateKey to implement the TextMarshaler interface
func (p PrivateKey) MarshalText() (o []byte, err error) {
	// Convert it to pem
	var block = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(p.key),
	}

	// Encrypt the pem
	if len(p.passphrase) > 0 {
		if block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(p.passphrase), x509.PEMCipherAES256); err != nil {
			err = errors.Wrap(err, "x509.EncryptPEMBlock failed")
			return
		}
	}

	// Encode to memory
	var b = pem.EncodeToMemory(block)

	// b64 encode
	o = make([]byte, b64.EncodedLen(len(b)))
	b64.Encode(o, b)
	return
}

// Scan implements the Scanner interface
func (p *PrivateKey) Scan(value interface{}) (err error) {
	var b []byte
	var ok bool
	if b, ok = value.([]byte); !ok {
		err = errors.New("value must be []byte")
		return
	}
	return p.UnmarshalText(b)
}

// UnmarshalText allows PrivateKey to implement the TextUnmarshaler interface
func (p *PrivateKey) UnmarshalText(i []byte) (err error) {
	// Base 64 decode
	var b = make([]byte, b64.DecodedLen(len(i)))
	var n int
	if n, err = b64.Decode(b, i); err != nil {
		err = errors.Wrap(err, "base64 decoding failed")
		return
	}
	b = b[:n]

	// Decode pem
	var block *pem.Block
	if block, _ = pem.Decode(b); block == nil {
		err = fmt.Errorf("No block found in pem %s", string(b))
		return
	}

	// Decrypt block
	b = block.Bytes
	if len(p.passphrase) > 0 {
		if b, err = x509.DecryptPEMBlock(block, []byte(p.passphrase)); err != nil {
			err = errors.Wrap(err, "x509.DecryptPEMBlock failed")
			return
		}
	}

	// Parse private key
	var rk *rsa.PrivateKey
	if rk, err = x509.ParsePKCS1PrivateKey(b); err != nil {
		err = errors.Wrap(err, "x509.ParsePKCS1PrivateKey failed")
		return
	}

	// Build private key
	var k *PrivateKey
	if k, err = newPrivateKey(rk, p.passphrase); err != nil {
		err = errors.Wrap(err, "creating new private key failed")
		return
	}
	// We need to assign the string field since marshaling generates a different result each time it's called
	k.string = string(i)
	*p = *k
	return
}
