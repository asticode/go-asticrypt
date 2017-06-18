package main

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astimail"
	"github.com/pkg/errors"
)

// handleMessageEmailAdd handles the "email.add" message
func handleMessageEmailAdd(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Unmarshal payload
	var email string
	if err = json.Unmarshal(m.Payload, &email); err != nil {
		err = errors.Wrap(err, "unmarshaling payload failed")
		return
	}

	// Add email
	var label string
	if err = sendEncryptedHTTPRequest(astimail.NameEmailAdd, email, &label); err != nil {
		err = errors.Wrap(err, "adding email failed")
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.added", Payload: label}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}

// handleMessageEmailList handles the "email.list" message
func handleMessageEmailList(w *astilectron.Window) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// List emails
	var emails []string
	if err = sendEncryptedHTTPRequest(astimail.NameEmailList, nil, &emails); err != nil {
		err = errors.Wrap(err, "listing emails failed")
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.listed", Payload: emails}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}
