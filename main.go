package main

import (
	"log"
	telegrambot "vkTest/internal/tg_bot"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	bot, err := telegrambot.NewBot(logger)
	if err != nil {
		logger.Fatal("failed to create a new BotAPI instance", zap.Error(err))
	}

}
