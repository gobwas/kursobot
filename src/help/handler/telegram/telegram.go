package telegram

import (
	"github.com/gobwas/telegram"
	"gopkg.in/telegram-bot-api.v2"
	"help"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, help.HelpContents)
	msg.ParseMode = telegram.ParseModeMarkdown
	bot.Send(msg)
	ctrl.Next()
}
