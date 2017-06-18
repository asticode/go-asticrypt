package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// handleMessages handles messages
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) {
	switch m.Name {
	case "index.login":
		handleMessageIndexLogin(w, m)
	case "index.show":
		handleMessageIndexShow(w)
	case "index.sign.up":
		handleMessageIndexSignUp(w, m)
	}
}

// processMessageError processes the message error
func processMessageError(w *astilectron.Window, err *error) {
	if *err != nil {
		astilog.Error(*err)
		if errSend := w.Send(bootstrap.MessageOut{Name: "error", Payload: errors.Cause(*err).Error()}); errSend != nil {
			astilog.Error("Sending error message failed")
		}
	}
}
