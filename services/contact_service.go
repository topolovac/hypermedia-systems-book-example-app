package services

import (
	"errors"
	"strings"
	"time"

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

func (c *ContactService) All(page int) []model.Contact {
	if page > 0 {
		start := (page - 1) * 10
		end := page * 10
		if end > len(c.Contacts) {
			end = len(c.Contacts)
		}
		return c.Contacts[start:end]
	}
	return c.Contacts
}

func (c *ContactService) Add(contact model.Contact) {
	c.Contacts = append(c.Contacts, contact)
}

func (c *ContactService) FindById(id int) (model.Contact, error) {
	for _, c := range c.Contacts {
		if c.Id == int32(id) {
			return c, nil
		}
	}
	return model.Contact{}, errors.New("contact not found")
}

func (c *ContactService) Update(new_contact model.Contact) error {
	for i, contact := range c.Contacts {
		if contact.Id == new_contact.Id {
			c.Contacts[i] = new_contact
			return nil
		}
	}
	return errors.New("contact not found")
}

func (c *ContactService) Delete(id int) error {
	for i, contact := range c.Contacts {
		if contact.Id == int32(id) {
			c.Contacts = append(c.Contacts[:i], c.Contacts[i+1:]...)
			return nil
		}
	}
	return errors.New("contact not found")
}

func (c *ContactService) FindByEmail(email string) (model.Contact, error) {
	for _, contact := range c.Contacts {
		if contact.Email == email {
			return contact, nil
		}
	}
	return model.Contact{}, errors.New("contact not found")
}

func (c *ContactService) GetLength() int {
	time.Sleep(1 * time.Second)
	return len(c.Contacts)
}
