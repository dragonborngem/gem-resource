package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	models "gem-resource/app/models"
	v "gem-resource/app/utils/view"

	"github.com/satori/go.uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
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

func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err == nil {
		fmt.Fprintf(w,"user logged \n")
		fmt.Fprintf(w,c.Value+"\n")
		loginSession := models.LoginSession{}
		db:=models.OpenDB()
		err := db.Debug().Where("session_value = ?",c.Value).First(&loginSession).Error
		if err != nil{
			v.RespondUnauthorized(w,"wrong session")
			return
		}
		//v.RespondSuccess(w," ")
		return
	}
		// config := jwt.Config{}
		// client := config.Client(context.Background())
		url := OauthConfig.AuthCodeURL("gem-secrect-code")
		fmt.Println(url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

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
		//t := template.Must(template.ParseFiles("index.html"))

		// Create a new random session token
		sessionToken := uuid.NewV4().String()
		NewSession := models.LoginSession{
			AccessToken:  AccessToken.AccessToken,
			TokenType:    AccessToken.TokenType,
			RefreshToken: AccessToken.RefreshToken,
			Expiry:       AccessToken.Expiry,
			SessionID:    "session_token",
			SessionValue: sessionToken,
			UserID:       "1",
		}
		NewSession.WriteToDB()

		//w.Header().Set("Authorization", "Bearer "+AccessToken.AccessToken)
		addCookie(w, "session_token", sessionToken)
		v.Respond(w, data)

		// _ = t.Execute(w, "")
		return
		// state := r.FormValue("state")
		// url := "http://localhost:8000/token?client_id=" +
		// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("index.html"))
		_ = t.Execute(w, "")
	})
	http.HandleFunc("protected-resource", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("index.html"))
		_ = t.Execute(w, "")
	})
	log.Fatal(http.ListenAndServe(":9096", nil))
}
