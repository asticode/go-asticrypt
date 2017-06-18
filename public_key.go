package astimail

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
)

// PublicKey represents a marshalable/unmarshalable public key
type PublicKey struct {
	hash   []byte
	key    *rsa.PublicKey
	string string
}

// newPublicKey creates a new PublicKey based on a *rsa.PublicKey
func newPublicKey(i interface{}) (k *PublicKey, err error) {
	// Init
	k = &PublicKey{}

	// Assert key
	var ok bool
	if k.key, ok = i.(*rsa.PublicKey); !ok {
		err = fmt.Errorf("Public key %s is not a *rsa.PublicKey", i)
		return
	}

	// Set string field
	var b []byte
	if b, err = k.MarshalText(); err != nil {
		err = errors.Wrap(err, "marshaling public key failed")
		return
	}
	k.string = string(b)

	// Set hash field
	var h = sha1.New()
	h.Write(b)
	k.hash = h.Sum(nil)
	return
}

// Hash hashes the public key
func (p PublicKey) Hash() []byte {
	return p.hash
}

// Key returns the *rsa.PublicKey
func (p PublicKey) Key() *rsa.PublicKey {
	return p.key
}

// String allows PublicKey to implement the Stringer interface
func (p PublicKey) String() string {
	return p.string
}

// MarshalText allows PublicKey to implement the TextMarshaler interface
func (p PublicKey) MarshalText() (o []byte, err error) {
	var b []byte
	if b, err = x509.MarshalPKIXPublicKey(p.key); err != nil {
		err = errors.Wrap(err, "x509.MarshalPKIXPublicKey failed")
		return
	}
	o = make([]byte, b64.EncodedLen(len(b)))
	b64.Encode(o, b)
	return
}

// UnmarshalText allows PublicKey to implement the TextUnmarshaler interface
func (p *PublicKey) UnmarshalText(i []byte) (err error) {
	// Base 64 decode
	var b = make([]byte, b64.DecodedLen(len(i)))
	var n int
	if n, err = b64.Decode(b, i); err != nil {
		err = errors.Wrap(err, "base64 decoding failed")
		return
	}
	b = b[:n]

	// Parse
	var in interface{}
	if in, err = x509.ParsePKIXPublicKey(b); err != nil {
		err = errors.Wrap(err, "x509.ParsePKIXPublicKey failed")
		return
	}

	// Build public key
	var k *PublicKey
	if k, err = newPublicKey(in); err != nil {
		err = errors.Wrap(err, "creating new public key failed")
		return
	}
	*p = *k
	return
}
