package clientchannel

import (
	// "vkTest/internal/keyboards"
	// tgutil "vkTest/internal/tg_util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type UpdateHandlerFunc func(update tgbotapi.Update, bot *tgbotapi.BotAPI)

type BotServer struct {
	Logger *zap.Logger
	State  map[int64]UpdateHandlerFunc
}

func New() *BotServer {
	bl := &BotServer{
		State: make(map[int64]UpdateHandlerFunc),
	}
	return bl
}

func (b *BotServer) Get(chadID int64) (f UpdateHandlerFunc, ok bool) {
	f, ok = b.State[chadID]
	return f, ok
}

func (b *BotServer) Add(chadID int64, f func(update tgbotapi.Update, bot *tgbotapi.BotAPI)) {
	b.State[chadID] = f
}

func (b *BotServer) DigitalSignature(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			if err != nil {
				b.Logger.Error("sending chattable error", zap.Error(err))
				return
			}
		}
		if update.CallbackQuery.Data == "Меню" {
			// logger.Printf(update.CallbackQuery.Data)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
			msg.ReplyMarkup = keyboards.StartKeyboard
			tgutil.SendBotMessage(msg, bot)
			delete(b.State, update.CallbackQuery.Message.Chat.ID)
		}

	}

}
