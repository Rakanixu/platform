package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

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

func HandleDropboxLogin(w http.ResponseWriter, r *http.Request) {
	url := globals.NewDropboxOauthConfig().AuthCodeURL(r.URL.Query().Get("user"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleDropboxCallback(w http.ResponseWriter, r *http.Request) {
	var da *DropboxAccount

	state := r.FormValue("state")
	if len(state) == 0 {
		fmt.Printf("invalid oauth state, got '%s'\n", state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := globals.NewDropboxOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	c := &http.Client{}
	b := []byte(`{"account_id":"` + token.Extra("account_id").(string) + `"}`)

	req, err := http.NewRequest("POST", globals.DropboxAccountEndpoint, bytes.NewBuffer(b))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		log.Printf("Getting user account failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer rsp.Body.Close()

	contents, err := ioutil.ReadAll(rsp.Body)
	if err := json.Unmarshal(contents, &da); err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
	}
	url := fmt.Sprintf("dropbox://%s", da.Email)

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
