// Package cmd contains common flags and routines for all CLIs.
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/notaryproject/notation/internal/envelope"
	"github.com/notaryproject/notation/pkg/configutil"
	"github.com/spf13/pflag"
)

var (
	PflagKey = &pflag.Flag{
		Name:      "key",
		Shorthand: "k",
		Usage:     "signing key name, for a key previously added to notation's key list.",
	}
	SetPflagKey = func(fs *pflag.FlagSet, p *string) {
		fs.StringVarP(p, PflagKey.Name, PflagKey.Shorthand, "", PflagKey.Usage)
	}

	PflagSignatureFormat = &pflag.Flag{
		Name:  "signature-format",
		Usage: "signature envelope format, options: 'jws', 'cose'",
	}
	SetPflagSignatureFormat = func(fs *pflag.FlagSet, p *string) {
		defaultSignatureFormat := envelope.JWS
		// load config to get signatureFormat
		config, err := configutil.LoadConfigOnce()
		if err == nil && config.SignatureFormat != "" {
			defaultSignatureFormat = config.SignatureFormat
		}

		fs.StringVar(p, PflagSignatureFormat.Name, defaultSignatureFormat, PflagSignatureFormat.Usage)
	}

	PflagTimestamp = &pflag.Flag{
		Name:      "timestamp",
		Shorthand: "t",
		Usage:     "timestamp the signed signature via the remote TSA",
	}
	SetPflagTimestamp = func(fs *pflag.FlagSet, p *string) {
		fs.StringVarP(p, PflagTimestamp.Name, PflagTimestamp.Shorthand, "", PflagTimestamp.Usage)
	}

	PflagExpiry = &pflag.Flag{
		Name:      "expiry",
		Shorthand: "e",
		Usage:     "optional expiry that provides a \"best by use\" time for the artifact. The duration is specified in minutes(m) and/or hours(h). For example: 12h, 30m, 3h20m",
	}
	SetPflagExpiry = func(fs *pflag.FlagSet, p *time.Duration) {
		fs.DurationVarP(p, PflagExpiry.Name, PflagExpiry.Shorthand, time.Duration(0), PflagExpiry.Usage)
	}

	PflagReference = &pflag.Flag{
		Name:      "reference",
		Shorthand: "r",
		Usage:     "original reference",
	}
	SetPflagReference = func(fs *pflag.FlagSet, p *string) {
		fs.StringVarP(p, PflagReference.Name, PflagReference.Shorthand, "", PflagReference.Usage)
	}

	PflagPluginConfig = &pflag.Flag{
		Name:      "plugin-config",
		Shorthand: "c",
		Usage:     "{key}={value} pairs that are passed as it is to a plugin, refer plugin's documentation to set appropriate values",
	}
	SetPflagPluginConfig = func(fs *pflag.FlagSet, p *[]string) {
		fs.StringArrayVarP(p, PflagPluginConfig.Name, PflagPluginConfig.Shorthand, nil, PflagPluginConfig.Usage)
	}
)

// KeyValueSlice is a flag with type int
type KeyValueSlice interface {
	Set(value string) error
	String() string
}

func ParseFlagPluginConfig(config []string) (map[string]string, error) {
	pluginConfig := make(map[string]string, len(config))
	for _, pair := range config {
		key, val, found := strings.Cut(pair, "=")
		if !found || key == "" || val == "" {
			return nil, fmt.Errorf("could not parse flag %s: key-value pair requires \"=\" as separator", PflagPluginConfig.Name)
		}
		pluginConfig[key] = val
	}
	return pluginConfig, nil
}
