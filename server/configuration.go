package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astimysql"
	"github.com/asticode/go-astipatch"
	"github.com/imdario/mergo"
	"github.com/rs/xlog"
)

// Flags
var (
	addrLocal          = flag.String("l", "", "the local addr")
	addrPublic         = flag.String("p", "", "the public addr")
	configPath         = flag.String("c", "", "the config path")
	googleClientID     = flag.String("gci", "", "the google client id")
	googleClientSecret = flag.String("gcs", "", "the google client secret")
	pathResources      = flag.String("r", "", "the resources path")
)

// Configuration represents a configuration
type Configuration struct {
	AddrLocal          string                  `toml:"addr_local"`
	AddrPublic         string                  `toml:"addr_public"`
	GoogleClientID     string                  `toml:"google_client_id"`
	GoogleClientSecret string                  `toml:"google_client_secret"`
	Logger             astilog.Configuration   `toml:"logger"`
	MySQL              astimysql.Configuration `toml:"mysql"`
	Patcher            astipatch.Configuration `toml:"patcher"`
	PathResources      string                  `toml:"path_resources"`
}

// newConfiguration creates a new configuration object
func newConfiguration() (c Configuration) {
	// Global config
	var gc = Configuration{
		Logger: astilog.Configuration{
			AppName: "go-astimail-server",
		},
	}

	// Local config
	if *configPath != "" {
		// Decode local config
		if _, err := toml.DecodeFile(*configPath, &gc); err != nil {
			xlog.Fatalf("%v while decoding the config path %s", err, *configPath)
		}
	}

	// Flag config
	c = Configuration{
		AddrLocal:          *addrLocal,
		AddrPublic:         *addrPublic,
		GoogleClientID:     *googleClientID,
		GoogleClientSecret: *googleClientSecret,
		Logger:             astilog.FlagConfig(),
		MySQL:              astimysql.FlagConfig(),
		Patcher:            astipatch.FlagConfig(),
		PathResources:      *pathResources,
	}

	// Merge configs
	if err := mergo.Merge(&c, gc); err != nil {
		xlog.Fatalf("%v while merging configs", err)
	}
	return
}
