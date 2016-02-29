package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram"
	"help"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, help.HelpContents))
	ctrl.Next()
}
