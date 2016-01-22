package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	conflib "github.com/yyoshiki41/go-gmail-drafts/lib"

	"golang.org/x/oauth2"
)

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	cmd := exec.Command("open", authURL)
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Browser Start Error: %v\n", err)
		fmt.Errorf("Go to the following link in your browser.\n")
	}
	fmt.Printf("URL: %v\nType the authorization code.\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return token
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	cacheFile := filepath.Join("./.credentials", url.QueryEscape("gmail_token.json"))

	config, err := conflib.CreateGmailConfig()
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	token := getTokenFromWeb(config)

	saveToken(cacheFile, token)
}
