package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SpentService interface {
	FetchByMerchantIdAndBucket(context.Context, uuid.UUID, string) (model.Spent, error)
	Create(context.Context, model.SpentEvent) error
}

type Spent struct {
	service SpentService
}

func NewSpent(service SpentService) *Spent {
	return &Spent{
		service: service,
	}
}

func (s Spent) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/spent/:merchantId", s.Fetch)
	v1.POST("/spent", s.Create)
}

func (s Spent) Fetch(c echo.Context) error {
	merchantId, err := uuid.Parse(c.Param("merchantId"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid merchant id"), "Invalid merchant id")
	}

	bucket := c.QueryParam("bucket")

	result, err := s.service.FetchByMerchantIdAndBucket(c.Request().Context(), merchantId, bucket)
	if err != nil {

		return ResponseError(c, err, "Failed to fetch merchant")
	}

	return c.JSON(http.StatusOK, result)
}

func (s Spent) Create(c echo.Context) error {
	req := model.SpentEvent{}

	if err := c.Bind(&req); err != nil {

		return ResponseError(c, err, "Failed to parse spent input")
	}

	err := s.service.Create(c.Request().Context(), req)
	if err != nil {

		return ResponseError(c, err, "Failed to create spent")
	}

	return c.JSON(http.StatusCreated, nil)
}
