package handlers

import (
	"net/http"
	"strconv"

	"contact.app/model"
	"contact.app/services"
	"contact.app/templates"
	"contact.app/utils"
	"github.com/go-playground/form"
	"github.com/labstack/echo/v4"
)

func NewContactsHandler(cs *services.ContactService) *ContactsHandler {
	decoder := form.NewDecoder()
	return &ContactsHandler{
		contact_service: cs,
		decoder:         decoder,
	}

}

type ContactsHandler struct {
	contact_service *services.ContactService
	decoder         *form.Decoder
}

func (h *ContactsHandler) RedirectToContacts(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/contacts")
}

func (h *ContactsHandler) RedirectToNotFound(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/404")
}

func (h *ContactsHandler) NotFoundView(c echo.Context) error {
	return utils.Render(c, templates.NotFound())
}

func (h *ContactsHandler) RedirectToOops(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/oops")
}

func (h *ContactsHandler) OopsView(c echo.Context) error {
	return utils.Render(c, templates.Oops())
}

func (h *ContactsHandler) ContactsView(c echo.Context) error {
	search := c.QueryParam("search")
	var contacts []model.Contact
	if search != "" {
		contacts = h.contact_service.Search(search)
	} else {
		contacts = h.contact_service.All()
	}
	return utils.Render(c, templates.Contacts(contacts, search))
}

func (h *ContactsHandler) NewContactView(c echo.Context) error {
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

func (h *ContactsHandler) ContactDetailView(c echo.Context) error {
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
