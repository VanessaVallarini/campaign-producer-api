package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/VanessaVallarini/campaign-producer-api/internal/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MerchantService interface {
	Fetch(context.Context, uuid.UUID) (model.Merchant, error)
	Create(context.Context, model.MerchantCreateRequest) error
}

type Merchant struct {
	service MerchantService
}

func NewMerchant(service MerchantService) *Merchant {
	return &Merchant{
		service: service,
	}
}

func (m Merchant) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/merchant/:id", m.Fetch)
	v1.POST("/merchant", m.Create)
}

func (m Merchant) Fetch(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid merchant id"), "Invalid merchant id")
	}

	result, err := m.service.Fetch(c.Request().Context(), id)
	if err != nil {

		return ResponseError(c, err, "Failed to fetch merchant")
	}

	return c.JSON(http.StatusOK, result)
}

func (m Merchant) Create(c echo.Context) error {
	req := model.MerchantCreateRequest{}

	if err := c.Bind(&req); err != nil {

		return ResponseError(c, err, "Failed to parse merchant input")
	}

	userEmail := util.GetHeaderFromContext(c, "x-user-email")
	req.CreatedBy = userEmail
	if userEmail == "" {
		return ResponseError(c, model.ErrInvalid.Wrap(nil, "Invalid request"), "Email is required")
	}

	err := m.service.Create(c.Request().Context(), req)
	if err != nil {

		return ResponseError(c, err, "Failed to create merchant")
	}

	return c.JSON(http.StatusCreated, nil)
}
