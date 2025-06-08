package controllers

import (
	"kairon/cmd/api/infrastructure"
	"kairon/cmd/api/presenter"
	"kairon/domain/model"
	"kairon/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	HandleGet(c echo.Context) error
	HandlePost(c echo.Context, user model.User) error
	HandlePut(c echo.Context, user model.User) error
	HandleDelete(c echo.Context) error
	HandleList(c echo.Context) error
}

type UserHandlerImp struct {
	userUsecase usecases.UserUsecase
}

func NewUserHandler(cu usecases.UserUsecase) UserHandler {
	return &UserHandlerImp{
		userUsecase: cu,
	}
}

func (h *UserHandlerImp) HandleGet(c echo.Context) error {
	cm, err := h.userUsecase.Read(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	c.JSON(http.StatusOK, cm)
	return nil
}

func (h *UserHandlerImp) HandlePost(c echo.Context, user model.User) error {
	cm, err := h.userUsecase.Create(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	c.JSON(http.StatusOK, cm)
	return nil
}

func (h *UserHandlerImp) HandlePut(c echo.Context, user model.User) error {
	changes, _ := c.Get("requestMap").(map[string]any)

	cm, err := h.userUsecase.Update(c.Param("id"), changes)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	c.JSON(http.StatusOK, cm)
	return nil
}

func (h *UserHandlerImp) HandleDelete(c echo.Context) error {
	err := h.userUsecase.Delete(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandlerImp) HandleList(c echo.Context) error {
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

	results, err := h.userUsecase.List(qo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	c.JSON(http.StatusOK, presenter.ListAPIResponse(results, qo.Offset, qo.Limit))
	return nil
}
