package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

type SlackTeamInfoResponse struct {
	OK   bool          `json:"ok"`
	Team SlackTeamInfo `json:"team"`
}
type SlackTeamInfo struct {
	Name        string            `json:"name"`
	Domain      string            `json:"domain"`
	EmailDomain string            `json:"email_domain"`
	Icon        map[string]string `json:"icon"`
}

func HandleSlackLogin(w http.ResponseWriter, r *http.Request) {
	url := globals.NewSlackOauthConfig().AuthCodeURL(globals.OauthStateString, oauth2.AccessTypeOffline)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleSlackCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != globals.OauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", globals.OauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := globals.NewSlackOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	response, err := http.Get("https://slack.com/api/team.info?token=" + token.AccessToken)
	defer response.Body.Close()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {

		fmt.Fprintf(w, err.Error())
	}
	sr := new(SlackTeamInfoResponse)
	if err := json.Unmarshal(contents, &sr); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	if !sr.OK {
		fmt.Fprintf(w, "Error $s", sr)
	}
	url := fmt.Sprintf("slack://%s", sr.Team.Name)
	if err := SaveDatasource(url, token); err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Close window, not reacheable by electron any more
	fmt.Fprintf(w, "%s", `
		<script>
		'use stric';
			(function() {
				window.close();
			}());
		</script>
	`)
}
