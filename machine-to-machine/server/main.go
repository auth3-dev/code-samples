package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const serverAddr = ":8082"

func writeReply(w http.ResponseWriter, d interface{}) {
	data, err := json.Marshal(d)
	if err != nil {
		log.Printf("marshalling data: %s", err)
		// you should handle error here...
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func writeUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Not authorized", 401)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer")
		if len(splitToken) != 2 {
			writeUnauthorized(w)
			return
		}

		token = strings.TrimSpace(splitToken[1])
		serviceName := ""
		data, err := verifyToken(w, token)
		if err != nil {
			log.Printf("verifying token", err)
			writeUnauthorized(w)
			return
		}
		serviceName = data["sub"].(string)

		writeReply(w, map[string]interface{}{"greet": serviceName})
		return
	})

	log.Printf("Starting server @ %q. Press CTRL+C to exit.", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func verifyToken(w http.ResponseWriter, accessToken string) (map[string]interface{}, error) {
	log.Printf("token: %s", accessToken)

	headers := map[string][]string{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Accept":       []string{"application/json"},
	}

	body := []byte("token=" + accessToken)

	introspectTokenURL := "https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/oauth2/introspect"

	req, err := http.NewRequest("POST", introspectTokenURL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Printf("creating POST request to %s: %s", introspectTokenURL, err)
		return nil, err
	}

	req.Header = headers

	client := &http.Client{}
	response, err := client.Do(req)

	// response, err := http.Get(
	// 	"https://" + os.Getenv("AUTH3_PROJECT_ID") + ".as.auth3.dev/userinfo?access_token=" + accessToken,
	// )
	if err != nil {
		log.Printf("getting userinfo: %s", err)
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("reading response: %s", err.Error())
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Printf("decoding useringo response: %s", err.Error())
		return nil, err
	}

	return data, nil
}
