package consul

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/config"
)

// Config is a Consul service discovery configuration
type Config struct {
	Enabled      bool     `toml:"enabled" override:"enabled"`
	Name         string   `toml:"name" override:"name"`
	Address      string   `toml:"address" override:"address"`
	Token        string   `toml:"token" override:"token,redact"`
	Datacenter   string   `toml:"datacenter" override:"datacenter"`
	TagSeparator string   `toml:"tag-separator" override:"tag-separator"`
	Scheme       string   `toml:"scheme" override:"scheme"`
	Username     string   `toml:"username" override:"username"`
	Password     string   `toml:"password" override:"password,redact"`
	Services     []string `toml:"services" override:"services"`
	// Path to CA file
	SSLCA string `toml:"ssl-ca" override:"ssl-ca"`
	// Path to host cert file
	SSLCert string `toml:"ssl-cert" override:"ssl-cert"`
	// Path to cert key file
	SSLKey string `toml:"ssl-key" override:"ssl-key"`
	// SSLServerName is used to verify the hostname for the targets.
	SSLServerName string `toml:"ssl-server-name" override:"ssl-server-name"`
	// Use SSL but skip chain & host verification
	InsecureSkipVerify bool `toml:"insecure-skip-verify" override:"insecure-skip-verify"`
}

// NewConfig creates a new Consul discovery with default values
func NewConfig() Config {
	return Config{
		Address:      "127.0.0.1:8500",
		TagSeparator: ",",
		Scheme:       "http",
		Services:     []string{},
	}
}

// ApplyConditionalDefaults adds defaults to Consul
func (c *Config) ApplyConditionalDefaults() {
	if c.TagSeparator == "" {
		c.TagSeparator = ","
	}
	if c.Scheme == "" {
		c.Scheme = "http"
	}
}

// Validate validates the consul configuration
func (c Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("consul discovery must be given a name")
	}
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("consul discovery requires a server address")
	}
	return nil
}

// Prom creates a prometheus configuration from Consul configuration
func (c Config) Prom(conf *config.ScrapeConfig) {
	conf.ServiceDiscoveryConfig.ConsulSDConfigs = []*config.ConsulSDConfig{
		&config.ConsulSDConfig{
			Server:       c.Address,
			Token:        c.Token,
			Datacenter:   c.Datacenter,
			TagSeparator: c.TagSeparator,
			Scheme:       c.Scheme,
			Username:     c.Username,
			Password:     c.Password,
			Services:     c.Services,
			TLSConfig: config.TLSConfig{
				CAFile:             c.SSLCA,
				CertFile:           c.SSLCert,
				KeyFile:            c.SSLKey,
				ServerName:         c.SSLServerName,
				InsecureSkipVerify: c.InsecureSkipVerify,
			},
		},
	}
}

// Service return discoverer type
func (c Config) Service() string {
	return "consul"
}

// ID returns the discoverers name
func (c Config) ID() string {
	return c.Name
}
