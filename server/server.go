package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"fmt"

	"github.com/asticode/go-asticrypt"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/template"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	r.GET("/oauth/:provider/redirect", handleOAuthRedirect)
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
	executeTemplate(rw, "/error.html", asticrypt.BodyError{Label: msgUser})
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

func handleOAuthRedirect(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Init
	const defaultUserErrorMsg = "OAuth failed"

	// Build auth URL
	var provider, authURL = p.ByName("provider"), ""
	switch provider {
	case "google":
		// Build config
		var config = &oauth2.Config{
			ClientID:     configuration.GoogleClientID,
			ClientSecret: configuration.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  "file://bite",
			Scopes:       []string{"https://mail.google.com"},
		}
		// TODO Use a legit state
		authURL = config.AuthCodeURL("state")
	default:
		handleErrorHTML(rw, fmt.Errorf("Invalid provider %s", provider), "building auth url", defaultUserErrorMsg)
		return
	}

	// Redirect
	http.Redirect(rw, r, authURL, http.StatusFound)
}

func handleErrorJSON(rw http.ResponseWriter, code int, err error, msgDev, msgUser string) {
	rw.WriteHeader(code)
	astilog.Error(errors.Wrap(err, msgDev+" failed"))
	if errWrite := json.NewEncoder(rw).Encode(asticrypt.BodyError{Label: msgUser}); errWrite != nil {
		astilog.Errorf("%s while writing", errWrite)
	}
}

func handleCreateUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Init
	const defaultUserErrorMsg = "Creating user failed"

	// Decode body
	var b asticrypt.BodyKey
	var err error
	if err = json.NewDecoder(r.Body).Decode(&b); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "decoding body", defaultUserErrorMsg)
		return
	}

	// Generate server private key
	// TODO Use passphrase?
	astilog.Debugf("Generating new private key")
	var srvPrvKey *asticrypt.PrivateKey
	if srvPrvKey, err = asticrypt.GeneratePrivateKey(""); err != nil {
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
	if err = json.NewEncoder(rw).Encode(asticrypt.BodyKey{Key: srvPrvKey.Public()}); err != nil {
		handleErrorJSON(rw, http.StatusInternalServerError, err, "writing", defaultUserErrorMsg)
		return
	}
}

func handleErrorEncrypted(rw http.ResponseWriter, u *User, err error, msgDev, msgUser string) {
	// Log
	astilog.Error(errors.Wrap(err, msgDev+" failed"))

	// Build body
	var b asticrypt.BodyMessage
	if b, err = asticrypt.NewBodyMessage(asticrypt.NameError, asticrypt.BodyError{Label: msgUser}, u.ServerPrivateKey, u.ServerPrivateKey.Public(), u.ClientPublicKey, time.Now()); err != nil {
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
	var b asticrypt.BodyMessage
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
	var m asticrypt.BodyMessageIn
	if m, err = b.Decrypt(u.ServerPrivateKey, u.ClientPublicKey, time.Now()); err != nil {
		handleErrorEncrypted(rw, u, err, "decrypting message", userErrorMsg)
		return
	}

	// Switch on name
	astilog.Debugf("m.Name is %s", m.Name)
	var data interface{}
	switch m.Name {
	case asticrypt.NameAccountAdd:
		data, userErrorMsg, err = handleAccountAdd(m.Payload, u)
	case asticrypt.NameAccountFetch:
		data, userErrorMsg, err = handleAccountFetch(m.Payload, u)
	case asticrypt.NameAccountList:
		data, userErrorMsg, err = handleAccountList(u)
	case asticrypt.NameReferences:
		data, userErrorMsg, err = handleReferences()
	default:
		err = errors.New("Unknown b.Name")
	}

	// Process error
	if err != nil {
		handleErrorEncrypted(rw, u, err, fmt.Sprintf("handling %s", m.Name), userErrorMsg)
		return
	}

	// Build body
	if b, err = asticrypt.NewBodyMessage(m.Name, data, u.ServerPrivateKey, u.ServerPrivateKey.Public(), u.ClientPublicKey, time.Now()); err != nil {
		handleErrorEncrypted(rw, u, err, "building body", userErrorMsg)
		return
	}

	// Write
	if err = json.NewEncoder(rw).Encode(b); err != nil {
		handleErrorEncrypted(rw, u, err, "writing", userErrorMsg)
		return
	}
}

func handleAccountAdd(payload json.RawMessage, u *User) (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Adding account failed"

	// Unmarshal payload
	var account string
	if err = json.Unmarshal(payload, &account); err != nil {
		err = errors.Wrap(err, "unmarshaling failed")
		return
	}

	// Fetch user based on the account
	if _, err = storage.UserFetchWithAccount(account); err != nil && err != errNotFound {
		err = errors.Wrap(err, "fetching account failed")
		return
	}

	// Account already exists
	if err == nil {
		userErrorMsg = "Account is already associated to a user"
		err = errors.New("Account already exists")
		return
	}

	// Create account
	var token string
	if token, err = storage.AccountCreate(account, u); err != nil {
		err = errors.Wrap(err, "creating account failed")
		return
	}
	astilog.Debugf("Token is %s", token)

	// TODO Send validation link

	// Set data
	data = "An account has been sent to you containing instructions to validate your account"
	return
}

func handleAccountFetch(payload json.RawMessage, u *User) (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Fetching account failed"

	// Unmarshal payload
	var account string
	if err = json.Unmarshal(payload, &account); err != nil {
		err = errors.Wrap(err, "unmarshaling failed")
		return
	}

	// Fetch user based on the account
	if _, err = storage.UserFetchWithAccount(account); err != nil && err != errNotFound {
		err = errors.Wrap(err, "fetching account failed")
		return
	}
	return
}

func handleAccountList(u *User) (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Listing accounts failed"

	// List accounts
	var es []*Account
	if es, err = storage.AccountList(u); err != nil {
		err = errors.Wrap(err, "listing account failed")
		return
	}

	// Build data
	var accounts = []string{}
	for _, e := range es {
		accounts = append(accounts, e.Addr)
	}
	data = accounts
	return
}

func handleReferences() (data interface{}, userErrorMsg string, err error) {
	// Init
	userErrorMsg = "Getting references failed"

	// Build data
	data = asticrypt.BodyReferences{
		GoogleClientID:     configuration.GoogleClientID,
		GoogleClientSecret: configuration.GoogleClientSecret,
		Now:                time.Now(),
	}
	return
}
