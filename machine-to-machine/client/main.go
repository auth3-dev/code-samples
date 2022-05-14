package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("AUTH3_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH3_CLIENT_SECRET"),
		Scopes:       []string{}, // add here required scopes, make sure you add them to whitelisted scopes on the client config too
		TokenURL:     "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/token",
	}

	client := config.Client(context.TODO())
	response, err := client.Get("http://localhost:8082")
	if err != nil {
		log.Printf("performing request: %s", err)
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
		log.Printf("decoding userinfo response: %s", err.Error())
		return
	}

	log.Printf("Received: %s", data["greet"])
}
