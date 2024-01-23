package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	handler := Handler{}

	e.GET("/", handler.RedirectToContacts)

	e.Start(":3000")
}

type Handler struct{}

func (h Handler) RedirectToContacts(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/contacts")
}
