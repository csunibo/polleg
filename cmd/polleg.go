package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kataras/muxie"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slog"

	"github.com/csunibo/polleg/api"
	"github.com/csunibo/polleg/api/proposal"
	"github.com/csunibo/polleg/auth"
	"github.com/csunibo/polleg/util"
)

type Config struct {
	Listen     string   `toml:"listen"`
	BaseURL    string   `toml:"base_url"`
	ClientURLs []string `toml:"client_urls"`

	DbURI                string        `toml:"db_uri" required:"true"`
	OAuthClientID        string        `toml:"oauth_client_id" required:"true"`
	OAuthClientSecret    string        `toml:"oauth_client_secret" required:"true"`
	OAuthSigningKey      string        `toml:"oauth_signing_key" required:"true"`
	OAuthSessionDuration time.Duration `toml:"oauth_session_duration"`
}

var (
	// Default config values
	config = Config{
		Listen:               "0.0.0.0:3001",
		BaseURL:              "http://localhost:3001",
		OAuthSessionDuration: time.Hour * 12,
	}
)

func main() {
	err := loadConfig()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		slog.Error("failed to parse baseURL", "err", err)
		os.Exit(1)
	}

	err = util.ConnectDb(config.DbURI)
	if err != nil {
		slog.Error("failed to connect to db", "err", err)
		os.Exit(1)
	}
	db := util.GetDb()
	err = db.AutoMigrate(&proposal.Proposal{}, &api.Question{}, &api.Answer{}, &api.Vote{})
	if err != nil {
		slog.Error("AutoMigrate failed", "err", err)
		os.Exit(1)
	}

	authenticator := auth.NewAuthenticator(&auth.Config{
		BaseURL:    baseURL,
		SigningKey: []byte(config.OAuthSigningKey),
		Expiration: config.OAuthSessionDuration,
	})

	mux := muxie.NewMux()
	mux.Use(util.NewCorsMiddleware(config.ClientURLs, true, mux))

	// authentication-less read-only queries
	mux.HandleFunc("/documents/:id", api.GetDocumentHandler)
	mux.HandleFunc("/questions/:id", api.GetQuestionHandler)

	// authenticated queries
	mux.Use(authenticator.Middleware)
	// insert new answer
	mux.HandleFunc("/answers", api.PutAnswerHandler)
	// put up/down votes to an answer
	mux.HandleFunc("/answers/:id/vote", api.PostVote)
	// insert new doc and quesions
	mux.HandleFunc("/documents", api.PutDocumentHandler)

	mux.HandleFunc("/answers/:id", api.DelAnswerHandler)
	// proposal managers
	mux.HandleFunc("/proposals", proposal.ProposalHandler)
	mux.HandleFunc("/proposals/:id", proposal.ProposalByIdHandler)
	mux.HandleFunc("/proposals/document/:id", proposal.ProposalByDocumentHandler)

	slog.Info("listening at", "address", config.Listen)
	err = http.ListenAndServe(config.Listen, mux)
	if err != nil {
		slog.Error("failed to serve", "err", err)
	}
}

func loadConfig() (err error) {
	file, err := os.Open("config.toml")
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close config file: %w", err)
	}

	return nil
}
