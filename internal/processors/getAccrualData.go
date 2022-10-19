package processors

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

func (p processors) getAccrualData(order int) (string, decimal.Decimal, error) {
	var respData accrualData
	var requestURL = p.sysAddr + "/api/orders/" + strconv.Itoa(order)
	p.logger.Debug("debugging call to accrual", zap.Any("url", requestURL))

	resp, err := http.Get(requestURL)
	if err != nil {
		p.logger.Error("Error when making request", zap.Error(err))
		return "", decimal.Zero, err
	}

	if resp.StatusCode == 500 {
		p.logger.Error("Code 500 when getting accrual")
		return "", decimal.Zero, err
	}

	if resp.StatusCode == 429 {
		p.logger.Debug("A lot of requests", zap.Any("headers", resp.Request))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error("Error reading body", zap.Error(err))
		return "", decimal.Zero, err
	}
	defer resp.Body.Close()

	p.logger.Debug("debugging call to accrual", zap.Any("body json", string(body)))

	if err := json.Unmarshal(body, &respData); err != nil {
		p.logger.Error("Error unmarshalling body", zap.Error(err))
		return "", decimal.Zero, err
	}

	return respData.Status, decimal.NewFromFloat(respData.Accrual), nil
}
