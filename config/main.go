package config

import (
	"io/ioutil"

	"github.com/akerl/timber/v2/log"
	"github.com/ghodss/yaml"
)

var logger = log.NewLogger("simplefin-exporter.config")

// Config defines the available configuration options
type Config struct {
	Port      int      `json:"port"`
	Interval  int      `json:"interval"`
	AccessURL string   `json:"access_url"`
	Ignore    []string `json:"ignore"`
}

// LoadConfig creates a config from a file path
func LoadConfig(customPath string) (Config, error) {
	var c Config

	path := customPath
	logger.InfoMsgf("loading config from path: %s", path)

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(contents, &c)
	logger.DebugMsgf("loaded config: %v+", c)
	return c, err
}
