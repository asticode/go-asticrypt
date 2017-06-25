package astimail

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Body names
const (
	NameEmailAdd   = "email.add"
	NameEmailFetch = "email.fetch"
	NameEmailList  = "email.list"
	NameError      = "error"
	NameReferences = "references"
)

// BodyError is a body containing an error
type BodyError struct {
	Label string `json:"label"`
}

// Error implements the error interface
func (b BodyError) Error() string {
	return b.Label
}

// BodyKey is a body containing a key
type BodyKey struct {
	Key *PublicKey `json:"key,omitempty"`
}

// BodyMessage is a body containing an encrypted message
type BodyMessage struct {
	Message *EncryptedMessage `json:"message,omitempty"`
	Key     *PublicKey        `json:"key,omitempty"`
}

// BodyMessageIn represents the body of a message coming in
type BodyMessageIn struct {
	BodyMessageOut
	Payload json.RawMessage `json:"payload"`
}

// BodyMessageOut represents the body of a message going out
type BodyMessageOut struct {
	CreatedAt time.Time   `json:"created_at"`
	Name      string      `json:"name"`
	Payload   interface{} `json:"payload"`
}

// NewBodyMessage builds a new body containing an encrypted message and a name
func NewBodyMessage(name string, i interface{}, prvSrc *PrivateKey, pubSrc, pubDst *PublicKey, now time.Time) (b BodyMessage, err error) {
	// Init
	b = BodyMessage{Key: pubSrc}

	// Encrypt message
	if b.Message, err = NewEncryptedMessage(BodyMessageOut{CreatedAt: now, Name: name, Payload: i}, prvSrc, pubDst); err != nil {
		err = errors.Wrap(err, "creating new encrypted message failed")
		return
	}
	return
}

// Decrypt decrypts the body containing the message
func (b *BodyMessage) Decrypt(prvSrc *PrivateKey, pubDst *PublicKey, now time.Time) (m BodyMessageIn, err error) {
	// Decrypt the message
	if err = b.Message.Decrypt(&m, prvSrc, pubDst); err != nil {
		err = errors.Wrap(err, "decrypting message failed")
		return
	}

	// Validate the message creation date
	if m.CreatedAt.After(now.Add(5*time.Second)) || m.CreatedAt.Before(now.Add(-5*time.Second)) {
		err = fmt.Errorf("Request creation date %s is invalid compared to now %s", m.CreatedAt, now)
		return
	}
	return
}

// BodyReferences represents a body containing references
type BodyReferences struct {
	Now time.Time `json:"now"`
}
