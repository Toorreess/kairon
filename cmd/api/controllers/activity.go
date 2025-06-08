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

type ActivityHandler interface {
	HandleGet(c echo.Context) error
	HandlePost(c echo.Context, activity model.Activity) error
	HandlePut(c echo.Context, activity model.Activity) error
	HandleDelete(c echo.Context) error
	HandleList(c echo.Context) error
}

type ActivityHandlerImp struct {
	activityUsecase usecases.ActivityUsecase
}

func NewActivityHandler(cu usecases.ActivityUsecase) ActivityHandler {
	return &ActivityHandlerImp{
		activityUsecase: cu,
	}
}

func (h *ActivityHandlerImp) HandleGet(c echo.Context) error {
	cm, err := h.activityUsecase.Read(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *ActivityHandlerImp) HandlePost(c echo.Context, activity model.Activity) error {
	cm, err := h.activityUsecase.Create(activity)
	if err != nil {
		log.Printf("Error creating activity: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *ActivityHandlerImp) HandlePut(c echo.Context, activity model.Activity) error {
	changes, _ := c.Get("requestMap").(map[string]any)

	cm, err := h.activityUsecase.Update(c.Param("id"), changes)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *ActivityHandlerImp) HandleDelete(c echo.Context) error {
	err := h.activityUsecase.Delete(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ActivityHandlerImp) HandleList(c echo.Context) error {
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

	results, err := h.activityUsecase.List(qo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, presenter.ListAPIResponse(results, qo.Offset, qo.Limit))
}
