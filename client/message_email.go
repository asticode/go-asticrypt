package main

import (
	"encoding/json"

	"crypto/tls"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astimail"
	"github.com/go-gomail/gomail"
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

	// Get google auth URL
	var googleAuthURL string

	// Send
	type BodyOut struct {
		Emails        []string `json:"emails"`
		GoogleAuthURL string   `json:"google_auth_url"`
	}
	if err = w.Send(bootstrap.MessageOut{Name: "email.listed", Payload: BodyOut{Emails: emails, GoogleAuthURL: googleAuthURL}}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// handleMessageEmailOpen handles the "email.open" message
func handleMessageEmailOpen(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	const defaultUserErrorMsg = "Opening email failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	type Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var b Body
	var err error
	if err = json.Unmarshal(m.Payload, &b); err != nil {
		msgError.update(err, "unmarshaling payload", defaultUserErrorMsg)
		return
	}

	// Dial smtp
	var d = gomail.NewDialer("smtp.gmail.com", 465, b.Email, b.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if _, errDial := d.Dial(); errDial != nil {
		msgError.update(err, "checking password", "")
		return
	}
	emails[b.Email] = b.Password

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "email.opened"}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}
