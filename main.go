package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := filepath.Join("./.credentials", url.QueryEscape("gmail_token.json"))

	token, err := tokenFromFile(cacheFile)
	if err != nil {
		log.Fatalf("Unable to retrieve token from cached credential file. %v", err)
	}
	return config.Client(ctx, token)
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// loadBodyFile returns mail body in the json template.
func loadBodyFile() (map[string]interface{}, error) {
	c := make(map[string]interface{})
	f, err := ioutil.ReadFile("./config/draft_tmpl.json")
	if err != nil {
		return c, err
	}
	json.Unmarshal(f, &c)
	return c, nil
}

// Convert UTF-8 to ISO2022JP
func toISO2022JP(str string) ([]byte, error) {
	reader := strings.NewReader(str)
	transformer := japanese.ISO2022JP.NewEncoder()

	return ioutil.ReadAll(transform.NewReader(reader, transformer))
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("./config/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b,
		gmail.GmailComposeScope, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}
	_ = srv

	draftMap, err := loadBodyFile()
	if err != nil {
		log.Fatalf("Unable to parse draft template file to map structure: %v", err)
	}

	toStr := "to: "
	if to, ok := draftMap["to"].(string); ok {
		toStr += to
	}
	toStr += "\n"

	subStr := "subject: "
	if subject, ok := draftMap["subject"].(string); ok {
		subStr += subject
	}
	subStr += "\n"
	t := time.Now()
	subStr = strings.Replace(subStr, "{{today}}", t.Format("01/02"), -1)

	msgStr := "\n"
	if message, ok := draftMap["message"].(string); ok {
		msgStr += message
	}
	msgStr = strings.Replace(msgStr, "{{today}}", t.Format("01/02"), -1)

	header, _ := toISO2022JP(toStr + subStr)
	msg := []byte(msgStr)
	bodyBytes := append(header, msg...)
	message := gmail.Message{}
	message.Raw = base64.URLEncoding.EncodeToString(bodyBytes)

	user := "me"
	draft := gmail.Draft{
		Message: &message,
	}
	_, err = srv.Users.Drafts.Create(user, &draft).Do()
	if err != nil {
		log.Fatalf("Unable to create drafts: %v", err)
	}
}
