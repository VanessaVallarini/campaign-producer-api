package api

import (
	"context"
	"net/http"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LedgerService interface {
	List(context.Context, uuid.UUID, model.ListingFilters) ([]model.Ledger, model.Paging, error)
}

type Ledger struct {
	service LedgerService
}

func NewLedger(service LedgerService) *Ledger {
	return &Ledger{
		service: service,
	}
}

func (l Ledger) Register(server *echo.Echo) {
	v1 := server.Group("v1")
	v1.GET("/ledger/:campaignId", l.List)
}

func (l Ledger) List(c echo.Context) error {
	id, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {

		return ResponseError(c, model.ErrInvalid.Wrap(err, "Invalid campaign id"), "Invalid campaign id")
	}

	filters := model.ListingFilters{}
	if err := c.Bind(&filters); err != nil {

		return ResponseError(c, err, "Failed to parse filters")
	}

	result, paging, err := l.service.List(c.Request().Context(), id, filters)
	if err != nil {

		return ResponseError(c, err, "Failed to list ledger record")
	}

	return c.JSON(http.StatusOK, model.PagedResults{Paging: paging, Data: result})
}
