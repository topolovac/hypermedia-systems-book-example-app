package services

import (
	"strings"

	"contact.app/model"
)

type ContactService struct {
	Contacts []model.Contact
}

func (c *ContactService) Search(search_param string) []model.Contact {
	matched_contacts := []model.Contact{}
	for _, contact := range c.Contacts {
		if strings.Contains(contact.Email, search_param) {
			matched_contacts = append(matched_contacts, contact)
		} else if strings.Contains(contact.First, search_param) {
			matched_contacts = append(matched_contacts, contact)
		} else if strings.Contains(contact.Last, search_param) {
			matched_contacts = append(matched_contacts, contact)
		}
	}
	return matched_contacts
}

func (c *ContactService) All() []model.Contact {
	return c.Contacts
}
