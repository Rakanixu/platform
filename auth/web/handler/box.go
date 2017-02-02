package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/oauth2"
)

//BoxUser struct
type BoxUser struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Login         string `json:"login"`
	CreatedAt     string `json:"created_at"`
	ModifiedAt    string `json:"modified_at"`
	Language      string `json:"language"`
	SpaceAmount   int64  `json:"space_amount"`
	SpaceUsed     int    `json:"space_used"`
	MaxUploadSize int    `json:"max_upload_size"`
	Status        string `json:"status"`
	JobTitle      string `json:"job_title"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	AvatarURL     string `json:"avatar_url"`
}

//HandleBoxLogin Oauth
func HandleBoxLogin(w http.ResponseWriter, r *http.Request) {
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
	url := globals.NewBoxOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//HandleBoxCallback handle response from BOX
func HandleBoxCallback(w http.ResponseWriter, r *http.Request) {
	var bu *BoxUser

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
	token, err := globals.NewBoxOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxAccountEndpoint, nil)
	if err != nil {
		CloseBrowserWindow(w, r)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	rsp, err := c.Do(req)
	if err != nil {
		log.Printf("Getting user account failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}
	defer rsp.Body.Close()

	contents, err := ioutil.ReadAll(rsp.Body)
	if err := json.Unmarshal(contents, &bu); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}
	url := fmt.Sprintf("box://%s", bu.Login)

	if err := SaveDatasource(uCtx, uID, url, token); err != nil {
		fmt.Fprintf(w, "Error adding data source %s \n", err.Error())
	}

	CloseBrowserWindow(w, r)

	if err := PublishNotification(uID); err != nil {
		fmt.Fprintf(w, "Error publishing notification msg %s \n", err.Error())
	}
}
