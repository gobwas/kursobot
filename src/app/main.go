package main

import (
	financeHandler "finance/handler/telegram"
	"finance/provider/yahoo"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gobwas/telegram"
	"github.com/gobwas/telegram/handler/canceler"
	"github.com/gobwas/telegram/handler/condition"
	"github.com/gobwas/telegram/handler/condition/matcher"
	"github.com/gobwas/telegram/handler/slugger"
	"github.com/kyokomi/emoji"
	"gopkg.in/telegram-bot-api.v2"
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
			log.Printf(
				"[%d][%s] incoming message: text:%q; query:%q",
				update.UpdateID, update.Message.From.String(),
				update.Message.Text, update.InlineQuery.Query,
			)
			ctrl.Next()
		}),
	)

	// logic
	app.UseOn("/help", helpHandler.New())

	fh := financeHandler.New(f)
	app.Use(condition.Condition{
		matcher.RegExp{
			Source:  matcher.SourceText,
			Pattern: regexp.MustCompile(`^\/([a-z]{3}$|rate.*$)`),
		},
		fh,
	})
	app.Use(condition.Condition{
		matcher.RegExp{
			Source:  matcher.SourceQuery,
			Pattern: regexp.MustCompile(`^([a-z]{3}|[a-z]{3} [a-z]{3})`),
		},
		fh,
	})

	// error handler
	app.UseErrFunc(func(ctrl *telegram.Control, bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {
		if update.Message.Text != "" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Sprintf(":space_invader:error: %s", err)))
		}
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