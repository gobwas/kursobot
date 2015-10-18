package controller
import (
	"finance"
	"net/http"
	"strconv"
)

type MainController struct {
	financeService finance.FinanceService
}

func New(financeService finance.FinanceService) *MainController {
	return &MainController{financeService: financeService}
}

func (self *MainController) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var rate *finance.Rate
	rate, err := self.financeService.GetRate(finance.USD, finance.RUB)
	if err != nil {
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
	rw.Write([]byte(strconv.FormatFloat(rate.Rate, 'E', 1, 32)))
}