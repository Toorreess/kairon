package controllers

import (
	"kairon/cmd/api/infrastructure"
	"kairon/cmd/api/presenter"
	"kairon/domain/model"
	"kairon/usecases"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	HandleGet(c echo.Context) error
	HandlePost(c echo.Context, order model.Order) error
	HandleDelete(c echo.Context) error
	HandleList(c echo.Context) error
	HandlePay(c echo.Context) error
	HandleCancel(c echo.Context) error
}

type OrderHandlerImp struct {
	orderUsecase usecases.OrderUsecase
}

func NewOrderHandler(cu usecases.OrderUsecase) OrderHandler {
	return &OrderHandlerImp{
		orderUsecase: cu,
	}
}

func (h *OrderHandlerImp) HandleGet(c echo.Context) error {
	cm, err := h.orderUsecase.Read(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *OrderHandlerImp) HandlePost(c echo.Context, order model.Order) error {
	cm, err := h.orderUsecase.Create(order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *OrderHandlerImp) HandleDelete(c echo.Context) error {
	err := h.orderUsecase.Delete(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *OrderHandlerImp) HandleList(c echo.Context) error {
	queryString := c.QueryParam("q")
	offsetStr := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	qo := infrastructure.QueryOpts{
		QueryString: queryString,
		Offset:      offset,
		Limit:       limit,
	}

	results, err := h.orderUsecase.List(qo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, presenter.ListAPIResponse(results, qo.Offset, qo.Limit))
}

func (h *OrderHandlerImp) HandlePay(c echo.Context) error {
	cm, err := h.orderUsecase.Pay(c.Param("id"))
	if err != nil {
		log.Printf("Error paying order: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, cm)
}

func (h *OrderHandlerImp) HandleCancel(c echo.Context) error {
	cm, err := h.orderUsecase.Cancel(c.Param("id"))
	if err != nil {
		log.Printf("Error cancelling order: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, cm)
}
