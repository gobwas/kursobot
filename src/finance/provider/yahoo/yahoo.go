package yahoo

import (
	"encoding/json"
	"finance"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Url                 string
	MaxIdleConnsPerHost int
	Timeout             int
	KeepAlive           int
	TLSHandshakeTimeout int
}

type YahooFinanceService struct {
	u      url.URL
	client *http.Client
}

func New(config Config) (*YahooFinanceService, error) {
	var u *url.URL
	u, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}

	parameters := url.Values{}
	parameters.Add("format", "json")
	parameters.Add("diagnostics", "false")
	parameters.Add("env", "store://datatables.org/alltableswithkeys")

	u.RawQuery = parameters.Encode()

	return &YahooFinanceService{*u, &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(config.Timeout) * time.Second,
				KeepAlive: time.Duration(config.KeepAlive) * time.Second,
			}).Dial,
			TLSHandshakeTimeout: time.Duration(config.TLSHandshakeTimeout) * time.Second,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		},
	}}, nil
}

func (self *YahooFinanceService) GetRate(ctx context.Context, from finance.Currency, to finance.Currency) (*finance.Rate, error) {
	// copy url
	u := self.u

	// prepare query
	q := u.Query()
	q.Add("q", fmt.Sprintf("select * from yahoo.finance.xchange where pair in (\"%s%s\")", from, to))
	u.RawQuery = q.Encode()

	var response *http.Response
	response, err := ctxhttp.Get(ctx, self.client, u.String())
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
