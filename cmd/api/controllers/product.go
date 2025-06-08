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

type ProductHandler interface {
	HandleGet(c echo.Context) error
	HandlePost(c echo.Context, product model.Product) error
	HandlePut(c echo.Context, product model.Product) error
	HandleDelete(c echo.Context) error
	HandleList(c echo.Context) error
}

type ProductHandlerImp struct {
	productUsecase usecases.ProductUsecase
}

func NewProductHandler(cu usecases.ProductUsecase) ProductHandler {
	return &ProductHandlerImp{
		productUsecase: cu,
	}
}

func (h *ProductHandlerImp) HandleGet(c echo.Context) error {
	cm, err := h.productUsecase.Read(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *ProductHandlerImp) HandlePost(c echo.Context, product model.Product) error {
	cm, err := h.productUsecase.Create(product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *ProductHandlerImp) HandlePut(c echo.Context, product model.Product) error {
	changes, _ := c.Get("requestMap").(map[string]any)

	cm, err := h.productUsecase.Update(c.Param("id"), changes)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *ProductHandlerImp) HandleDelete(c echo.Context) error {
	err := h.productUsecase.Delete(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ProductHandlerImp) HandleList(c echo.Context) error {
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

	results, err := h.productUsecase.List(qo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, presenter.ListAPIResponse(results, qo.Offset, qo.Limit))
}
