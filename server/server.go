package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astimail"
	"github.com/asticode/go-astitools/template"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

var templates *template.Template

func serve(addr, pathResources string) (err error) {
	// Parse templates
	if templates, err = astitemplate.ParseDirectory(filepath.Join(pathResources, "templates"), ".html"); err != nil {
		return
	}

	// Build router
	var r = httprouter.New()

	// HTML
	r.GET("/", handleHomepage)
	r.GET("/validate_email/:token", handleValidateEmail)
	r.ServeFiles("/static/*filepath", http.Dir(filepath.Join(pathResources, "static")))

	// JSON
	r.POST("/users", handleCreateUser)

	// Encrypted
	r.POST("/encrypted", handleEncryptedMessages)

	// Listen
	astilog.Debugf("Listening on %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, adaptHandler(r)); err != nil {
			astilog.Error(err)
		}
	}()
	return
}

func adaptHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		astilog.Debugf("handling %s", r.URL.Path)
		h.ServeHTTP(rw, r)
	})
}

func handleErrorHTML(rw http.ResponseWriter, err error, msgDev, msgUser string) {
	astilog.Error(errors.Wrap(err, msgDev+" failed"))
	executeTemplate(rw, "/error.html", astimail.BodyError{Label: msgUser})
}

func executeTemplate(rw http.ResponseWriter, name string, data interface{}) {
	// Check if template exists
	var t *template.Template
	if t = templates.Lookup(name); t == nil {
		handleErrorHTML(rw, errors.New("template not found"), "looking up template", "")
		return
	}

	// Execute template
	astilog.Debugf("Executing template %s", name)
	if err := t.Execute(rw, data); err != nil {
		handleErrorHTML(rw, err, "executing template", "")
		return
	}
}

func handleHomepage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	executeTemplate(rw, "/index.html", nil)
}

func handleValidateEmail(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Init
	const defaultUserErrorMsg = "Validating email failed"

	// Fetch email
	var e *Email
	var err error
	if e, err = storage.EmailFetchWithValidationToken(p.ByName("token")); err != nil {
		handleErrorHTML(rw, err, "fetching email", defaultUserErrorMsg)
		return
	}

	// Validate email
	if err = storage.EmailValidate(e); err != nil {
		handleErrorHTML(rw, err, "validating email", defaultUserErrorMsg)
		return
	}

	// Execute template
	executeTemplate(rw, "/email_validated.html", nil)
}

func handleErrorJSON(rw http.ResponseWriter, code int, err error, msgDev, msgUser string) {
	rw.WriteHeader(code)
	astilog.Error(errors.Wrap(err, msgDev+" failed"))
	if errWrite := json.NewEncoder(rw).Encode(astimail.BodyError{Label: msgUser}); errWrite != nil {
		astilog.Errorf("%s while writing", errWrite)
	}
}

func handleCreateUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Init
	const defaultUserErrorMsg = "Creating user failed"

	// Decode body
	var b astimail.BodyKey
	var err error
	if err = json.NewDecoder(r.Body).Decode(&b); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "decoding body", defaultUserErrorMsg)
		return
	}

	// Generate server private key
	// TODO Use passphrase?
	astilog.Debugf("Generating new private key")
	var srvPrvKey *astimail.PrivateKey
	if srvPrvKey, err = astimail.GeneratePrivateKey(""); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "generating server private key", defaultUserErrorMsg)
		return
	}

	// Fetch user
	if _, err = storage.UserFetchWithKey(b.Key); err != nil && err != errNotFound {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "fetching user", defaultUserErrorMsg)
		return
	} else if err == nil {
		handleErrorJSON(rw, http.StatusInternalServerError, errors.New("user already exists"), "creating user", defaultUserErrorMsg)
		return
	}

	// Create user
	if err = storage.UserCreate(b.Key, srvPrvKey); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "creating user", defaultUserErrorMsg)
		return
	}

	// Write
	if err = json.NewEncoder(rw).Encode(astimail.BodyKey{Key: srvPrvKey.Public()}); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "writing", defaultUserErrorMsg)
		return
	}
}

func handleErrorEncrypted(rw http.ResponseWriter, u *User, err error, msgDev, msgUser string) {
	// Log
	astilog.Error(errors.Wrap(err, msgDev+" failed"))

	// Build body
	var b astimail.BodyMessage
	if b, err = astimail.NewBodyMessage(astimail.NameError, astimail.BodyError{Label: msgUser}, u.ServerPrivateKey, u.ServerPrivateKey.Public(), u.ClientPublicKey, time.Now()); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "building body", msgUser)
		return
	}

	// Write
	if err = json.NewEncoder(rw).Encode(b); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "writing", msgUser)
		return
	}
}

func handleEncryptedMessages(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Init
	var userErrorMsg = "Communicating with server failed"

	// Decode body
	var b astimail.BodyMessage
	var err error
	if err = json.NewDecoder(r.Body).Decode(&b); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "decoding body", userErrorMsg)
		return
	}

	// Fetch user based on the key
	var u *User
	if u, err = storage.UserFetchWithKey(b.Key); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "decoding body", userErrorMsg)
		return
	}

	// Decrypt message
	var m astimail.BodyMessageIn
	if m, err = b.Decrypt(u.ServerPrivateKey, u.ClientPublicKey, time.Now()); err != nil {
		handleErrorEncrypted(rw, u, err, "decrypting message", userErrorMsg)
		return
	}

	// Switch on name
	var data interface{}
	switch m.Name {
	case astimail.NameEmailAdd:
		data, userErrorMsg, err = handleEmailAdd(m.Payload, u)
	case astimail.NameEmailFetch:
		data, userErrorMsg, err = handleEmailFetch(m.Payload, u)
	case astimail.NameEmailList:
		data, userErrorMsg, err = handleEmailList(u)
	default:
		err = errors.New("Unknown b.Name")
	}

	// Process error
	if err != nil {
		handleErrorEncrypted(rw, u, err, fmt.Sprintf("handling %s", m.Name), userErrorMsg)
		return
	}

	// Build body
	if b, err = astimail.NewBodyMessage(m.Name, data, u.ServerPrivateKey, u.ServerPrivateKey.Public(), u.ClientPublicKey, time.Now()); err != nil {
		handleErrorEncrypted(rw, u, err, "building body", userErrorMsg)
		return
	}

	// Write
	if err = json.NewEncoder(rw).Encode(b); err != nil {
		handleErrorEncrypted(rw, u, err, "writing", userErrorMsg)
		return
	}
}

func handleEmailAdd(payload json.RawMessage, u *User) (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Adding email failed"

	// Unmarshal payload
	var email string
	if err = json.Unmarshal(payload, &email); err != nil {
		err = errors.Wrap(err, "unmarshaling failed")
		return
	}

	// Validate email
	if !govalidator.IsEmail(email) {
		userErrorMsg = "Email is invalid"
		err = fmt.Errorf("validating email %s failed", email)
		return
	}

	// Fetch user based on the email
	if _, err = storage.UserFetchWithEmail(email); err != nil && err != errNotFound {
		err = errors.Wrap(err, "fetching email failed")
		return
	}

	// Email already exists
	if err == nil {
		userErrorMsg = "Email is already associated to a user"
		err = errors.New("Email already exists")
		return
	}

	// Create email
	var token string
	if token, err = storage.EmailCreate(email, u); err != nil {
		err = errors.Wrap(err, "creating email failed")
		return
	}
	astilog.Debugf("Token is %s", token)

	// TODO Send validation link

	// Set data
	data = "An email has been sent to you containing instructions to follow"
	return
}

func handleEmailFetch(payload json.RawMessage, u *User) (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Fetching email failed"

	// Unmarshal payload
	var email string
	if err = json.Unmarshal(payload, &email); err != nil {
		err = errors.Wrap(err, "unmarshaling failed")
		return
	}

	// Fetch user based on the email
	if _, err = storage.UserFetchWithEmail(email); err != nil && err != errNotFound {
		err = errors.Wrap(err, "fetching email failed")
		return
	}
	return
}

func handleEmailList(u *User) (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Listing emails failed"

	// List emails
	var es []*Email
	if es, err = storage.EmailList(u); err != nil {
		err = errors.Wrap(err, "listing email failed")
		return
	}

	// Build data
	var emails = []string{}
	for _, e := range es {
		emails = append(emails, e.Addr)
	}
	data = emails
	return
}
