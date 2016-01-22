# go-gmail-drafts

Create a new draft using the Gmail API.

## Installation

```shell
$ go get -u github.com/yyoshiki41/go-gmail-drafts/...
```

## Preparation
- Create project in the [Google Developers Console] (https://console.developers.google.com/flows/enableapi?apiid=gmail) and turn on the Gmail API.
- Download credential file (`client_secret.json`) .

## Configuration
### 1. Put your credential file

```shell
$ cd $GOPATH/src/github.com/yyoshiki41/go-gmail-drafts
$ mv /path/to/client_secret.json config/client_secret.json
```

### 2. Save access token to cachefile

Open the browser automatically when exec `go run` command.

After authorize Google APIs, paste authorization code into the command-line prompt.

```shell
$ go run savetoken/main.go
URL: https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=hogehoge.apps.googleusercontent.com&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fgmail.compose+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fgmail.modify&state=state-token
Type the authorization code.
### Paste authorization code here! ###
Saving credential file to: .credentials/gmail_token.json

$ cat .credentials/gmail_token.json | jq .
{
  "access_token": "xxxxx",
  "token_type": "Bearer",
  "refresh_token": "xxxxx",
  "expiry": "2016-01-22T11:11:47.804107777+09:00"
}
```

### Example

Create a gmail template.

```shell
$ cp config/draft_tmpl.json.skel config/draft_tmpl.json
$ vim config/draft_tmpl.json
{
  "to": "xxx@gmail.com",
  "subject": "Daily Reports",
  "message": "Hello!\nThis is a draft."
}
```

Run!

```shell
$ go run main.go
```

## License 
The MIT License

## Author
Yoshiki Nakagawa
