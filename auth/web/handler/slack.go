package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

//SlackTeamInfoResponse  data
type SlackTeamInfoResponse struct {
	OK   bool          `json:"ok"`
	Team SlackTeamInfo `json:"team"`
}

//SlackTeamInfo data
type SlackTeamInfo struct {
	Name        string            `json:"name"`
	Domain      string            `json:"domain"`
	EmailDomain string            `json:"email_domain"`
	Icon        map[string]string `json:"icon"`
}

//HandleSlackLogin Slack oauth2 redirect
func HandleSlackLogin(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	uuid, err := utils.NewUUID()
	if err != nil {
		fmt.Printf("UUID generation failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	if err := SaveTmpToken(uuid, jwt); err != nil {
		fmt.Printf("Save tmp token failed with error: '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	nt, err := utils.Encrypt([]byte(globals.ENCRYTION_KEY_32), []byte(uuid)) // Encryption
	if err != nil {
		fmt.Printf("Encryption failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	// Code conversion from bytes to hexadecimal string to be send over the wire
	url := globals.NewSlackOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//HandleSlackCallback Slack response handler
func HandleSlackCallback(w http.ResponseWriter, r *http.Request) {
	euID, err := hex.DecodeString(r.FormValue("state"))                // Convert the code we sent in hex format to bytes
	uuid, err := utils.Decrypt([]byte(globals.ENCRYTION_KEY_32), euID) // Decrypt the bytes into bytes --> string(bytes) was the encrypted string
	if err != nil {
		fmt.Printf("Decryption failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	if len(uuid) == 0 {
		fmt.Printf("invalid oauth state, got '%s'\n", uuid)
		CloseBrowserWindow(w, r)
		return
	}

	// Get userId and context
	uID, uCtx, err := RetrieveUserAndContextFromUUID(string(uuid))
	if err != nil {
		fmt.Printf("Retrieving user_id and context failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	code := r.FormValue("code")
	token, err := globals.NewSlackOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}
	response, err := http.Get("https://slack.com/api/team.info?token=" + token.AccessToken)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf(err.Error())
		CloseBrowserWindow(w, r)
		return
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf(err.Error())
		CloseBrowserWindow(w, r)
		return
	}
	sr := new(SlackTeamInfoResponse)
	if err := json.Unmarshal(contents, &sr); err != nil {
		fmt.Printf(err.Error())
		CloseBrowserWindow(w, r)
		return
	}
	if !sr.OK {
		fmt.Printf("Error %v", sr)
		CloseBrowserWindow(w, r)
		return
	}
	url := fmt.Sprintf("slack://%s", sr.Team.Name)
	if err := SaveDatasource(uCtx, uID, url, token); err != nil {
		fmt.Fprintf(w, err.Error())
		CloseBrowserWindow(w, r)
		return
	}

	// Close window, not reacheable by electron any more
	CloseBrowserWindow(w, r)

	if err := PublishNotification(uCtx, string(uID)); err != nil {
		fmt.Fprintf(w, "Error publishing notification msg %s \n", err.Error())
	}
}
