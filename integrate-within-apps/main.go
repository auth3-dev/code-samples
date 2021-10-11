package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

const serverAddr = ":8082"

func main() {
	config := &oauth2.Config{
		RedirectURL:  "http://localhost" + serverAddr + "/callback",
		ClientID:     os.Getenv("AUTH3_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH3_CLIENT_SECRET"),
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/auth",
			TokenURL: "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/token",
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token, _ := r.Cookie("token")
		username, _ := r.Cookie("username")

		if token == nil || username == nil {
			w.Write([]byte(
				"<h1>Welcome, anonymous. Please <a href=\"/login\">login</h1>."),
			)
			return
		}

		// TODO: get user
		w.Write([]byte(
			fmt.Sprintf("<h1>Welcome, %s.", username.Value)),
		)
		return
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		state := generateStateOauthCookie(w)
		url := config.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		oauthState, _ := r.Cookie("oauthstate")

		if r.URL.Query()["state"][0] != oauthState.Value {
			log.Println("invalid oauth state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		token, err := config.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			log.Printf("code exchange errored: %s", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// by here you have a valid Access Token
		log.Printf("Access token successfully obtained: %s", token.AccessToken)

		var expiration = time.Now().Add(60 * time.Minute)
		cookie := http.Cookie{Name: "token", Value: token.AccessToken, Expires: expiration}
		http.SetCookie(w, &cookie)

		setUserInfo(w, token.AccessToken)

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	})

	log.Printf("Starting server @ %q. Press CTRL+C to exit.", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func setUserInfo(w http.ResponseWriter, accessToken string) {
	response, err := http.Get(
		"https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/userinfo?access_token=" + accessToken,
	)
	if err != nil {
		log.Printf("getting userinfo: %s", err)
		return
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("reading response: %s", err.Error())
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Printf("decoding useringo response: %s", err.Error())
		return
	}

	var expiration = time.Now().Add(60 * time.Minute)
	log.Printf("setting username to %s", data["email"].(string))

	cookie := http.Cookie{Name: "username", Value: data["email"].(string), Expires: expiration}
	http.SetCookie(w, &cookie)

	return
}
