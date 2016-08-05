package handler

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/slack"
)

var (
	slackOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8082/auth/slack/callback",
		ClientID:     "2506087186.66729631906",
		ClientSecret: "53ea1f0afa4560b7e070964fb2b0c5d6",
		Scopes:       []string{"files:read", "files:write:user"},
		Endpoint:     slack.Endpoint,
	}
	// Some random string, random for each request
	oauthSlackStateString = "randomsdsdahfoashfouahsfohasofhoashfaf"
)

func HandleSlackLogin(w http.ResponseWriter, r *http.Request) {
	url := slackOauthConfig.AuthCodeURL(oauthSlackStateString)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleSlackCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := slackOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	//response, err := http.Get("https://www.googleapis.com/drive/v3/files?corpus=user&key=" + token.AccessToken)
	//defer response.Body.Close()
	//contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", token.AccessToken)
}
