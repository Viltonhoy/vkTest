package internal

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	clientchannel "vkTest/internal/client_channel"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"golang.org/x/mod/sumdb/storage"
)

type Bot struct {
	Logger  *zap.Logger
	TgBot   *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
	Store   storage.Storage
}

func NewBot(logger *zap.Logger, store *storage.Storage) (*Bot, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		logger.Error("failed to create a new BotAPI instance", zap.Error(err))
		return nil, err
	}

	bot.Debug = true

	logger.Debug("Authorized on ", zap.String("account", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	return &Bot{
		Logger:  logger,
		TgBot:   bot,
		Updates: updates,
		Store:   *store,
	}, err
}

func (b *Bot) BotWorker(cc *clientchannel.BotServer) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case update := <-b.Updates:
			if update.Message == nil && update.CallbackQuery == nil {
				continue
			}

			var chatID int64

			if update.Message != nil {
				continue
			}

			if update.CallbackQuery != nil {
				chatID = update.CallbackQuery.Message.Chat.ID
			}

			if f, ok := cc.Get(chatID); ok {
				f(update, b.TgBot)
				continue
			}

		}
	}
}

func (b *Bot) switcherCallback(update tgbotapi.Update, msg tgbotapi.MessageConfig, c *clientchannel.BotServer) {

}
