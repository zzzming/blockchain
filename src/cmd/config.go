package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/apex/log"
	"github.com/ghodss/yaml"
)

type Config struct {
	WalletDir    string `json:"WalletDir" yaml:"WalletDir"`
	DatabaseDir  string `json:"DatabaseDir" yaml:"DatabaseDir"`
	Port         string `json:"Port" yaml:"Port"`
	POWDifficult int    `json:"POWDifficult" yaml:"POWDifficult"`
}

// ReadConfigFile reads configuration file.
func ReadConfigFile(configFile string) (Config, error) {
	var config Config
	fileBytes, err := ioutil.ReadFile(filepath.Clean(configFile))
	if err != nil {
		log.Infof("failed to load configuration file %s", configFile)
		return config, err
	}

	if hasJSONPrefix(fileBytes) {
		err = json.Unmarshal(fileBytes, &config)
		if err != nil {
			return config, err
		}
	} else {
		err = yaml.Unmarshal(fileBytes, &config)
		if err != nil {
			return config, err
		}
	}

	// Next section allows env variable overwrites config file value
	fields := reflect.TypeOf(config)
	// pointer to struct
	values := reflect.ValueOf(&config)
	// struct
	st := values.Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i).Name
		f := st.FieldByName(field)

		if f.Kind() == reflect.String {
			envV := os.Getenv(field)
			if len(envV) > 0 && f.IsValid() && f.CanSet() {
				f.SetString(strings.TrimSuffix(envV, "\n")) // ensure no \n at the end of line that was introduced by loading k8s secrete file
			}
			os.Setenv(field, f.String())
		}
	}

	return config, nil
}

var jsonPrefix = []byte("{")

func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}

// AssignString returns the first non-empty string
// It is equivalent the following in Javascript
// var value = val0 || val1 || val2 || default
func AssignString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
