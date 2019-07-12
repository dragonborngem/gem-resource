package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	v "gem-resource/app/utils/view"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var OauthConfig *oauth2.Config

// Endpoint is Google's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:   "http://localhost:8000/auth",
	TokenURL:  "http://localhost:8000/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	RedirectURL := os.Getenv("REDIRECT_URL")
	ClientID := os.Getenv("CLIENT_ID")
	ClientSecret := os.Getenv("CLIENT_SECRET")
	SCOPES := os.Getenv("SCOPES")
	OauthConfig = &oauth2.Config{
		RedirectURL:  RedirectURL,
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       []string{SCOPES},
		Endpoint:     Endpoint,
	}
}

func main() {
	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		htmlIndex := OauthConfig.AuthCodeURL("gem-secrect-code", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		url := OauthConfig.AuthCodeURL("gem-secrect-code")
		fmt.Println(url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		htmlIndex = `
		Hello, I'm protected `
		fmt.Fprintf(w, htmlIndex)
	})
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		AuthCode := r.FormValue("code")
		AccessToken, err := OauthConfig.Exchange(context.Background(), AuthCode)
		fmt.Printf("%+v\n", AccessToken)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer r.Body.Close()
		data := v.Message(true, "success")
		data["token"] = AccessToken
		v.Respond(w, data)
		return
		// state := r.FormValue("state")
		// url := "http://localhost:8000/token?client_id=" +
		// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
	log.Fatal(http.ListenAndServe(":9096", nil))
}
