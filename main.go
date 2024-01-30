package main

import (
	"log"
	"net/http"
	"strconv"

	"contact.app/model"
	"contact.app/services"
	"contact.app/templates"
	"contact.app/utils"
	"github.com/go-playground/form"
	"github.com/labstack/echo/v4"
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

	handler := &Handler{
		contact_service: cs,
		decoder:         form.NewDecoder(),
	}

	e.GET("/", handler.RedirectToContacts)

	e.GET("/404", handler.NotFoundView)

	e.GET("/oops", handler.OopsView)

	e.GET("/contacts", handler.ContactsView)

	e.GET("/contacts/new", handler.NewContactView)

	e.POST("/contacts/new", handler.NewContactView)

	e.GET("/contacts/:id", handler.ContactDetailView)

	e.RouteNotFound("/", handler.NotFoundView)
	e.RouteNotFound("/*", handler.NotFoundView)

	err := e.Start(":3000")
	if err != nil {
		log.Panic(err)
	}
}

type Handler struct {
	contact_service *services.ContactService
	decoder         *form.Decoder
}

func (h *Handler) RedirectToContacts(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/contacts")
}

func (h *Handler) RedirectToNotFound(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/404")
}

func (h *Handler) NotFoundView(c echo.Context) error {
	return utils.Render(c, templates.NotFound())
}

func (h *Handler) RedirectToOops(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/oops")
}

func (h *Handler) OopsView(c echo.Context) error {
	return utils.Render(c, templates.Oops())
}

func (h *Handler) ContactsView(c echo.Context) error {
	search := c.QueryParam("search")
	var contacts []model.Contact
	if search != "" {
		contacts = h.contact_service.Search(search)
	} else {
		contacts = h.contact_service.All()
	}
	return utils.Render(c, templates.Contacts(contacts, search))
}

func (h *Handler) NewContactView(c echo.Context) error {
	contact := &templates.NewContactForm{}
	errors := templates.FormErrors{}

	if c.Request().Method == http.MethodPost {
		values, err := c.FormParams()
		if err != nil {
			return c.Redirect(http.StatusUnprocessableEntity, c.Path())
		}

		err = h.decoder.Decode(&contact, values)

		if err != nil {
			return c.Redirect(http.StatusUnprocessableEntity, c.Path())
		}

		if contact.First == "" {
			errors.First = "Please enter first name"
		}
		if contact.Last == "" {
			errors.Last = "Please enter last name"
		}
		if contact.Phone == "" {
			errors.Phone = "Please enter phone number"
		}
		if contact.Email == "" || !utils.IsEmail(contact.Email) {
			errors.Email = "Please enter valid email"
		}

		if !errors.HasErrors() {
			h.contact_service.Add(model.Contact{
				First: contact.First,
				Last:  contact.Last,
				Phone: contact.Phone,
				Email: contact.Email,
			})
			return h.RedirectToContacts(c)
		}
	}

	return utils.Render(c, templates.NewContactView(contact, errors))
}

func (h *Handler) ContactDetailView(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.RedirectToNotFound(c)
	}

	contact, err := h.contact_service.FindById(id)
	if err != nil {
		return h.RedirectToNotFound(c)
	}

	return utils.Render(c, templates.ContactDetail(&contact))
}
