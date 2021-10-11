package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var clientID = os.Getenv("AUTH3_CLIENT_ID")
var projectID = os.Getenv("AUTH3_PROJECT_ID")

// This URI is also broadcasted by the metadata endpoint to support OIDC dynamic configuration.
var authServerURI = "https://" + projectID + ".as.auth3.dev"
var authEndpoint = authServerURI + "/device/auth"
var tokenEndpoint = authServerURI + "/oauth2/token"

type AuthResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

func main() {
	data := url.Values{}
	data.Set("client_id", clientID)

	req, err := http.NewRequest("POST", authEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(resp)
	}

	response := AuthResponse{}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Please visit the following URL:\n\n%s\n\nAnd type the following code when requested:\n\t\t%s\n",
		response.VerificationURI,
		response.UserCode,
	)

	ticker := time.NewTicker(time.Duration(response.Interval) * time.Second)
	quit := make(chan struct{})
	token := ""
	run := true

	for run == true {
		select {
		case <-ticker.C:
			token = getToken(response.DeviceCode, response.Interval)
			if token != "" {
				close(quit)
			} else {
				fmt.Printf(".")
			}
		case <-quit:
			ticker.Stop()
			run = false
		}
	}

	fmt.Printf("\nToken: %s\n", token)
	// perform protected calls
}

func getToken(deviceCode string, intervalInSeconds int) string {
	client := &http.Client{}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		panic(resp)
	}

	response := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	if err, ok := response["error"].(string); ok {
		// authorization_pending is the normal error returned when the user is still completing the auth flow
		if err != "" && err != "authorization_pending" {
			panic(err)
		}
	}

	tok := ""
	if t, ok := response["access_token"].(string); ok {
		tok = t
	}

	return tok
}
