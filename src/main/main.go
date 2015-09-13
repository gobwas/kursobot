package main

import (
	//	"finance/yahoo"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Syfaro/telegram-bot-api"
	"http"
	"io/ioutil"
	"log"
	"net/url"
)

type Config struct {
	Scheme string
	Host   string
	Token  string
	Debug  bool
}

func main() {
	//	var financeService *yahoo.YahooFinanceService
	//	financeService, err := yahoo.New(yahoo.Config{Url: "https://query.yahooapis.com/v1/public/yql"})
	//	if err != nil {
	//		log.Panic(err)
	//		return
	//	}
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
	if _, err := bot.SetWebhook(tgbotapi.NewWebhookWithCert(webHookUrl.String(), "server.crt")); err != nil {
		log.Panic("Could not set webhook", err)
		return
	}

	bot.ListenForWebhook()
	go http.ListenAndServeTLS("0.0.0.0:443", "server.crt", "server.key", nil)

	for update := range bot.Updates {
		log.Printf("%+v\n", update)
	}
}
