package finance

type Currency string

const (
	RUB Currency = "RUB"
	EUR = "EUR"
	USD = "USD"
)

type Rate struct {
	Rate float64
}

type FinanceService interface {
	GetRate(from Currency, to Currency) (*Rate, error)
}