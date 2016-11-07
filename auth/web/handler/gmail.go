package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleGmailLogin(w http.ResponseWriter, r *http.Request) {
	t := []byte(r.URL.Query().Get("user"))                          // String to encrypt
	nt, err := globals.Encrypt([]byte(globals.ENCRYTION_KEY_32), t) // Encryption
	if err != nil {
		log.Printf("Encryption failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Code conversion from bytes to hexadecimal string to be send over the wire
	url := globals.NewGmailOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGmailCallback(w http.ResponseWriter, r *http.Request) {
	userInfo := new(GoogleUserInfo)

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
	token, err := globals.NewGmailOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}
	url := fmt.Sprintf("gmail://%s", userInfo.Email)

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
