package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Queries         []string
	IntervalSeconds int
	TelegramToken   string
	ChatID          int64
	Proxy           string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Ignore error, it might not exist if using pure env vars

	queriesStr := os.Getenv("DEPOP_QUERIES")
	if queriesStr == "" {
		return nil, fmt.Errorf("DEPOP_QUERIES is required (comma-separated)")
	}
	queries := strings.Split(queriesStr, ",")
	for i := range queries {
		queries[i] = strings.TrimSpace(queries[i])
	}

	intervalStr := os.Getenv("DEPOP_INTERVAL")
	interval := 30
	if intervalStr != "" {
		if val, err := strconv.Atoi(intervalStr); err == nil {
			interval = val
		}
	}

	tgToken := os.Getenv("TG_TOKEN")
	if tgToken == "" {
		return nil, fmt.Errorf("TG_TOKEN is required")
	}

	chatIDStr := os.Getenv("TG_CHAT_ID")
	if chatIDStr == "" {
		return nil, fmt.Errorf("TG_CHAT_ID is required")
	}
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TG_CHAT_ID: %v", err)
	}

	proxy := os.Getenv("DEPOP_PROXY")

	return &Config{
		Queries:         queries,
		IntervalSeconds: interval,
		TelegramToken:   tgToken,
		ChatID:          chatID,
		Proxy:           proxy,
	}, nil
}
