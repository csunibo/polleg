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
		Listen:               "0.0.0.0:3000",
		BaseURL:              "http://localhost:3000",
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
	if err := db.AutoMigrate(&api.Question{}, &api.Answer{}, &api.Vote{}); err != nil {
		slog.Error("AutoMigrate failed", "err", err)
		os.Exit(1)
	}

	authenticator := auth.NewAuthenticator(&auth.Config{
		BaseURL:      baseURL,
		ClientID:     config.OAuthClientID,
		ClientSecret: config.OAuthClientSecret,
		SigningKey:   []byte(config.OAuthSigningKey),
		Expiration:   config.OAuthSessionDuration,
	})

	// Routes
	mux := muxie.NewMux()
	mux.Use(util.NewCorsMiddleware(config.ClientURLs, true, mux))

	// authentication api
	mux.HandleFunc("/login", authenticator.LoginHandler)
	mux.HandleFunc("/login/callback", authenticator.CallbackHandler)

	mux.Use(authenticator.Middleware)
	mux.HandleFunc("/whoami", auth.WhoAmIHandler)

	// insert new answer
	mux.HandleFunc("/answers", api.PutAnswerHandler)
	// get answer
	mux.HandleFunc("/answers/:id", api.GetAnswerById)
	// get answers by question
	mux.HandleFunc("/answers/:id/vote", api.PostVote)

	// Same api
	mux.HandleFunc("/answers/by-question/:id", api.GetAnswersByQuestion)
	mux.HandleFunc("/question/:id/answers", api.GetAnswersByQuestion)

	// Same api
	mux.HandleFunc("/questions/by-document/:id", api.GetQuestionsByDoc)
	mux.HandleFunc("/documents/:id", api.GetQuestionsByDoc)

	// insert new doc and quesions
	mux.HandleFunc("/documents", api.PutDocument)
	// get all doc's questions
	mux.HandleFunc("/question/:id", api.GetQuestionsById)

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
