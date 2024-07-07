package main

import (
	"bot-adviser/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.Token)
	// TODO tgClient = telegram.New(token)
	// TODO fetcher = fetcher.New()
	// TODO processor = processor.New()

	// TODO consumer.Start(fetcher, processor)
}

func token() string {
	return fmt.Sprintf("Hello!")
}
