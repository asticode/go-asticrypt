package astimail

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Body names
const (
	NameEmailCreate = "email.create"
	NameEmailFetch  = "email.fetch"
)

// BodyKey is a body containing a key
type BodyKey struct {
	Key *PublicKey `json:"key,omitempty"`
}

// BodyMessage is a body containing an encrypted message
type BodyMessage struct {
	Name    string            `json:"name"`
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
	Payload   interface{} `json:"payload"`
}

// NewBodyMessage builds a new body containing an encrypted message and a name
func NewBodyMessage(name string, i interface{}, prvSrc *PrivateKey, pubSrc, pubDst *PublicKey, now time.Time) (b BodyMessage, err error) {
	// Init
	b = BodyMessage{Key: pubSrc, Name: name}

	// Encrypt message
	if b.Message, err = NewEncryptedMessage(BodyMessageOut{CreatedAt: now, Payload: i}, prvSrc, pubDst); err != nil {
		err = errors.Wrap(err, "creating new encrypted message failed")
		return
	}
	return
}

// Decrypt decrypts the body containing the message
func (b *BodyMessage) Decrypt(o interface{}, prvSrc *PrivateKey, pubDst *PublicKey, now time.Time) (err error) {
	// Decrypt the message
	var m BodyMessageIn
	if err = b.Message.Decrypt(&m, prvSrc, pubDst); err != nil {
		err = errors.Wrap(err, "decrypting message failed")
		return
	}

	// Validate the message creation date
	if m.CreatedAt.After(now.Add(5*time.Second)) || m.CreatedAt.Before(now.Add(-5*time.Second)) {
		err = fmt.Errorf("Request creation date %s is invalid compared to now %s", m.CreatedAt, now)
		return
	}

	// Unmarshal payload
	if err = json.Unmarshal(m.Payload, o); err != nil {
		err = errors.Wrap(err, "unmarshaling payload failed")
		return
	}
	return
}
