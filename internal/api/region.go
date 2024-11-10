package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/VanessaVallarini/campaign-producer-api/internal/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RegionService interface {
	Fetch(context.Context, uuid.UUID) (model.Region, error)
	Create(context.Context, model.RegionCreateRequest) error
	ListHistory(context.Context, uuid.UUID, model.ListingFilters) ([]model.RegionHistory, model.Paging, error)
}

type Region struct {
	service RegionService
}

func NewRegion(service RegionService) *Region {
	return &Region{
		service: service,
	}
}

func (s Region) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/region/:id", s.Fetch)
	v1.POST("/region", s.Create)
	v1.GET("/history/region/:id", s.List)
}

func (s Region) Fetch(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid region id"), "Invalid region id")
	}

	result, err := s.service.Fetch(c.Request().Context(), id)
	if err != nil {

		return ResponseError(c, err, "Failed to fetch region")
	}

	return c.JSON(http.StatusOK, result)
}

func (s Region) Create(c echo.Context) error {
	req := model.RegionCreateRequest{}

	if err := c.Bind(&req); err != nil {

		return ResponseError(c, err, "Failed to parse region input")
	}

	userEmail := util.GetHeaderFromContext(c, "x-user-email")
	req.CreatedBy = userEmail
	if userEmail == "" {
		return ResponseError(c, model.ErrInvalid.Wrap(nil, "Invalid request"), "Email is required")
	}

	err := s.service.Create(c.Request().Context(), req)
	if err != nil {

		return ResponseError(c, err, "Failed to create region")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (s Region) List(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid Region id"), "Invalid Region id")
	}

	filters := model.ListingFilters{}
	if err := c.Bind(&filters); err != nil {

		return ResponseError(c, err, "Failed to parse filters")
	}

	result, paging, err := s.service.ListHistory(c.Request().Context(), id, filters)
	if err != nil {

		return ResponseError(c, err, "Failed to list history Region record")
	}

	return c.JSON(http.StatusOK, model.PagedResults{Paging: paging, Data: result})
}
