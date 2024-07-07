package main

import (
	"bot-adviser/clients/telegram"
	event_consumer "bot-adviser/consumer/event-consumer"
	tgEvent "bot-adviser/events/telegram"
	"bot-adviser/internal/config"
	"bot-adviser/storage/files"
	"fmt"
	"log"
)

func main() {
	cfg := config.MustLoad()
	tgClient := telegram.New(cfg.TelegramBotHost, cfg.Token)
	eventsProcessor := tgEvent.New(
		tgClient,
		files.New(cfg.StoragePath))

	fmt.Println(eventsProcessor)
	log.Print("Service has been started...")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}
