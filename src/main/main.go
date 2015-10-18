package main

import (
	"finance"
	"finance/yahoo"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/matcher"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type SSL struct {
	Certificate string
	Key         string
}

type Config struct {
	Telegram telegram.Config
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		return
	}

	cfg, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Panic(err)
	}

	var config Config
	if _, err := toml.Decode(string(cfg), &config); err != nil {
		log.Panic(err)
	}

	var financeService *yahoo.YahooFinanceService
	financeService, err = yahoo.New(yahoo.Config{Url: "https://query.yahooapis.com/v1/public/yql"})
	if err != nil {
		log.Panic(err)
	}

	app, err := telegram.New(config.Telegram)
	if err != nil {
		log.Panic("Could not init app : ", err)
	}

	app.Use(telegram.Condition{matcher.Equal{"/help"}, telegram.HandlerFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I do not have help yet! =("))
		ctrl.Stop()
	})})

	app.Use(telegram.Condition{matcher.RegExp{regexp.MustCompile(`/(usd|eur)`)}, telegram.HandlerFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		if match, ok := ctrl.Context().Value(telegram.MATCH).(matcher.Match); ok {
			rate, err := financeService.GetRate(finance.Currency(strings.ToUpper(match.Slugs[0].Value)), finance.RUB)
			if err != nil {
				ctrl.Throw(err)
				return
			}

			bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatFloat(rate.Rate, 'f', 2, 64)))
			ctrl.Stop()
		} else {
			ctrl.Throw(fmt.Errorf("Unexpected"))
		}
	})})

	app.UseFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I do not know this route yet"))
		ctrl.Stop()
	})

	app.UseErrFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {
		bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I've got error =("))
		ctrl.Stop()
	})

	log.Fatal(app.Listen())
}
