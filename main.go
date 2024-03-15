package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"contact.app/handlers"
	"contact.app/services"
)

func main() {
	e := echo.New()

	cs := &services.ContactService{
		Contacts: FakeContacts,
	}

	handler := handlers.NewContactsHandler(cs)

	e.Static("/static", "static")

	e.GET("/", handler.RedirectToContacts)

	e.GET("/404", handler.NotFoundView)

	e.GET("/oops", handler.OopsView)

	e.GET("/contacts", handler.ContactsView)

	e.GET("/contacts/new", handler.NewContactView)

	e.POST("/contacts/new", handler.NewContactView)

	e.GET("/contacts/:id", handler.ContactDetailView)

	e.GET("/contacts/:id/edit", handler.EditContactView)

	e.POST("/contacts/:id/edit", handler.EditContactView)

	e.DELETE("/contacts/:id", handler.DeleteContact)

	e.GET("/contacts/validate/email", handler.ValidateEmail)

	e.RouteNotFound("/", handler.NotFoundView)
	e.RouteNotFound("/*", handler.NotFoundView)

	err := e.Start(":3000")
	if err != nil {
		log.Panic(err)
	}
}
