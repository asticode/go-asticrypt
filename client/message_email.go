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
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var email string
	var err error
	if err = json.Unmarshal(m.Payload, &email); err != nil {
		msgError.update(err, "unmarshaling payload failed", "")
		return
	}

	// Add email
	var label string
	if err = sendEncryptedHTTPRequest(astimail.NameEmailAdd, email, &label); err != nil {
		msgError.update(err, "adding email failed", "")
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.added", Payload: label}); err != nil {
		msgError.update(err, "sending message failed", "")
		return
	}
}

// handleMessageEmailList handles the "email.list" message
func handleMessageEmailList(w *astilectron.Window) {
	// Process errors
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// List emails
	var emails []string
	var err error
	if err = sendEncryptedHTTPRequest(astimail.NameEmailList, nil, &emails); err != nil {
		msgError.update(err, "listing emails failed", "")
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.listed", Payload: emails}); err != nil {
		msgError.update(err, "sending message failed", "")
		return
	}
}
