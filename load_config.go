package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
)

const filePath = "config/client_secret.json"

// readClient reads config json file.
func readClient() ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

// CreateGmailConfig returns OAuth 2.0 config.
func CreateGmailConfig() (*oauth2.Config, error) {
	b, err := readClient()
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, gmail.GmailComposeScope, gmail.GmailModifyScope)
}
