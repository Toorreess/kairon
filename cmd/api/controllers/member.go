package controllers

import (
	"kairon/cmd/api/infrastructure"
	"kairon/cmd/api/presenter"
	"kairon/config"
	"kairon/domain/model"
	"kairon/usecases"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MemberHandler interface {
	HandleGet(c echo.Context) error
	HandlePost(c echo.Context, member model.Member) error
	HandlePut(c echo.Context, member model.Member) error
	HandleDelete(c echo.Context) error
	HandleList(c echo.Context) error
	HandleSendEmail(c echo.Context) error
}

type MemberHandlerImp struct {
	memberUsecase usecases.MemberUsecase
}

func NewMemberHandler(cu usecases.MemberUsecase) MemberHandler {
	return &MemberHandlerImp{
		memberUsecase: cu,
	}
}

func (h *MemberHandlerImp) HandleGet(c echo.Context) error {
	cm, err := h.memberUsecase.Read(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *MemberHandlerImp) HandlePost(c echo.Context, member model.Member) error {
	cm, err := h.memberUsecase.Create(member)
	if err != nil {
		log.Printf("Error creating member: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *MemberHandlerImp) HandlePut(c echo.Context, member model.Member) error {
	changes, _ := c.Get("requestMap").(map[string]any)

	cm, err := h.memberUsecase.Update(c.Param("id"), changes)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	return c.JSON(http.StatusOK, cm)
}

func (h *MemberHandlerImp) HandleDelete(c echo.Context) error {
	err := h.memberUsecase.Delete(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *MemberHandlerImp) HandleList(c echo.Context) error {
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

	results, err := h.memberUsecase.List(qo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, presenter.ListAPIResponse(results, qo.Offset, qo.Limit))
}

func (h *MemberHandlerImp) HandleSendEmail(c echo.Context) error {
	memberID := c.Param("id")
	cm, err := h.memberUsecase.Read(memberID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, presenter.APIResponse(http.StatusNotFound, err.Error()))
	}

	var req model.EmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, presenter.APIResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	if err := h.memberUsecase.SendEmail(
		config.C.Smtp.Host, config.C.Smtp.Email, config.C.Smtp.Password, config.C.Smtp.Port,
		cm.Email, req.Subject, req.Body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, presenter.APIResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, presenter.APIResponse(http.StatusOK, "Email sent succesfully"))
}
