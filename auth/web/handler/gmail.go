package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"github.com/micro/go-micro/client"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

//HandleGmailLogin handle Gmail login
func HandleGmailLogin(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	uuid, err := utils.NewUUID()
	if err != nil {
		fmt.Printf("UUID generation failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	if err := SaveTmpToken(uuid, jwt); err != nil {
		log.Printf("Save tmp token failed with error: '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	nt, err := utils.Encrypt([]byte(globals.ENCRYTION_KEY_32), []byte(uuid)) // Encryption
	if err != nil {
		log.Printf("Encryption failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	// Code conversion from bytes to hexadecimal string to be send over the wire
	url := globals.NewGmailOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//HandleGmailCallback Gmail response handler
func HandleGmailCallback(w http.ResponseWriter, r *http.Request) {
	userInfo := new(GoogleUserInfo)

	euID, err := hex.DecodeString(r.FormValue("state"))                // Convert the code we sent in hex format to bytes
	uuid, err := utils.Decrypt([]byte(globals.ENCRYTION_KEY_32), euID) // Decrypt the bytes into bytes --> string(bytes) was the encrypted string
	if err != nil {
		log.Printf("Decryption failed with '%s'\n", err)
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
		log.Printf("Retrieving user_id and context failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	code := r.FormValue("code")
	token, err := globals.NewGmailOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}
	url := fmt.Sprintf("gmail://%s", userInfo.Email)

	if err := SaveDatasource(uCtx, client.NewClient(), uID, url, token); err != nil {
		fmt.Fprintf(w, "Error adding data source %s \n", err.Error())
	}

	CloseBrowserWindow(w, r)

	if err := PublishNotification(uCtx, uID); err != nil {
		fmt.Fprintf(w, "Error publishing notification msg %s \n", err.Error())
	}
}
