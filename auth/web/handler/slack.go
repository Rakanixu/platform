package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/oauth2"
	"io/ioutil"
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
	t := []byte(r.URL.Query().Get("user"))                          // String to encrypt
	nt, err := globals.Encrypt([]byte(globals.ENCRYTION_KEY_32), t) // Encryption
	if err != nil {
		fmt.Printf("Encryption failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Code conversion from bytes to hexadecimal string to be send over the wire
	url := globals.NewSlackOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleSlackCallback(w http.ResponseWriter, r *http.Request) {
	euID, err := hex.DecodeString(r.FormValue("state"))                 // Convert the code we sent in hex format to bytes
	uID, err := globals.Decrypt([]byte(globals.ENCRYTION_KEY_32), euID) // Decrypt the bytes into bytes --> string(bytes) was the encrypted string
	if err != nil {
		fmt.Printf("Decryption failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(uID) == 0 {
		fmt.Printf("invalid oauth state, got '%s'\n", uID)
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
	if err := SaveDatasource(globals.NewSystemContext(), string(uID), url, token); err != nil {
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

	if err := PublishNotification(string(uID)); err != nil {
		fmt.Fprintf(w, "Error publishing notification msg %s \n", err.Error())
	}
}
