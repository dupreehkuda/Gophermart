package processors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type accrualData struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

func (p processors) getAccrualData(order int) (string, int, error) {
	var respData accrualData
	var requestURL = p.sysAddr + "/" + string(rune(order))

	resp, err := http.Get(requestURL)
	if err != nil {
		p.logger.Error("Error when making request", zap.Error(err))
		return "", 0, err
	}

	if resp.StatusCode == 500 {
		p.logger.Error("Code 500 when getting accrual")
		return "", 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error("Error reading body", zap.Error(err))
		return "", 0, err
	}
	defer resp.Body.Close()

	p.logger.Debug("body", zap.Any("json", string(body)))

	if err := json.Unmarshal(body, &respData); err != nil {
		p.logger.Error("Error unmarshalling body", zap.Error(err))
		return "", 0, err
	}

	return respData.Status, respData.Accrual, nil
}
