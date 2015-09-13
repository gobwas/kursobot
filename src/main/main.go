package main

import (
	"finance"
	"finance/yahoo"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type SSL struct {
	Certificate string
	Key         string
}

type Config struct {
	Scheme string
	Host   string
	Token  string
	Debug  bool
	SSL    SSL
}

func main() {
	var financeService *yahoo.YahooFinanceService
	financeService, err := yahoo.New(yahoo.Config{Url: "https://query.yahooapis.com/v1/public/yql"})
	if err != nil {
		log.Panic(err)
		return
	}

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		return
	}

	cfg, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Panic(err)
		return
	}

	var config Config
	if _, err := toml.Decode(string(cfg), &config); err != nil {
		log.Panic(err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic("Could not initialize bot: ", err)
		return
	}
	bot.Debug = config.Debug
	log.Println("Initialized bot")

	webHookUrl := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   config.Token,
	}
	if _, err := bot.SetWebhook(tgbotapi.WebhookConfig{URL: &webHookUrl, Certificate: config.SSL.Certificate}); err != nil {
		log.Panic("Could not set webhook", err)
		return
	}
	bot.ListenForWebhook()
	go http.ListenAndServeTLS(":443", config.SSL.Certificate, config.SSL.Key, nil)

	for update := range bot.Updates {
		if config.Debug {
			log.Printf("Got an update: %+v\n", update)
		}

		switch update.Message.Text {
		case "/help":
			bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I do not have help yet! =("))
		case "/usd":
			rate, err := financeService.GetRate(finance.USD, finance.RUB)
			if err != nil {
				log.Println("Got error: %v", err)
				bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I've got error =("))
			} else {
				bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatUint(rate.Rate, 10)))
			}
		case "/eur":
			rate, err := financeService.GetRate(finance.EUR, finance.RUB)
			if err != nil {
				log.Println("Got error: %v", err)
				bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I've got error =("))
			} else {
				bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatUint(rate.Rate, 10)))
			}
		}
	}
}
