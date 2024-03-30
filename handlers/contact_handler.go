package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	return c.Redirect(http.StatusSeeOther, "/contacts")
}

func (h *ContactsHandler) RedirectToNotFound(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/404")
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

	page := 1
	pp, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		fmt.Println(err)
	} else {
		page = pp
	}

	var contacts []model.Contact
	if search != "" {
		contacts = h.contact_service.Search(search)
	} else {
		contacts = h.contact_service.All(page)
	}

	if c.Request().Header.Get("HX-Trigger") == "search" {
		// added delay to simulate slower response
		time.Sleep(500 * time.Millisecond)
		return utils.Render(c, templates.ContactRows(contacts))
	}

	return utils.Render(c, templates.Contacts(contacts, search, page))
}

func (h *ContactsHandler) ContactsCount(c echo.Context) error {
	count := h.contact_service.GetLength()
	return c.HTML(200, fmt.Sprintf("<span>(%d)</span>", count))
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

func (h *ContactsHandler) EditContactView(c echo.Context) error {
	errors := templates.FormErrors{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.RedirectToNotFound(c)
	}

	existing_contact, err := h.contact_service.FindById(id)
	if err != nil {
		return h.RedirectToNotFound(c)
	}

	contact := &templates.EditContactForm{
		Id:    strconv.Itoa(int(existing_contact.Id)),
		Email: existing_contact.Email,
		First: existing_contact.First,
		Last:  existing_contact.Last,
		Phone: existing_contact.Phone,
	}

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
			err := h.contact_service.Update(model.Contact{
				Id:    existing_contact.Id,
				First: contact.First,
				Last:  contact.Last,
				Phone: contact.Phone,
				Email: contact.Email,
			})
			if err != nil {
				return h.RedirectToOops(c)
			}
			return h.RedirectToContacts(c)
		}
	}

	return utils.Render(c, templates.EditContactView(contact, errors))
}

func (h *ContactsHandler) DeleteContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println("id", id)
	if err != nil {
		fmt.Println("err1:", err)
		return h.RedirectToNotFound(c)
	}

	err = h.contact_service.Delete(id)
	if err != nil {
		fmt.Println("err2:", err)
		return h.RedirectToNotFound(c)
	}

	trigger := c.Request().Header.Get("HX-Trigger")
	fmt.Println(trigger)
	if trigger == "delete-btn" {
		return c.HTML(200, "")
	}

	return h.RedirectToContacts(c)
}

func (h *ContactsHandler) ValidateEmail(c echo.Context) error {
	email := c.QueryParam("email")

	if email == "" {
		return c.HTML(http.StatusOK, "email is required")
	}

	if _, err := h.contact_service.FindByEmail(email); err == nil {
		return c.HTML(http.StatusOK, "email is already taken")
	}

	if !utils.IsEmail(email) {
		return c.HTML(http.StatusOK, "email is invalid")
	}

	return c.HTML(http.StatusOK, "")
}
