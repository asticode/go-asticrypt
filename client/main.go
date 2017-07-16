package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"github.com/asticode/go-asticrypt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	clientPrivateKey   *asticrypt.PrivateKey
	accounts           = make(map[string]string)
	httpClient         = &http.Client{}
	googleClientID     string
	googleClientSecret string
	now                time.Time
	pathConfiguration  string
	pathExecutable     string
	serverPublicKey    *asticrypt.PublicKey
	ServerPublicAddr   string
	Version            string
)

//go:generate go-bindata -pkg $GOPACKAGE -o resources.go resources/...
func main() {
	// TODO For test purposes
	ServerPublicAddr = "http://127.0.0.1:4000"

	// Parse flags
	flag.Parse()

	// Build logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Fetch executable path
	var err error
	if pathExecutable, err = os.Executable(); err != nil {
		astilog.Fatal(errors.Wrap(err, "fetching executable path failed"))
	}
	pathExecutable = filepath.Dir(pathExecutable)

	// Build paths
	pathConfiguration = filepath.Join(pathExecutable, "local.toml")

	// Run bootstrap
	if err = bootstrap.Run(bootstrap.Options{
		AstilectronOptions: astilectron.Options{
			AppName: "Asticrypt",
		},
		Debug:          true,
		Homepage:       "index.html",
		MessageHandler: handleMessages,
		// RestoreAssets:  RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(720),
			Width:           astilectron.PtrInt(720),
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
