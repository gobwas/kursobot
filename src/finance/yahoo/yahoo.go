package yahoo

import (
	"encoding/json"
	"finance"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type YahooFinanceService struct {
	Url string
}

type Config struct {
	Url string
}

type Response struct {
	Query Result `json:"query"`
}

type Result struct {
	Count   int        `json:"count"`
	Created string     `json:"created"`
	Lang    string     `json:"lang"`
	Results ResultRate `json:"results"`
}

type ResultRate struct {
	Rate YahooRate `json:"rate"`
}

type YahooRate struct {
	Id   string `json:"id"`
	Name string
	Rate string
	Date string
	Time string
	Ask  string
	Bid  string
}

func New(config Config) (*YahooFinanceService, error) {
	var Url *url.URL
	Url, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}

	parameters := url.Values{}
	parameters.Add("format", "json")
	parameters.Add("diagnostics", "false")
	parameters.Add("env", "store://datatables.org/alltableswithkeys")

	Url.RawQuery = parameters.Encode()

	return &YahooFinanceService{Url: Url.String()}, nil
}

func (self *YahooFinanceService) GetRate(from finance.Currency, to finance.Currency) (*finance.Rate, error) {
	// prepare params
	parameters := url.Values{}
	parameters.Add("q", fmt.Sprintf("select * from yahoo.finance.xchange where pair in (\"%s%s\")", from, to))
	Url := self.Url + "&" + parameters.Encode()

	var response *http.Response
	response, err := http.Get(Url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var r Response
	if err := json.Unmarshal(contents, &r); err != nil {
		return nil, err
	}

	rate, err := strconv.ParseFloat(r.Query.Results.Rate.Rate, 64)
	if err != nil {
		return nil, err
	}

	return &finance.Rate{Rate: rate}, nil
}
