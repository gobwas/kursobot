package main

import (
	financeHandler "finance/handler/telegram"
	"finance/provider/yahoo"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/canceler"
	"github.com/gobwas/telegram/handler/condition"
	"github.com/gobwas/telegram/handler/condition/matcher"
	"github.com/gobwas/telegram/handler/slugger"
	helpHandler "help/handler/telegram"
	"io/ioutil"
	"log"
	"regexp"
)

type SSL struct {
	Certificate string
	Key         string
}

type Config struct {
	Telegram telegram.Config
	Yahoo    yahoo.Config
	Canceler duration
}

func main() {
	configPath := flag.String("c", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		return
	}

	// read and decode configuration
	var config Config
	cfg, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Panic(err)
	}
	if _, err := toml.Decode(string(cfg), &config); err != nil {
		log.Panic(err)
	}

	log.Println("config:", fmt.Sprintf("%#v", config))

	// initialize yahoo finance
	var f *yahoo.YahooFinanceService
	f, err = yahoo.New(config.Yahoo)
	if err != nil {
		log.Panic(err)
	}

	// initialize telegram framework
	app, err := telegram.New(config.Telegram)
	if err != nil {
		log.Panic("could not init app : ", err)
	}

	// helper handlers
	app.Use(
		&slugger.Slugger{},
		&canceler.Canceler{config.Canceler.Duration},

		telegram.HandlerFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
			log.Println("got message:", fmt.Sprintf(update.Message.Text))
			ctrl.Next()
		}),
	)

	// logic
	app.UseOn("/help", helpHandler.New())
	app.Use(condition.Condition{
		matcher.RegExp{regexp.MustCompile(`^\/([a-z]{3}$|rate.*$)`)},
		financeHandler.New(f),
	})
	// todo use gopkg.in/telegram-bot-api.v2
	//	app.Use(telegram.HandlerFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//		log.Println("INLINE", update.InlineQuery.Query)
	//		if update.InlineQuery.Query == "usd" {
	//			bot.AnswerInlineQuery(tgbotapi.InlineConfig{
	//				InlineQueryID: update.InlineQuery.ID,
	//				Results: []tgbotapi.InlineQueryResult{
	//					tgbotapi.InlineQueryResult{
	//						Type: tgboty,
	//						ID:   "ddd",
	//					},
	//				},
	//			})
	//		}
	//		ctrl.Next()
	//	}))

	// error handler
	app.UseErrFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("sad, but bot got error: %s", err)))
		log.Println("got error", err)
		ctrl.Stop()
	})

	app.Use(telegram.HandlerFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		log.Println("message processing complete")
		ctrl.Next()
	}))

	// start listen
	log.Fatal(app.Listen())
}
