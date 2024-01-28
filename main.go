package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"contact.app/model"
	"contact.app/services"
	"contact.app/templates"
	"contact.app/utils"
)

func main() {
	e := echo.New()

	contacts := []model.Contact{
		{Id: 1, First: "John", Last: "Doe", Phone: "555-555-5555", Email: "johndoe@email.com"},
		{Id: 2, First: "Jane", Last: "Doe", Phone: "555-555-5555", Email: "janedoe@email.com"},
		{Id: 3, First: "Alice", Last: "Smith", Phone: "555-555-5555", Email: "alicesmith@email.com"},
		{Id: 4, First: "Bob", Last: "Johnson", Phone: "555-555-5555", Email: "bobjohnson@email.com"},
		{Id: 5, First: "Eva", Last: "Brown", Phone: "555-555-5555", Email: "evabrown@email.com"},
		{Id: 6, First: "David", Last: "Wilson", Phone: "555-555-5555", Email: "davidwilson@email.com"},
		{Id: 7, First: "Grace", Last: "Miller", Phone: "555-555-5555", Email: "gracemiller@email.com"},
	}

	cs := &services.ContactService{
		Contacts: contacts,
	}

	handler := Handler{
		contact_service: cs,
	}

	e.GET("/", handler.RedirectToContacts)

	e.GET("/contacts", handler.ContactsView)

	e.Start(":3000")
}

type Handler struct {
	contact_service *services.ContactService
}

func (h Handler) RedirectToContacts(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/contacts")
}

func (h Handler) ContactsView(c echo.Context) error {
	search := c.QueryParam("search")
	var contacts []model.Contact
	if search != "" {
		contacts = h.contact_service.Search(search)
	} else {
		contacts = h.contact_service.All()
	}
	return utils.Render(c, templates.Contacts(contacts, search))
}
