package controllers

import (
	"kairon/cmd/api/presenter"
	"kairon/usecases"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ReportHandler interface {
	HandleGetFinancialReport(c echo.Context) error
}

type ReportHandlerImp struct {
	reportUsecase usecases.ReportUsecase
}

func NewReportHandler(cu usecases.ReportUsecase) ReportHandler {
	return &ReportHandlerImp{
		reportUsecase: cu,
	}
}

func (r *ReportHandlerImp) HandleGetFinancialReport(c echo.Context) error {
	startDateStr := c.QueryParam("startDate")
	endDateStr := c.QueryParam("endDate")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	cm, err := r.reportUsecase.GenerateFinancialReport(startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}
