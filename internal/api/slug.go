package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/VanessaVallarini/campaign-producer-api/internal/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SlugService interface {
	Fetch(context.Context, uuid.UUID) (model.Slug, error)
	Create(context.Context, model.SlugCreateRequest) error
	ListHistory(context.Context, uuid.UUID, model.ListingFilters) ([]model.SlugHistory, model.Paging, error)
}

type Slug struct {
	service SlugService
}

func NewSlug(service SlugService) *Slug {
	return &Slug{
		service: service,
	}
}

func (s Slug) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/slug/:id", s.Fetch)
	v1.POST("/slug", s.Create)
	v1.GET("/history/slug/:id", s.List)
}

func (s Slug) Fetch(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid slug id"), "Invalid slug id")
	}

	result, err := s.service.Fetch(c.Request().Context(), id)
	if err != nil {
		return ResponseError(c, err, "Failed to fetch slug")
	}

	return c.JSON(http.StatusOK, result)
}

func (s Slug) Create(c echo.Context) error {
	req := model.SlugCreateRequest{}

	if err := c.Bind(&req); err != nil {

		return ResponseError(c, err, "Failed to parse slug input")
	}

	userEmail := util.GetHeaderFromContext(c, "x-user-email")
	req.CreatedBy = userEmail
	if userEmail == "" {
		return ResponseError(c, model.ErrInvalid.Wrap(nil, "Invalid request"), "Email is required")
	}

	err := s.service.Create(c.Request().Context(), req)
	if err != nil {

		return ResponseError(c, err, "Failed to create slug")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (s Slug) List(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid slug id"), "Invalid slug id")
	}

	filters := model.ListingFilters{}
	if err := c.Bind(&filters); err != nil {

		return ResponseError(c, err, "Failed to parse filters")
	}

	result, paging, err := s.service.ListHistory(c.Request().Context(), id, filters)
	if err != nil {

		return ResponseError(c, err, "Failed to list history slug record")
	}

	return c.JSON(http.StatusOK, model.PagedResults{Paging: paging, Data: result})
}
