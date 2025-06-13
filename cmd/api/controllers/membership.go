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

type MembershipHandler interface {
	HandleGet(c echo.Context) error
	HandlePost(c echo.Context, membership model.Membership) error
	HandlePut(c echo.Context, membership model.Membership) error
	HandleDelete(c echo.Context) error
	HandleList(c echo.Context) error
}

type MembershipHandlerImp struct {
	membershipUsecase usecases.MembershipUsecase
}

func NewMembershipHandler(cu usecases.MembershipUsecase) MembershipHandler {
	return &MembershipHandlerImp{
		membershipUsecase: cu,
	}
}

func (h *MembershipHandlerImp) HandleGet(c echo.Context) error {
	cm, err := h.membershipUsecase.Read(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *MembershipHandlerImp) HandlePost(c echo.Context, membership model.Membership) error {
	cm, err := h.membershipUsecase.Create(membership)
	if err != nil {
		log.Printf("Error creating membership: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *MembershipHandlerImp) HandlePut(c echo.Context, membership model.Membership) error {
	changes, _ := c.Get("requestMap").(map[string]any)

	cm, err := h.membershipUsecase.Update(c.Param("id"), changes)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *MembershipHandlerImp) HandleDelete(c echo.Context) error {
	err := h.membershipUsecase.Delete(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *MembershipHandlerImp) HandleList(c echo.Context) error {
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

	results, err := h.membershipUsecase.List(qo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, presenter.ListAPIResponse(results, qo.Offset, qo.Limit))
}
