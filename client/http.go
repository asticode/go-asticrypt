package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/asticode/go-astimail"
	"github.com/pkg/errors"
)

// signup signs up
func signup(password string) (err error) {
	// Generate private key
	var cltPrvKey *astimail.PrivateKey
	if cltPrvKey, err = astimail.GeneratePrivateKey(password); err != nil {
		err = errors.Wrap(err, "generating private key failed")
		return
	}

	// Marshal body
	var b []byte
	if b, err = json.Marshal(astimail.BodyKey{Key: cltPrvKey.Public()}); err != nil {
		err = errors.Wrap(err, "marshaling body failed")
		return
	}

	// Create new request
	var r *http.Request
	if r, err = http.NewRequest(http.MethodPost, ServerPublicAddr+"/users", bytes.NewReader(b)); err != nil {
		err = errors.Wrap(err, "creating http request failed")
		return
	}

	// Send request
	var resp *http.Response
	if resp, err = httpClient.Do(r); err != nil {
		err = errors.Wrap(err, "sending request failed")
		return
	}
	defer resp.Body.Close()

	// Unmarshal body
	var body astimail.BodyKey
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		err = errors.Wrap(err, "unmarshaling body failed")
		return
	}

	// Set keys
	clientPrivateKey = &astimail.PrivateKey{}
	*clientPrivateKey = *cltPrvKey
	serverPublicKey = &astimail.PublicKey{}
	*serverPublicKey = *body.Key
	return
}
