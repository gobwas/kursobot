package telegram

import (
	"errors"
	"finance"
	"fmt"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/slugger"
	"github.com/kyokomi/emoji"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v2"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

var ErrRateExpectTwoArguments = errors.New("rate method should be called with 2 arguments")
var ErrRoubleExpectOneArgument = errors.New("rouble method should be called with single argument")

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

	if len(call.Query) > 0 {
		switch len(call.Query) {
		case 1:
			h.roubleInline(ctrl, bot, update, call)
		default:
			h.rateInline(ctrl, bot, update, call)
		}
	} else {
		switch call.Args[0] {
		case "rate":
			h.rate(ctrl, bot, update, call)
		default:
			h.rouble(ctrl, bot, update, call)
		}
	}

}

func (h *Handler) rate(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, call slugger.Call) {
	if len(call.Args) != 3 {
		ctrl.Throw(ErrRateExpectTwoArguments)
		return
	}
	if err := validateCurrency(call.Args[1]); err != nil {
		ctrl.Throw(err)
		return
	}
	if err := validateCurrency(call.Args[2]); err != nil {
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
func (h *Handler) rateInline(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, call slugger.Call) {
	if len(call.Query) != 2 {
		ctrl.Throw(ErrRateExpectTwoArguments)
		return
	}
	if err := validateCurrency(call.Query[0]); err != nil {
		ctrl.Throw(err)
		return
	}
	if err := validateCurrency(call.Query[1]); err != nil {
		ctrl.Throw(err)
		return
	}

	from, to := toCurrency(call.Query[0]), toCurrency(call.Query[1])
	rate, err := h.f.GetRate(ctrl.Context(), from, to)
	if err != nil {
		ctrl.Throw(err)
		return
	}

	bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results: []interface{}{
			inlineResultArticle(from, to, rate.Rate),
		},
	})

	ctrl.Next()
}
func (h *Handler) rouble(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, call slugger.Call) {
	if len(call.Args) != 1 {
		ctrl.Throw(ErrRoubleExpectOneArgument)
		return
	}
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
func (h *Handler) roubleInline(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, call slugger.Call) {
	if len(call.Query) != 1 {
		ctrl.Throw(ErrRoubleExpectOneArgument)
		return
	}
	if err := validateCurrency(call.Query[0]); err != nil {
		ctrl.Throw(err)
		return
	}

	from := toCurrency(call.Query[0])
	rate, err := h.f.GetRate(ctrl.Context(), from, finance.RUB)
	if err != nil {
		ctrl.Throw(err)
		return
	}

	bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results: []interface{}{
			inlineResultArticle(from, finance.RUB, rate.Rate),
		},
	})

	ctrl.Next()
}

func inlineResultArticle(from, to finance.Currency, rate float64) tgbotapi.InlineQueryResultArticle {
	text := fmt.Sprintf("1 %s = %.2f %s", from, rate, to)
	return tgbotapi.InlineQueryResultArticle{
		Type:        telegram.InlineQueryResultArticleType,
		MessageText: text,
		Title:       fmt.Sprintf("%s/%s", from, to),
		ID:          fmt.Sprintf("%X", rand.Int63()),
		Description: fmt.Sprintf("%s\n%s", text, emoji.Sprint(":sunglasses:")),
		ParseMode:   "Markdown",
		ThumbURL:    "http://vasi.net/uploads/posts/2012-07/1342465091_0.jpg",
	}
}

func validateCurrency(c string) error {
	if len(c) != 3 {
		return fmt.Errorf("invalid %q argument: currency should have 3 characters", c)
	}

	return nil
}

func toCurrency(s string) finance.Currency {
	return finance.Currency(strings.ToUpper(s[0:3]))
}
