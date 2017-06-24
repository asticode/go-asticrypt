package main

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astimail"
)

// handleMessageEmailAdd handles the "email.add" message
func handleMessageEmailAdd(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	const defaultUserErrorMsg = "Adding email failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var email string
	var err error
	if err = json.Unmarshal(m.Payload, &email); err != nil {
		msgError.update(err, "unmarshaling payload", defaultUserErrorMsg)
		return
	}

	// Add email
	var label string
	if err = sendEncryptedHTTPRequest(astimail.NameEmailAdd, email, &label); err != nil {
		msgError.update(err, "adding email", defaultUserErrorMsg)
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.added", Payload: label}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// handleMessageEmailList handles the "email.list" message
func handleMessageEmailList(w *astilectron.Window) {
	// Process errors
	const defaultUserErrorMsg = "Listing emails failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// List emails
	var emails []string
	var err error
	if err = sendEncryptedHTTPRequest(astimail.NameEmailList, nil, &emails); err != nil {
		msgError.update(err, "listing emails", defaultUserErrorMsg)
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.listed", Payload: emails}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}
