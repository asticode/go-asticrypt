package asticrypt

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"

	"github.com/pkg/errors"
)

// EncryptedMessage represents an encrypted message
type EncryptedMessage struct {
	Hash      []byte `json:"hash,omitempty"`
	IV        []byte `json:"iv,omitempty"`
	Key       []byte `json:"key,omitempty"`
	Message   []byte `json:"message,omitempty"`
	Signature []byte `json:"signature,omitempty"`
}

// NewEncryptedMessage encrypts a message
func NewEncryptedMessage(i interface{}, prvSrc *PrivateKey, pubDst *PublicKey) (em *EncryptedMessage, err error) {
	// Generate random key
	var key = make([]byte, aesKeyBits/8)
	if _, err = rand.Read(key); err != nil {
		err = errors.Wrap(err, "generating random key failed")
		return
	}

	// Create AES block
	var b cipher.Block
	if b, err = aes.NewCipher(key); err != nil {
		err = errors.Wrap(err, "create AES block failed")
		return
	}

	// Generate random IV
	em = &EncryptedMessage{IV: make([]byte, b.BlockSize())}
	if _, err = rand.Read(em.IV); err != nil {
		err = errors.Wrap(err, "generate random IV failed")
		return
	}

	// Marshal message
	var msg []byte
	if msg, err = json.Marshal(i); err != nil {
		err = errors.Wrap(err, "marshaling message failed")
		return
	}

	// AES encrypt the message
	em.Message = make([]byte, len(msg))
	var cfb = cipher.NewCFBEncrypter(b, em.IV)
	cfb.XORKeyStream(em.Message, msg)

	// RSA encrypt the AES key
	if em.Key, err = rsa.EncryptOAEP(sha512.New(), rand.Reader, pubDst.Key(), key, nil); err != nil {
		err = errors.Wrap(err, "rsa.EncryptOAEP failed")
		return
	}

	// Hash message
	var h = sha512.New()
	h.Write(em.Message)
	em.Hash = h.Sum(nil)

	// Sign the message
	if em.Signature, err = rsa.SignPKCS1v15(rand.Reader, prvSrc.Key(), crypto.SHA512, em.Hash); err != nil {
		err = errors.Wrap(err, "rsa.SignPKCS1v15 failed")
		return
	}
	return
}

// Decrypt decrypts a message
func (m EncryptedMessage) Decrypt(o interface{}, prvSrc *PrivateKey, pubDst *PublicKey) (err error) {
	// Verify signature
	if err = rsa.VerifyPKCS1v15(pubDst.Key(), crypto.SHA512, m.Hash, m.Signature); err != nil {
		err = errors.Wrap(err, "rsa.VerifyPKCS1v15 failed")
		return
	}

	// RSA decrypt the AES key
	var key []byte
	if key, err = rsa.DecryptOAEP(sha512.New(), rand.Reader, prvSrc.Key(), m.Key, nil); err != nil {
		err = errors.Wrap(err, "rsa.DecryptOAEP failed")
		return
	}

	// Create AES block
	var c cipher.Block
	if c, err = aes.NewCipher(key); err != nil {
		err = errors.Wrap(err, "creating AES block failed")
		return
	}

	// AES decrypt the message
	var msg = make([]byte, len(m.Message))
	var cfb = cipher.NewCFBDecrypter(c, m.IV)
	cfb.XORKeyStream(msg, m.Message)

	// Unmarshal message
	if err = json.Unmarshal(msg, o); err != nil {
		err = errors.Wrap(err, "unmarshaling message failed")
		return
	}
	return
}
