package middleware

import "github.com/labstack/echo/v4"

type AuthContext struct {
	echo.Context
	UserID int64
	ClientID int64
}
