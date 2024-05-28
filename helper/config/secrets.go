// Package config provides configuration and secret management
package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const infisicalURL string = "https://app.infisical.com/api/"

var loginRes loginResult

// loginToSecretProvider will perform the request to get acces token
func loginToSecretProvider() {
	if !Config.SECRETS.Use {
		fmt.Println("Secret managment platform disabled. Will not be used.")
		return
	}

	data := url.Values{
		"clientSecret": {Config.SECRETS.ClientSecret},
		"clientId":     {Config.SECRETS.ClientID},
	}

	req, err := http.NewRequest(http.MethodPost, infisicalURL+"v1/auth/universal-auth/login", strings.NewReader(data.Encode()))
	if err != nil {
		panic(fmt.Sprintf("could not get login request: %v", err))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		panic(fmt.Sprintf("could not make login request: %v", err))
	}
	defer resp.Body.Close()

	// Read the response body
	if err = json.NewDecoder(resp.Body).Decode(&loginRes); err != nil {
		panic(fmt.Sprintf("could not get login data: %v", err))
	}
}

func getKeyFromSecretProvider(key string) string {
	if !Config.SECRETS.Use {
		fmt.Println("Secret managment platform disabled. Cannot load any key.")
		return ""
	}

	url := fmt.Sprintf(
		"%sv3/secrets/raw/%s?environment=%s?workspaceId=?",
		infisicalURL,
		key,
		Config.SECRETS.Environment,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(fmt.Sprintf("could not get http request: %v", err))
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", loginRes.TokenType, loginRes.AccessToken))

	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		panic(fmt.Sprintf("could not make secret loading request: %v", err))
	}
	defer resp.Body.Close()

	// Read the response body
	var secret secret
	if err = json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		panic(fmt.Sprintf("could not get secret data: %v", err))
	}
	return secret.Secret.SecretValue
}
