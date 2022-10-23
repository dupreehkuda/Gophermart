package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type accrualData struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (a actions) getAccrualData(order int) (string, decimal.Decimal, error) {
	var respData accrualData
	var requestURL = a.sysAddr + "/api/orders/" + strconv.Itoa(order)

	resp, err := http.Get(requestURL)
	if err != nil {
		a.logger.Error("Error when making request", zap.Error(err))
		return "", decimal.Zero, err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		a.logger.Error("Code 500 when getting accrual")
		return "", decimal.Zero, err
	}

	// todo: add repeating request logic
	if resp.StatusCode == http.StatusTooManyRequests {
		a.logger.Debug("A lot of requests", zap.Any("headers", resp.Request))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error("Error reading body", zap.Error(err))
		return "", decimal.Zero, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &respData); err != nil {
		a.logger.Error("Error unmarshalling body", zap.Error(err))
		return "", decimal.Zero, err
	}

	return respData.Status, decimal.NewFromFloat(respData.Accrual), nil
}
