package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type accrualData struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

// updateOrderData is a job for calling accrual service for order info
func (s *Service) updateOrderData() {
	for s.active {
		order := <-s.OrderQueue
		var requestURL = s.addr + "/api/orders/" + strconv.Itoa(order)

		resp, err := http.Get(requestURL)
		if err != nil {
			s.logger.Error("Error when making request", zap.Error(err))
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			s.logger.Info("A lot of requests", zap.Any("headers", resp.Request))
			time.Sleep(time.Second * 2)
			s.OrderQueue <- order

		case http.StatusInternalServerError:
			s.logger.Error("Code 500 when getting accrual")
			s.active = false

		case http.StatusOK:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				s.logger.Error("Error reading body", zap.Error(err))
			}

			var respData accrualData
			if err := json.Unmarshal(body, &respData); err != nil {
				s.logger.Error("Error unmarshalling body", zap.Error(err))
			}

			err = s.storage.UpdateAccrual(order, respData.Status, decimal.NewFromFloat(respData.Accrual))
			if err != nil {
				s.logger.Error("Error storage call", zap.Error(err))
			}
		}
	}
}
