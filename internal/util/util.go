package util

import "github.com/labstack/echo/v4"

func GetHeaderFromContext(ctx echo.Context, headerName string) string {
	return ctx.Request().Header.Get(headerName)
}
