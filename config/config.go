package config

import (
	"encoding/base64"
)

type Config struct {
	URL     string                 `json:"JiraUrl"`
	Creds   Creds                  `json:"Credentials"`
	Filters map[string]interface{} `json:"Filters"`
	Fields  []string               `json:"FieldsToRetrieve"`
}

type Creds struct {
	Username string
	Password string
}

func NewCreds(username string, password string) Creds {
	return Creds{Username: username, Password: password}
}
func NewConfig(url string, filters map[string]interface{}, fields []string, output string) *Config {
	return &Config{URL: url, Filters: filters, Fields: fields}
}

func (c *Config) SetCreds(creds Creds) {
	c.Creds = creds
}
func (c *Config) AuthToken() string {
	return encodeStringToBase64(c.Creds.Username + ":" + c.Creds.Password)
}

func encodeStringToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}
