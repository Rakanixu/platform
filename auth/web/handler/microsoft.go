package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/onedrive"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func HandleMicrosoftLogin(w http.ResponseWriter, r *http.Request) {
	url := globals.NewMicrosoftOauthConfig().AuthCodeURL(r.URL.Query().Get("user"), oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleMicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if len(state) == 0 {
		fmt.Printf("invalid oauth state, got '%s'\n", state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := globals.NewMicrosoftOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get user name
	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	u := globals.OneDriveEndpoint + "drives"
	req, err := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()

	var drivesRsp *onedrive.DrivesListResponse
	if err := json.NewDecoder(res.Body).Decode(&drivesRsp); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	url := fmt.Sprintf("onedrive://%s", drivesRsp.Value[0].Owner.User.DisplayName)

	if err := SaveDatasource(globals.NewSystemContext(), state, url, token); err != nil {
		fmt.Fprintf(w, "Error adding data source %s \n", err.Error())
	}

	fmt.Fprintf(w, "%s", `
		<script>
		'use stric';
			(function() {
				window.close();
			}());
		</script>
	`)
}
