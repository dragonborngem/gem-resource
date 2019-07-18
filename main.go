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

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var OauthConfig *oauth2.Config

// Endpoint is Google's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:   "http://localhost:9094/auth",
	TokenURL:  "http://localhost:9094/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

type Claims struct {
	jwt.StandardClaims
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	RedirectURL := os.Getenv("REDIRECT_URL")
	
	if os.Getenv("MODE") != "" {
		RedirectURL = os.Getenv("REDIRECT_URLL")
		Endpoint = oauth2.Endpoint{
			AuthURL:   "https://gem-auth.herokuapp.com/auth",
			TokenURL:  "https://gem-auth.herokuapp.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		}
	}
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


//ApiLayerFromCookie Self validation to check token validate
func ApiLayerFromCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
	if err == nil {
		loginSession := models.LoginSession{}
		db:=models.OpenDB()
		err := db.Debug().Where("session_value = ?",c.Value).First(&loginSession).Error
		if err != nil{

			v.RespondUnauthorized(w,"wrong session")
			return
		}
		accessToken := loginSession.AccessToken

		// Kiểm tra xem có tồn tại token không
		if accessToken == "" {
			v.Respond(w, v.Message(false, "An authorization header is required"))
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("12345678"), nil
		})
		if err != nil {
			if (err.Error() == "Token is expired"){
				signingMethod := jwt.GetSigningMethod("HS256")
				fmt.Println(token)
				fmt.Println("=========================================")
				expirationTime := time.Now().Add(24 * time.Hour)
				claims := &Claims{
					StandardClaims: jwt.StandardClaims{
						// In JWT, the expiry time is expressed as unix milliseconds
						ExpiresAt: expirationTime.Unix(),
						Audience: os.Getenv("CLIENT_ID"),
						Subject: loginSession.UserID,
					},
				}
				newToken := jwt.NewWithClaims(signingMethod, claims)
				tokenString, _ := newToken.SignedString([]byte("12345678"))
				fmt.Println(tokenString)
				fmt.Println(token)
				fmt.Println("===================")
				newToken, verr := jwt.Parse(tokenString, func(newToken *jwt.Token) (interface{}, error) {
					if _, ok := newToken.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("12345678"), nil
				})
				if verr != nil{
					v.Respond(w, v.Message(false, verr.Error()))
					return
				}
				}else{
					v.Respond(w, v.Message(false, err.Error()))
					return
				}
			}
		next.ServeHTTP(w, r)
	}else{
		v.RespondUnauthorized(w,"no session")
	}	
	})
}


func protected(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Hola")
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
	url := OauthConfig.AuthCodeURL("gem-secrect-code")
	fmt.Println(url)	
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	})
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("call?")
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

	finalHandler := http.HandlerFunc(protected)
	http.Handle("/protected-resource", ApiLayerFromCookie(finalHandler))
	if os.Getenv("MODE") != "" {
		log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+"9096", nil))
	}
}
