package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/onedrive"
	"golang.org/x/oauth2"
)

//HandleMicrosoftLogin Microsoft oauth2 redirect
func HandleMicrosoftLogin(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	uuid, err := globals.NewUUID()
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

	nt, err := globals.Encrypt([]byte(globals.ENCRYTION_KEY_32), []byte(uuid)) // Encryption
	if err != nil {
		log.Printf("Encryption failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	// Code conversion from bytes to hexadecimal string to be send over the wire
	url := globals.NewMicrosoftOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//HandleMicrosoftCallback M$ response handler
func HandleMicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	euID, err := hex.DecodeString(r.FormValue("state"))                  // Convert the code we sent in hex format to bytes
	uuid, err := globals.Decrypt([]byte(globals.ENCRYTION_KEY_32), euID) // Decrypt the bytes into bytes --> string(bytes) was the encrypted string
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
	token, err := globals.NewMicrosoftOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
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
		CloseBrowserWindow(w, r)
		return
	}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		CloseBrowserWindow(w, r)
		return
	}
	defer res.Body.Close()

	var drivesRsp *onedrive.DrivesListResponse
	if err := json.NewDecoder(res.Body).Decode(&drivesRsp); err != nil {
		log.Println(err)
		CloseBrowserWindow(w, r)
		return
	}

	url := fmt.Sprintf("onedrive://%s", drivesRsp.Value[0].Owner.User.DisplayName)

	if err := SaveDatasource(uCtx, uID, url, token); err != nil {
		fmt.Fprintf(w, "Error adding data source %s \n", err.Error())
	}

	CloseBrowserWindow(w, r)

	if err := PublishNotification(uCtx, string(uID)); err != nil {
		fmt.Fprintf(w, "Error publishing notification msg %s \n", err.Error())
	}
}
