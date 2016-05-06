package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	PagesDir    string `yaml:"pages_dir"`
	Debug       bool   `yaml:"debug"`
}

// C is the global conf
var C conf

// Load loads the configuration file into C and fills it with default
// configuration if fields are empty
// (Listen to 127.0.0.1, on port 8080, watching the "pages" directory)
func Load(fp string) error {
	var err error
	var c []byte

	if c, err = ioutil.ReadFile(fp); err != nil {
		return err
	}
	if C.PagesDir == "" {
		C.PagesDir = "pages"
	}
	if C.Host == "" {
		C.Host = "127.0.0.1"
	}
	if C.Port == 0 {
		C.Port = 8080
	}
	return yaml.Unmarshal(c, &C)
}
