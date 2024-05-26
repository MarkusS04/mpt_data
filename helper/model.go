// Package helper provides general functions, like config handling
package helper

type loginResult struct {
	AccessToken       string `json:"accessToken"`
	ExpiresIn         uint   `json:"expiresIn"`
	AccessTokenMaxTTL uint   `json:"accessTokenMaxTTL"`
	TokenType         string `json:"tokenType"`
}

type secret struct {
	Secret struct {
		ID string `json:"id"`
		// _id           string `json:"_id"`
		Workspace     string `json:"workspace"`
		Environment   string `json:"environment"`
		Version       uint   `json:"version"`
		Type          string `json:"type"`
		SecretKey     string `json:"secretKey"`
		SecretValue   string `json:"secretValue"`
		SecretComment string `json:"secretComment"`
	} `json:"secret"`
}
