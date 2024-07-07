package main

import (
	"bot-adviser/clients/telegram"
	"bot-adviser/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	tgClient := telegram.New(cfg.TelegramBotHost, cfg.Token)
	// TODO fetcher = fetcher.New()
	// TODO processor = processor.New()

	// TODO consumer.Start(fetcher, processor)
}

func token() string {
	return fmt.Sprintf("Hello!")
}
