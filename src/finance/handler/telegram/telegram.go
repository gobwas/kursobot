package telegram

import (
	"errors"
	"finance"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/slugger"
	"golang.org/x/net/context"
	"reflect"
	"strconv"
	"strings"
)

type FinanceService interface {
	GetRate(context.Context, finance.Currency, finance.Currency) (*finance.Rate, error)
}

type Handler struct {
	f FinanceService
}

func New(f FinanceService) *Handler {
	return &Handler{f}
}

const prefix = "/"

func (h *Handler) Serve(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var call slugger.Call
	if v := ctrl.Context().Value(reflect.TypeOf(call)); v == nil {
		ctrl.Throw(errors.New("slugger should be used before"))
		return
	} else {
		call = v.(slugger.Call)
	}

	if len(call.Args) < 1 {
		ctrl.Throw(errors.New("wrong call"))
	}

	switch call.Args[0] {
	case "rate":
		h.rate(ctrl, bot, update, call)
	default:
		h.rouble(ctrl, bot, update, call)
	}
}

func validateCurrency(c string) error {
	if len(c) != 3 {
		return fmt.Errorf("wrong %q argument: currency should have 3 characters", c)
	}

	return nil
}

func validateRate(call slugger.Call) error {
	if len(call.Args) != 3 {
		return errors.New("rate method should be called with 2 arguments")
	}

	for _, c := range call.Args[1:] {
		if e := validateCurrency(c); e != nil {
			return e
		}
	}

	return nil
}

func toCurrency(s string) finance.Currency {
	return finance.Currency(strings.ToUpper(s[0:3]))
}

func (h *Handler) rate(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, call slugger.Call) {
	if err := validateRate(call); err != nil {
		ctrl.Throw(err)
		return
	}

	rate, err := h.f.GetRate(ctrl.Context(), toCurrency(call.Args[1]), toCurrency(call.Args[2]))
	if err != nil {
		ctrl.Throw(err)
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatFloat(rate.Rate, 'f', 2, 64)))
	ctrl.Next()
}

func (h *Handler) rouble(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, call slugger.Call) {
	if err := validateCurrency(call.Args[0]); err != nil {
		ctrl.Throw(err)
		return
	}

	rate, err := h.f.GetRate(ctrl.Context(), toCurrency(call.Args[0]), finance.RUB)
	if err != nil {
		ctrl.Throw(err)
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatFloat(rate.Rate, 'f', 2, 64)))
	ctrl.Next()
}
