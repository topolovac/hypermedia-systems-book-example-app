package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"contact.app/model"
	"contact.app/templates"
	"contact.app/utils"
)

func main() {
	e := echo.New()

	handler := Handler{}

	e.GET("/", handler.RedirectToContacts)

	e.GET("/contacts", handler.ContactsView)

	e.Start(":3000")
}

type Handler struct{}

func (h Handler) RedirectToContacts(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/contacts")
}

func (h Handler) ContactsView(c echo.Context) error {

	contacts := []model.Contact{
		{Id: 1, First: "John", Last: "Doe", Phone: "555-555-5555", Email: "johndoe@email.com"},
		{Id: 2, First: "Jane", Last: "Doe", Phone: "555-555-5555", Email: "janedoe@email.com"},
	}

	return utils.Render(c, templates.Contacts(contacts))
}
