package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

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

func HandleBoxLogin(w http.ResponseWriter, r *http.Request) {
	t := []byte(r.URL.Query().Get("user"))                          // String to encrypt
	nt, err := globals.Encrypt([]byte(globals.ENCRYTION_KEY_32), t) // Encryption
	if err != nil {
		log.Printf("Encryption failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Code conversion from bytes to hexadecimal string to be send over the wire
	url := globals.NewBoxOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleBoxCallback(w http.ResponseWriter, r *http.Request) {
	var bu *BoxUser

	euID, err := hex.DecodeString(r.FormValue("state"))                 // Convert the code we sent in hex format to bytes
	uID, err := globals.Decrypt([]byte(globals.ENCRYTION_KEY_32), euID) // Decrypt the bytes into bytes --> string(bytes) was the encrypted string
	if err != nil {
		log.Printf("Decryption failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(uID) == 0 {
		fmt.Printf("invalid oauth state, got '%s'\n", uID)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := globals.NewBoxOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxAccountEndpoint, nil)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	rsp, err := c.Do(req)
	if err != nil {
		log.Printf("Getting user account failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer rsp.Body.Close()

	contents, err := ioutil.ReadAll(rsp.Body)
	if err := json.Unmarshal(contents, &bu); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}
	url := fmt.Sprintf("box://%s", bu.Login)

	if err := SaveDatasource(globals.NewSystemContext(), string(uID), url, token); err != nil {
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
