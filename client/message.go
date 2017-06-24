package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astimail"
	"github.com/pkg/errors"
)

// handleMessages handles messages
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) {
	switch m.Name {
	case "email.add":
		handleMessageEmailAdd(w, m)
	case "email.list":
		handleMessageEmailList(w)
	case "index.login":
		handleMessageIndexLogin(w, m)
	case "index.show":
		handleMessageIndexShow(w)
	case "index.sign.up":
		handleMessageIndexSignUp(w, m)
	}
}

// messageError represents a message error
type messageError struct {
	err     error
	userMsg string
}

// update updates the message error
func (e *messageError) update(err error, devMsg string, userMsg string) {
	e.err = errors.Wrap(err, devMsg+" failed")
	if bodyError, ok := err.(astimail.BodyError); ok {
		e.userMsg = bodyError.Label
	} else {
		e.userMsg = userMsg
	}
}

// processMessageError processes the message error
func processMessageError(w *astilectron.Window, msgError *messageError) {
	if msgError.err != nil {
		astilog.Error(msgError.err)
		if errSend := w.Send(bootstrap.MessageOut{Name: "error", Payload: msgError.userMsg}); errSend != nil {
			astilog.Error("Sending error message failed")
		}
	}
}
