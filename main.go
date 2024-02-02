package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"contact.app/handlers"
	"contact.app/model"
	"contact.app/services"
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

	handler := handlers.NewContactsHandler(cs)

	e.GET("/", handler.RedirectToContacts)

	e.GET("/404", handler.NotFoundView)

	e.GET("/oops", handler.OopsView)

	e.GET("/contacts", handler.ContactsView)

	e.GET("/contacts/new", handler.NewContactView)

	e.POST("/contacts/new", handler.NewContactView)

	e.GET("/contacts/:id", handler.ContactDetailView)

	e.GET("/contacts/:id/edit", handler.EditContactView)

	e.POST("/contacts/:id/edit", handler.EditContactView)

	e.RouteNotFound("/", handler.NotFoundView)
	e.RouteNotFound("/*", handler.NotFoundView)

	err := e.Start(":3000")
	if err != nil {
		log.Panic(err)
	}
}
