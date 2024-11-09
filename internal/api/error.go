package api

import (
	"context"
	"net/http"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const layout = "2006-01-02T15:04:05.999999Z07:00"

func ErrorHandlerResponse(err error) *Error {
	var errorResponse *Error

	baseError := &model.BaseErrorWrapper{}
	errors.As(err, &baseError)

	switch true {
	case errors.Is(err, model.ErrNotFound):
		errorResponse = NewError(http.StatusNotFound, baseError.ErrorDetails())
	case errors.Is(err, model.ErrInvalid):
		errorResponse = NewError(http.StatusBadRequest, baseError.ErrorDetails())
	case errors.Is(err, model.ErrInternal):
		errorResponse = NewError(http.StatusInternalServerError, baseError.ErrorDetails())
	case errors.Is(err, model.ErrForbidden):
		errorResponse = NewError(http.StatusForbidden, baseError.ErrorDetails())
	default:
		errorResponse = NewError(http.StatusInternalServerError, model.Error{Message: err.Error(), Code: "INTERNAL_ERROR"})
	}

	return errorResponse
}

type Error struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
	err        model.Error
}

func NewError(statusCode int, err model.Error) *Error {
	return &Error{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Timestamp:  time.Now().Format(layout),
		err:        err,
	}
}

type logFunc func(context.Context, error, string, ...zapcore.Field)

func ResponseError(c echo.Context, err error, message string) error {
	errorResponse := ErrorHandlerResponse(err)
	return c.JSON(errorResponse.StatusCode, errorResponse)
}
