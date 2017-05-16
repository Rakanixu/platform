package handler

import (
	"bytes"
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

//DropboxAccount data
type DropboxAccount struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName       string `json:"given_name"`
		Surname         string `json:"surname"`
		FamiliarName    string `json:"familiar_name"`
		DisplayName     string `json:"display_name"`
		AbbreviatedName string `json:"abbreviated_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Disabled      bool   `json:"disabled"`
	Locale        string `json:"locale"`
	ReferralLink  string `json:"referral_link"`
	IsPaired      bool   `json:"is_paired"`
	AccountType   struct {
		Tag string `json:".tag"`
	} `json:"account_type"`
	Country string `json:"country"`
	Team    struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		SharingPolicies struct {
			SharedFolderMemberPolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_folder_member_policy"`
			SharedFolderJoinPolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_folder_join_policy"`
			SharedLinkCreatePolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_link_create_policy"`
		} `json:"sharing_policies"`
	} `json:"team"`
	TeamMemberID string `json:"team_member_id"`
}

//HandleDropboxLogin redirect
func HandleDropboxLogin(w http.ResponseWriter, r *http.Request) {
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
	// Dropbox does not follow oauth2 spec. They do define a new flag force_reapprove Boolean.
	url := globals.NewDropboxOauthConfig().AuthCodeURL(fmt.Sprintf("%0x", nt), oauth2.SetAuthURLParam("force_reapprove", "true"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//HandleDropboxCallback repsonse from Dropbox
func HandleDropboxCallback(w http.ResponseWriter, r *http.Request) {
	var da *DropboxAccount

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
	token, err := globals.NewDropboxOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}

	c := &http.Client{}
	b := []byte(`{"account_id":"` + token.Extra("account_id").(string) + `"}`)

	req, err := http.NewRequest("POST", globals.DropboxAccountEndpoint, bytes.NewBuffer(b))
	if err != nil {
		CloseBrowserWindow(w, r)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		log.Printf("Getting user account failed with '%s'\n", err)
		CloseBrowserWindow(w, r)
		return
	}
	defer rsp.Body.Close()

	contents, err := ioutil.ReadAll(rsp.Body)
	if err := json.Unmarshal(contents, &da); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}
	url := fmt.Sprintf("dropbox://%s", da.Email)

	if err := SaveDatasource(uCtx, client.NewClient(), uID, url, token); err != nil {
		fmt.Fprintf(w, "Error adding data source %s \n", err.Error())
	}

	CloseBrowserWindow(w, r)

	if err := PublishNotification(uCtx, uID); err != nil {
		log.Println("Error publishing notification msg (Service.Create)", err)
	}
}
