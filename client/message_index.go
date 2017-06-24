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
)

// handleMessageIndexShow handles the "index.show" message
func handleMessageIndexShow(w *astilectron.Window) {
	// Process errors
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

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
	var err error
	if err = w.Send(m); err != nil {
		msgError.update(err, "sending message", "")
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
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var password string
	var err error
	if err = json.Unmarshal(m.Payload, &password); err != nil {
		msgError.update(err, "unmarshaling payload", "")
		return
	}

	// Generate private key
	var cltPrvKey *astimail.PrivateKey
	astilog.Debug("Generating new private key")
	if cltPrvKey, err = astimail.GeneratePrivateKey(password); err != nil {
		msgError.update(err, "generating private key", "")
		return
	}

	// Send HTTP request
	var body astimail.BodyKey
	if err = sendHTTPRequest(http.MethodPost, "/users", astimail.BodyKey{Key: cltPrvKey.Public()}, &body); err != nil {
		msgError.update(err, "sending http request", "")
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
		msgError.update(err, "creating configuration file", "")
		return
	}
	defer f.Close()

	// Write configuration
	if err = toml.NewEncoder(f).Encode(Configuration{
		ClientPrivateKey: clientPrivateKey,
		ServerPublicKey:  serverPublicKey,
	}); err != nil {
		msgError.update(err, "creating configuration file", "")
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "index.signed.up"}); err != nil {
		msgError.update(err, "sending message", "")
		return
	}
}

// handleMessageIndexLogin handles the "index.login" message
func handleMessageIndexLogin(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var password string
	var err error
	if err = json.Unmarshal(m.Payload, &password); err != nil {
		msgError.update(err, "unmarshaling payload", "")
		return
	}

	// Parse configuration
	var c = Configuration{ClientPrivateKey: &astimail.PrivateKey{}}
	c.ClientPrivateKey.SetPassphrase(password)
	if _, err = toml.DecodeFile(pathConfiguration, &c); err != nil {
		msgError.update(err, "decoding toml file", "")
		return
	}

	// Set keys
	clientPrivateKey = &astimail.PrivateKey{}
	*clientPrivateKey = *c.ClientPrivateKey
	serverPublicKey = &astimail.PublicKey{}
	*serverPublicKey = *c.ServerPublicKey

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "index.logged.in"}); err != nil {
		msgError.update(err, "sending message", "")
		return
	}
}
