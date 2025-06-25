package handlers

import (
	"pgpockets/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
	logger *zap.Logger
}

func NewDashboardHandler(dashboardService services.DashboardService, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		logger: logger,
	}
}

func (d *DashboardHandler) GetExchangeRates(c *fiber.Ctx) error {
	rates, err := d.dashboardService.GetExchangeRates()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch exchange rates",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"rates": rates,
	})
}