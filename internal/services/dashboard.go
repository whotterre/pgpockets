package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pgpockets/internal/repositories"

	"go.uber.org/zap"
)

type DashboardService interface {
	GetExchangeRates() (map[string]float64, error)
}

type dashboardService struct {
	userRepo repositories.DashboardRepository
	logger   *zap.Logger
	apiKey string
}

func NewDashboardService(
	userRepo repositories.DashboardRepository,
	logger *zap.Logger,
	apiKey string,
) *dashboardService {
	return &dashboardService{
		userRepo: userRepo,
		logger:   logger,
		apiKey:  apiKey,
	}
}

func (s *dashboardService) GetExchangeRates() (map[string]float64, error) {
	// Make a get request to the exchange rates API
	url := fmt.Sprintf("https://api.exchangeratesapi.io/v1/latest?access_key=%s", s.apiKey)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		s.logger.Error("Failed to make request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Failed to get exchange rates", zap.Int("status_code", resp.StatusCode))
		return nil, err
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		s.logger.Error("Failed to decode response body", zap.Error(err))
		return nil, err
	}
	return result.Rates, nil
}
