package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/VanessaVallarini/campaign-producer-api/internal/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OwnerService interface {
	Fetch(context.Context, uuid.UUID) (model.Owner, error)
	Create(context.Context, model.OwnerCreateRequest) error
}

type Owner struct {
	service OwnerService
}

func NewOwner(service OwnerService) *Owner {
	return &Owner{
		service: service,
	}
}

func (o Owner) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/owner/:id", o.Fetch)
	v1.POST("/owner", o.Create)
}

func (o Owner) Create(c echo.Context) error {
	req := model.OwnerCreateRequest{}

	if err := c.Bind(&req); err != nil {

		return ResponseError(c, err, "Failed to parse owner input")
	}

	userEmail := util.GetHeaderFromContext(c, "x-user-email")
	req.CreatedBy = userEmail
	if userEmail == "" {
		return ResponseError(c, model.ErrInvalid.Wrap(nil, "Invalid request"), "Email is required")
	}

	err := o.service.Create(c.Request().Context(), req)
	if err != nil {

		return ResponseError(c, err, "Failed to create owner")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (o Owner) Fetch(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid owner id"), "Invalid owner id")
	}

	result, err := o.service.Fetch(c.Request().Context(), id)
	if err != nil {
		return ResponseError(c, err, "Failed to fetch owner")
	}

	return c.JSON(http.StatusOK, result)
}
