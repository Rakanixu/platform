package handler

import (
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func HandleMicrosoftLogin(w http.ResponseWriter, r *http.Request) {
	url := globals.NewMicrosoftOauthConfig().AuthCodeURL(globals.OauthStateString, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleMicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := globals.NewMicrosoftOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	url := fmt.Sprintf("onedrive://%s", token.Extra("user_id"))

	if err := SaveDatasource(url, token); err != nil {
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
