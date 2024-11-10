package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/VanessaVallarini/campaign-producer-api/internal/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CampaignService interface {
	Fetch(context.Context, uuid.UUID) (model.Campaign, error)
	Create(context.Context, model.CampaignCreateRequest) error
	ListHistory(context.Context, uuid.UUID, model.ListingFilters) ([]model.CampaignHistory, model.Paging, error)
}

type Campaign struct {
	service CampaignService
}

func NewCampaign(service CampaignService) *Campaign {
	return &Campaign{
		service: service,
	}
}

func (s Campaign) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/campaign/:id", s.Fetch)
	v1.POST("/campaign", s.Create)
	v1.GET("/history/campaign/:id", s.List)
}

func (s Campaign) Fetch(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid campaign id"), "Invalid campaign id")
	}

	result, err := s.service.Fetch(c.Request().Context(), id)
	if err != nil {

		return ResponseError(c, err, "Failed to fetch campaign")
	}

	return c.JSON(http.StatusOK, result)
}

func (s Campaign) Create(c echo.Context) error {
	req := model.CampaignCreateRequest{}

	if err := c.Bind(&req); err != nil {

		return ResponseError(c, err, "Failed to parse campaign input")
	}

	userEmail := util.GetHeaderFromContext(c, "x-user-email")
	req.CreatedBy = userEmail
	if userEmail == "" {
		return ResponseError(c, model.ErrInvalid.Wrap(nil, "Invalid request"), "Email is required")
	}

	err := s.service.Create(c.Request().Context(), req)
	if err != nil {

		return ResponseError(c, err, "Failed to create campaign")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (s Campaign) List(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid campaign id"), "Invalid campaign id")
	}

	filters := model.ListingFilters{}
	if err := c.Bind(&filters); err != nil {

		return ResponseError(c, err, "Failed to parse filters")
	}

	result, paging, err := s.service.ListHistory(c.Request().Context(), id, filters)
	if err != nil {

		return ResponseError(c, err, "Failed to list history campaign record")
	}

	return c.JSON(http.StatusOK, model.PagedResults{Paging: paging, Data: result})
}
