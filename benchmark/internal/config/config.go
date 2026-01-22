package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PlayerServiceURL      string
	MatchMakingServiceURL string
	HTTPTimeout           time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
}

func NewConfig() *Config {

	_ = godotenv.Load()

	playerServiceURL := os.Getenv("PLAYER_SERVICE_URL")
	matchMakingServiceURL := os.Getenv("MATCH_MAKING_SERVICE_URL")
	maxIdleConns := os.Getenv("MAX_IDLE_CONNS")
	maxIdleConnsPerHost := os.Getenv("MAX_IDLE_CONNS_PER_HOST")
	maxConnsPerHost := os.Getenv("MAX_CONNS_PER_HOST")

	if playerServiceURL == "" {
		playerServiceURL = "http://localhost:5108"
	}
	if matchMakingServiceURL == "" {
		matchMakingServiceURL = "http://localhost:8000"
	}
	if maxIdleConns == "" {
		maxIdleConns = "100"
	}
	if maxIdleConnsPerHost == "" {
		maxIdleConnsPerHost = "100"
	}
	if maxConnsPerHost == "" {
		maxConnsPerHost = "0"
	}
	maxIdleConnsInt, err := strconv.Atoi(maxIdleConns)
	if err != nil {
		panic(err)
	}
	maxIdleConnsPerHostInt, err := strconv.Atoi(maxIdleConnsPerHost)
	if err != nil {
		panic(err)
	}
	maxConnsPerHostInt, err := strconv.Atoi(maxConnsPerHost)
	if err != nil {
		panic(err)
	}

	return &Config{
		PlayerServiceURL:      playerServiceURL,
		MatchMakingServiceURL: matchMakingServiceURL,
		HTTPTimeout:           time.Second * 10,
		MaxIdleConns:          maxIdleConnsInt,
		MaxIdleConnsPerHost:   maxIdleConnsPerHostInt,
		MaxConnsPerHost:       maxConnsPerHostInt,
	}
}
