package finance

type Currency string

const (
	RUB Currency = "RUB"
	EUR          = "EUR"
	USD          = "USD"
)

type Rate struct {
	Id   string
	Name string
	Rate float64
	Date string
	Ask  float64
	Bid  float64
}

type FinanceService interface {
	GetRate(from Currency, to Currency) (*Rate, error)
}
