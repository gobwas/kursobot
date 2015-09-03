package yahoo

import (
	"finance"
	"net/url"
	"fmt"
	"net/http"
	"io/ioutil"
)

type YahooFinanceService struct {
	Url string
}

type Config struct {
	Url string
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
	parameters.Add("q", "select * from yahoo.finance.xchange where pair in (\"EURUSD\",\"EURRUB\", \"USDRUB\")")
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

	fmt.Println("Request!")
	fmt.Println(fmt.Sprintf("%s => %s", Url, string(contents)))

	return &finance.Rate{Rate:0}, nil
}