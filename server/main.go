package main

import (
	"flag"

	"os"
	"os/signal"
	"syscall"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astimysql"
	"github.com/asticode/go-astipatch"
	"github.com/asticode/go-astitools/flag"
	"github.com/jmoiron/sqlx"
)

var channelQuit = make(chan bool)

func main() {
	// Parse flags
	var s = astiflag.Subcommand()
	flag.Parse()

	// Build configuration
	var c = newConfiguration()

	// Build logger
	astilog.SetLogger(astilog.New(c.Logger))

	// Build db
	var db *sqlx.DB
	var err error
	if db, err = astimysql.New(c.MySQL); err != nil {
		astilog.Fatalf("%s while creating db", err)
	}

	// Build patcher
	var p = astipatch.NewPatcherSQL(db, astipatch.NewStorerSQL(db))

	// Build storage
	storage = newStorageMySQL(db)

	// Handle signals
	handleSignals()

	// Switch on subcommand
	switch s {
	case "db-init":
		if err = p.Init(); err != nil {
			astilog.Fatal(err)
		}
		astilog.Info("db-init successful")
	case "db-migrate", "db-rollback":
		// Load patches
		if err = p.Load(c.Patcher); err != nil {
			astilog.Fatal(err)
		}

		// Exec
		if s == "db-migrate" {
			if err = p.Patch(); err != nil {
				astilog.Fatal(err)
			}
		} else {
			if err = p.Rollback(); err != nil {
				astilog.Fatal(err)
			}
		}
		astilog.Infof("%s successful", s)
	default:
		// Serve
		if err := serve(c.AddrLocal, c.PathResources); err != nil {
			astilog.Fatalf("%s while serving", err)
		}

		// Wait
		wait()
	}
}

func handleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			astilog.Debugf("Received signal %s", sig)
			channelQuit <- true
		}
	}()
}

func wait() {
	for {
		select {
		case <-channelQuit:
			return
		}
	}
}
