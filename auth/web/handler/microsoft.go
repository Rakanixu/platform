package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

var (
	AzureConnectEndpoint = oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/token",
	}
	redirect_uri         = "http://localhost:8082/auth/microsoft/callback"
	resource             = "https://graph.microsoft.com/"
	microsoftOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8082/auth/microsoft/callback",
		ClientID:     "6544d648-c29c-43f1-8e3c-7a903e8524c4",
		ClientSecret: "lGFfMwyErTWXDjTFRDYeypIctjfjn1CcbgtWRfRv/kw=",
		Endpoint:     AzureConnectEndpoint,
	}
	// Some random string, random for each request
	oauthMStateString = "randomsdsdahfoashfouahsfohasofhoashfaf"
)

func HandleMicrosoftLogin(w http.ResponseWriter, r *http.Request) {
	url := microsoftOauthConfig.AuthCodeURL(oauthStateString)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleMicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthMStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	log.Print(r.URL)
	resp, err := http.PostForm(AzureConnectEndpoint.TokenURL, url.Values{"grant_type": {"authorization_code"}, "redirect_uri": {redirect_uri}, "client_id": {microsoftOauthConfig.ClientID}, "client_secret": {microsoftOauthConfig.ClientSecret}, "code": {code}, "resource": {resource}})
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	//response, err := http.Get("https://www.googleapis.com/drive/v3/files?corpus=user&key=" + token.AccessToken)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "Token: %s\n", contents)
}
