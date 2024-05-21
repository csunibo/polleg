package auth

import (
	"net/url"
	"time"
)

type Config struct {
	BaseURL    *url.URL      // The base URL from where csunibo/upld is being served from
	SigningKey []byte        // The key to sign the JWTs with
	Expiration time.Duration // How long should user sessions last?
}

type Authenticator struct {
	baseURL    *url.URL
	expiration time.Duration
	signingKey []byte
}

type User struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
}

func NewAuthenticator(config *Config) *Authenticator {
	authenticator := Authenticator{
		baseURL:    config.BaseURL,
		signingKey: config.SigningKey,
		expiration: config.Expiration,
	}
	return &authenticator
}
