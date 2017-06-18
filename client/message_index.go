package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astimail"
	"github.com/pkg/errors"
)

// handleMessageIndexShow handles the "index.show" message
func handleMessageIndexShow(w *astilectron.Window) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Build message
	var m = bootstrap.MessageOut{Name: "index.show"}
	if _, errStat := os.Stat(pathConfiguration); errStat != nil {
		m.Payload = "signup"
	} else if clientPrivateKey == nil || serverPublicKey == nil {
		m.Payload = "login"
	} else {
		m.Payload = "index"
	}

	// Send
	if err = w.Send(m); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}

// Configuration represents a configuration
type Configuration struct {
	ClientPrivateKey *astimail.PrivateKey `toml:"client_private_key"`
	ServerPublicKey  *astimail.PublicKey  `toml:"server_public_key"`
}

// handleMessageIndexSignUp handles the "index.sign.up" message
func handleMessageIndexSignUp(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Unmarshal payload
	var password string
	if err = json.Unmarshal(m.Payload, &password); err != nil {
		err = errors.Wrap(err, "unmarshaling payload failed")
		return
	}

	// Generate private key
	var cltPrvKey *astimail.PrivateKey
	astilog.Debug("Generating new private key")
	if cltPrvKey, err = astimail.GeneratePrivateKey(password); err != nil {
		err = errors.Wrap(err, "generating private key failed")
		return
	}

	// Send HTTP request
	var body astimail.BodyKey
	if err = sendHTTPRequest(http.MethodPost, "/users", astimail.BodyKey{Key: cltPrvKey.Public()}, &body); err != nil {
		err = errors.Wrap(err, "sending http request failed")
		return
	}

	// Set keys
	clientPrivateKey = &astimail.PrivateKey{}
	*clientPrivateKey = *cltPrvKey
	serverPublicKey = &astimail.PublicKey{}
	*serverPublicKey = *body.Key

	// Create configuration file
	var f *os.File
	if f, err = os.Create(pathConfiguration); err != nil {
		err = errors.Wrap(err, "creating configuration file failed")
		return
	}
	defer f.Close()

	// Write configuration
	if err = toml.NewEncoder(f).Encode(Configuration{
		ClientPrivateKey: clientPrivateKey,
		ServerPublicKey:  serverPublicKey,
	}); err != nil {
		err = errors.Wrap(err, "writing configuration failed")
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "index.signed.up"}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}

// handleMessageIndexLogin handles the "index.login" message
func handleMessageIndexLogin(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Unmarshal payload
	var password string
	if err = json.Unmarshal(m.Payload, &password); err != nil {
		err = errors.Wrap(err, "unmarshaling payload failed")
		return
	}

	// Parse configuration
	var c = Configuration{ClientPrivateKey: &astimail.PrivateKey{}}
	c.ClientPrivateKey.SetPassphrase(password)
	if _, err = toml.DecodeFile(pathConfiguration, &c); err != nil {
		err = errors.Wrap(err, "decoding toml file failed")
		return
	}

	// Set keys
	clientPrivateKey = &astimail.PrivateKey{}
	*clientPrivateKey = *c.ClientPrivateKey
	serverPublicKey = &astimail.PublicKey{}
	*serverPublicKey = *c.ServerPublicKey

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "index.logged.in"}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}
