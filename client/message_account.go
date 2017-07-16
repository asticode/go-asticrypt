package main

import (
	"encoding/json"

	"crypto/tls"

	"fmt"

	"github.com/asticode/go-asticrypt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/go-gomail/gomail"
)

// handleMessageAccountAdd handles the "account.add" message
func handleMessageAccountAdd(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	const defaultUserErrorMsg = "Adding account failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	var account string
	var err error
	if err = json.Unmarshal(m.Payload, &account); err != nil {
		msgError.update(err, "unmarshaling payload", defaultUserErrorMsg)
		return
	}

	// Add account
	var label string
	if err = sendEncryptedHTTPRequest(asticrypt.NameAccountAdd, account, &label); err != nil {
		msgError.update(err, "adding account", defaultUserErrorMsg)
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "account.added", Payload: label}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// handleMessageAccountList handles the "account.list" message
func handleMessageAccountList(w *astilectron.Window) {
	// Process errors
	const defaultUserErrorMsg = "Listing accounts failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// List accounts
	var es []string
	var err error
	if err = sendEncryptedHTTPRequest(asticrypt.NameAccountList, nil, &es); err != nil {
		msgError.update(err, "listing accounts", defaultUserErrorMsg)
		return
	}

	// Build accounts
	type Account struct {
		Addr    string `json:"addr"`
		AuthURL string `json:"auth_url"`
	}
	var accounts []Account
	for _, e := range es {
		accounts = append(accounts, Account{Addr: e, AuthURL: fmt.Sprintf("%s/oauth/google?account=%s", ServerPublicAddr, e)})
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "account.listed", Payload: accounts}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}

// handleMessageAccountOpen handles the "account.open" message
func handleMessageAccountOpen(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	const defaultUserErrorMsg = "Opening account failed"
	var msgError = &messageError{}
	defer processMessageError(w, msgError)

	// Unmarshal payload
	type Body struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	var b Body
	var err error
	if err = json.Unmarshal(m.Payload, &b); err != nil {
		msgError.update(err, "unmarshaling payload", defaultUserErrorMsg)
		return
	}

	// Dial smtp
	var d = gomail.NewDialer("smtp.gmail.com", 465, b.Account, b.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if _, errDial := d.Dial(); errDial != nil {
		msgError.update(err, "checking password", "")
		return
	}
	accounts[b.Account] = b.Password

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "account.opened"}); err != nil {
		msgError.update(err, "sending message", defaultUserErrorMsg)
		return
	}
}
