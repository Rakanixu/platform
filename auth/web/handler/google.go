package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kazoup/platform/datasource/srv/proto/datasource"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8082/desktop/GoogleCallback",
		ClientID:     "928848534435-kjubrqvl1sp50sfs3icemj2ma6v2an5j.apps.googleusercontent.com",
		ClientSecret: "zZAQz3zP5xnpLaA1S_q6YNhy",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "randomsdsdahfoashfouahsfohasofhoashfaf"
)

/* {
 "id": "101000145849438728639",
 "email": "radekdymacz@gmail.com",
 "verified_email": true,
 "name": "radek dymacz",
 "given_name": "radek",
 "family_name": "dymacz",
 "link": "https://plus.google.com/101000145849438728639",
 "picture": "https://lh3.googleusercontent.com/-LAOs29oR6RA/AAAAAAAAAAI/AAAAAAAAHCY/iCzIUKZ5mGo/photo.jpg",
 "gender": "male",
 "locale": "en-GB"
}*/
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	userInfo := new(GoogleUserInfo)
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	//response, err := http.Get("https://www.googleapis.com/drive/v3/files?corpus=user&key=" + token.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}

	t := &go_micro_srv_datasource.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.String(),
	}
	c := go_micro_srv_datasource.NewDataSourceClient("go.micro.srv.desktop", nil)
	endpoint := &go_micro_srv_datasource.Endpoint{
		Url:   fmt.Sprintf("googledrive://%s", userInfo.Email),
		Token: t,
	}
	req := &go_micro_srv_datasource.CreateRequest{
		Endpoint: endpoint,
	}

	_, err = c.Create(context.TODO(), req)
	if err != nil {
		fmt.Fprintf(w, "Error adding data source %s \n", err.Error())
	}

	fmt.Fprintf(w, "Status: New Google Drive added. You can close the window.\n Info %s Resp: %s \n", userInfo, contents)
}
