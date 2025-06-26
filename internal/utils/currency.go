package utils

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"regexp"

	"github.com/shopspring/decimal"
)

func IsValidCurrencyFormat(currency string) bool {
	regexp := regexp.MustCompile(`^[A-Z]{3}$`)
	return regexp.MatchString(currency)
}

// Gets current exchange rate for a currency
func GetExchangeRatesForPair(currencyA, currencyB, apiKey string, logger *zap.Logger) (decimal.Decimal, decimal.Decimal, error) {
	// Check if the currency was passed
	url := fmt.Sprintf("https://api.exchangeratesapi.io/v1/latest?access_key=%s", apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Failed to make currency exchange rate inquiry request", zap.Error(err))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to make request", zap.Error(err))
		return decimal.Zero, decimal.Zero, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to get exchange rates", zap.Int("status_code", resp.StatusCode))
		return decimal.Zero, decimal.Zero, err
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("Failed to decode response body", zap.Error(err))
		return decimal.Zero, decimal.Zero, err
	}

	strRateA := fmt.Sprintf("%f", result.Rates[currencyA])
	strRateB := fmt.Sprintf("%f", result.Rates[currencyB])

	decimalA, err := decimal.NewFromString(strRateA)
	if err != nil {
		logger.Error("Failed to parse currencyA rate", zap.Error(err))
		return decimal.Zero, decimal.Zero, err
	}

	decimalB, err := decimal.NewFromString(strRateB)
	if err != nil {
		logger.Error("Failed to parse currencyB rate", zap.Error(err))
		return decimal.Zero, decimal.Zero, err
	}

	return decimalA, decimalB, nil
}
