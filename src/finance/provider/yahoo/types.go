//go:generate easyjson -all
package yahoo

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
