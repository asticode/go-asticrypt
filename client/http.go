package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astimail"
	"github.com/pkg/errors"
)

// sendHTTPRequest sends an HTTP request
func sendHTTPRequest(method string, pattern string, in interface{}, out interface{}) (err error) {
	// Marshal body
	var b []byte
	if b, err = json.Marshal(in); err != nil {
		err = errors.Wrap(err, "marshaling body failed")
		return
	}

	// Create new request
	var r *http.Request
	if r, err = http.NewRequest(method, ServerPublicAddr+pattern, bytes.NewReader(b)); err != nil {
		err = errors.Wrap(err, "creating http request failed")
		return
	}

	// Send request
	var resp *http.Response
	astilog.Debugf("Sending %s request to %s", r.Method, r.URL.Path)
	if resp, err = httpClient.Do(r); err != nil {
		err = errors.Wrap(err, "sending request failed")
		return
	}
	astilog.Debug("Request done")
	defer resp.Body.Close()

	// Unmarshal body
	if out != nil {
		if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
			err = errors.Wrap(err, "unmarshaling body failed")
			return
		}
	}
	return
}

// sendEncryptedHTTPRequest sends an encrypted HTTP request
func sendEncryptedHTTPRequest(name string, in interface{}, out interface{}) (err error) {
	// Build body
	var bout astimail.BodyMessage
	if bout, err = astimail.NewBodyMessage(name, in, clientPrivateKey, clientPrivateKey.Public(), serverPublicKey, time.Now()); err != nil {
		err = errors.Wrap(err, "building body failed")
		return
	}

	// Send HTTP request
	var bin astimail.BodyMessage
	if err = sendHTTPRequest(http.MethodPost, "/encrypted", bout, &bin); err != nil {
		err = errors.Wrap(err, "sending HTTP request failed")
		return
	}

	// Decrypt body
	var m astimail.BodyMessageIn
	if m, err = bin.Decrypt(clientPrivateKey, serverPublicKey, time.Now()); err != nil {
		err = errors.Wrap(err, "decrypting message failed")
		return
	}

	// Validate name
	if m.Name != name {
		err = fmt.Errorf("input name %s != message name %s", name, m.Name)
		return
	}

	// Unmarshal payload
	if out != nil {
		if err = json.Unmarshal(m.Payload, out); err != nil {
			err = errors.Wrap(err, "unmarshaling payload failed")
			return
		}
	}
	return
}
