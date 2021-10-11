package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	if os.Getenv("AUTH3_CLIENT_ID") == "" {
		panic("Please set AUTH3_CLIENT_ID")
	}

	if os.Getenv("AUTH3_PROJECT_ID") == "" {
		panic("Please set AUTH3_PROJECT_ID")
	}

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Generate your URI on https://console.auth3.dev/ using the Authorize URL builder.
	var authn string = "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/auth?" +
		"response_type=token&" +
		"client_id=" + os.Getenv("AUTH3_CLIENT_ID") + "&" +
		"state=" + state + "&" +
		"nonce=&" + state + "&" +
		"redirect_uri=https://" + os.Getenv("AUTH3_PROJECT_ID") + ".login.auth3.dev/success"

	fmt.Printf("Please open up your browser and visit the following URL:\n\n%s\n\n", authn)
	fmt.Printf("Once you are done paste the token here to continue: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}
