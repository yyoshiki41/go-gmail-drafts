package config

import (
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
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
