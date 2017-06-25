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

// handleMessageIndex handles the "index" message
func handleMessageIndex(w *astilectron.Window) {
	// Process errors
	const defaultUserErrorMsg = "Showing index failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Build message
	var m = bootstrap.MessageOut{Name: "indexed"}
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
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// fetchReferences fetches references
func fetchReferences() (err error) {
	// Fetch references
	var body astimail.BodyReferences
	if err = sendEncryptedHTTPRequest(astimail.NameReferences, nil, &body); err != nil {
		err = errors.Wrap(err, "sending encrypted http request failed")
		return
	}

	// Update references
	now = body.Now
	return
}

// Configuration represents a configuration
type Configuration struct {
	ClientPrivateKey *astimail.PrivateKey `toml:"client_private_key"`
	ServerPublicKey  *astimail.PublicKey  `toml:"server_public_key"`
}

// handleMessageSignUp handles the "sign.up" message
func handleMessageSignUp(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	const defaultUserErrorMsg = "Signing up failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var password string
	var err error
	if err = json.Unmarshal(m.Payload, &password); err != nil {
		msgError.update(err, "unmarshaling payload", defaultUserErrorMsg)
		return
	}

	// Generate private key
	var cltPrvKey *astimail.PrivateKey
	astilog.Debug("Generating new private key")
	if cltPrvKey, err = astimail.GeneratePrivateKey(password); err != nil {
		msgError.update(err, "generating private key", defaultUserErrorMsg)
		return
	}

	// Send HTTP request
	var body astimail.BodyKey
	if err = sendHTTPRequest(http.MethodPost, "/users", astimail.BodyKey{Key: cltPrvKey.Public()}, &body); err != nil {
		msgError.update(err, "sending http request", defaultUserErrorMsg)
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
		msgError.update(err, "creating configuration file", defaultUserErrorMsg)
		return
	}
	defer f.Close()

	// Write configuration
	if err = toml.NewEncoder(f).Encode(Configuration{
		ClientPrivateKey: clientPrivateKey,
		ServerPublicKey:  serverPublicKey,
	}); err != nil {
		msgError.update(err, "creating configuration file", defaultUserErrorMsg)
		return
	}

	// Fetch references
	if err = fetchReferences(); err != nil {
		msgError.update(err, "fetching references", defaultUserErrorMsg)
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "signed.up"}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// handleMessageLogin handles the "login" message
func handleMessageLogin(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	const defaultUserErrorMsg = "Logging in failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var password string
	var err error
	if err = json.Unmarshal(m.Payload, &password); err != nil {
		msgError.update(err, "unmarshaling payload", defaultUserErrorMsg)
		return
	}

	// Parse configuration
	var c = Configuration{ClientPrivateKey: &astimail.PrivateKey{}}
	c.ClientPrivateKey.SetPassphrase(password)
	if _, err = toml.DecodeFile(pathConfiguration, &c); err != nil {
		msgError.update(err, "decoding toml file", defaultUserErrorMsg)
		return
	}

	// Set keys
	clientPrivateKey = &astimail.PrivateKey{}
	*clientPrivateKey = *c.ClientPrivateKey
	serverPublicKey = &astimail.PublicKey{}
	*serverPublicKey = *c.ServerPublicKey

	// Fetch references
	if err = fetchReferences(); err != nil {
		msgError.update(err, "fetching references", defaultUserErrorMsg)
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "logged.in"}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// handleMessageLogout handles the "logout" message
func handleMessageLogout(w *astilectron.Window) {
	// Process errors
	const defaultUserErrorMsg = "Logging out failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Set keys
	clientPrivateKey = nil
	serverPublicKey = nil

	// Send
	var err error
	if err = w.Send(bootstrap.MessageOut{Name: "logged.out"}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}
