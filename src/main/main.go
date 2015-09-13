package main

import (
	//	"finance/yahoo"
	"crypto/tls"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	if _, err := bot.SetWebhook(tgbotapi.WebhookConfig{URL: &webHookUrl, Certificate: config.SSL.Certificate}); err != nil {
		log.Panic("Could not set webhook", err)
		return
	}

	bot.ListenForWebhook()

	go func() {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		server := &http.Server{Addr: ":443", TLSConfig: tlsConfig}
		server.ListenAndServeTLS(config.SSL.Certificate, config.SSL.Key)
	}()
	//	go http.ListenAndServeTLS(":443", config.SSL.Certificate, config.SSL.Key, nil)

	for update := range bot.Updates {
		log.Printf("%+v\n", update)
	}
}
