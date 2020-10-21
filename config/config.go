package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Configuration struct {
	JiraURL          string                 `json:"JiraUrl"`
	Credentials      Credentials            `json:"Credentials"`
	Filters          map[string]interface{} `json:"Filters"`
	FieldsToRetrieve []string               `json:"FieldsToRetrieve"`
	DownloadPath     string                 `json:"DownloadPath"`
	AuthToken        string
}

type Credentials struct {
	Username string
	Password string
}

func ensureFile(confgFile string) (error, string) {
	if confgFile == "" {
		return errors.New("empty config not allowed"), ""
	}

	if !filepath.IsAbs(confgFile) {
		wd, _ := os.Getwd()
		confgFile = filepath.Join(wd, confgFile)
	}

	// ensure configFile is a valid file
	i, err := os.Stat(confgFile)
	if err != nil {
		return err, ""
	}

	if i.IsDir() {
		return errors.Wrapf(err, "invalid config file"), ""
	}

	return nil, confgFile
}

func encodeStringToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

func New(confgFile string) (error, *Configuration) {
	var c *Configuration

	err, confgFile := ensureFile(confgFile)
	if err != nil {
		return err, nil
	}

	fmt.Println(" Fetching data based on the configuration file => " + "'" + confgFile + "'")

	jsonFile, err := os.Open(confgFile)
	if err != nil {
		return errors.Wrapf(err, "failed to open config file"), nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal([]byte(byteValue), &c)
	if err != nil {
		return errors.Wrapf(err, "failed to parse config file"), nil
	}

	c.AuthToken = encodeStringToBase64(c.Credentials.Username + ":" + c.Credentials.Password)

	return nil, c
}
