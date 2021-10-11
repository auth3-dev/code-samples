package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"golang.org/x/oauth2"
)

const serverPort = ":8082"

func main() {
	if os.Getenv("AUTH3_CLIENT_ID") == "" {
		panic("Please set AUTH3_CLIENT_ID")
	}

	// This is optional, to verify the client credentials: if set, make sure to allow the client
	// to authenticate on the Token Endpoint via `post`.
	// if os.Getenv("AUTH3_CLIENT_SECRET") == "" {
	// 	panic("Please set AUTH3_CLIENT_SECRET")
	// }

	if os.Getenv("AUTH3_PROJECT_ID") == "" {
		panic("Please set AUTH3_PROJECT_ID")
	}

	conf := &oauth2.Config{
		ClientID: os.Getenv("AUTH3_CLIENT_ID"),
		// This is optional, to verify the client credentials: if set, make sure to allow the client
		// to authenticate on the Token Endpoint via `post`.
		// ClientSecret: os.Getenv("AUTH3_CLIENT_SECRET"),
		Scopes: []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/auth",
			TokenURL: "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/token",
		},
		// it's important that this is always set for public clients even if you just need to `Exchange()`
		// otherwise the library will try to validate credentials which might result in 400s if the client is not
		// configured to authenticate on the token endpoint.
		RedirectURL: "http://localhost" + serverPort + "/callback",
	}

	// initialize the code verifier
	var codeVerifier, err = cv.CreateCodeVerifierWithLength(96)
	if err != nil {
		panic(err)
	}

	// Create code_challenge with S256 method
	codeChallenge := codeVerifier.CodeChallengeS256()
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	authn := conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	fmt.Printf("Please visit the URL: %s\n", authn)

	ch := make(chan string, 1)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != state {
			log.Println("invalid oauth state")
			return
		}

		token, err := conf.Exchange(
			context.Background(),
			r.FormValue("code"),
			oauth2.SetAuthURLParam("code_verifier", codeVerifier.String()),
		)
		if err != nil {
			log.Printf("code exchange errored: %s", err.Error())
			return
		}

		ch <- token.AccessToken
		return
	})

	log.Printf("Starting server @ %q. Press CTRL+C to exit.", serverPort)
	go func() {
		log.Fatal(http.ListenAndServe(serverPort, nil))
	}()

	fmt.Printf("Token: %s\n", <-ch)
}
