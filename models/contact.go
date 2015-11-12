// Definition of the structures and SQL interaction functions
package models

import (
	"time"

	// "github.com/asaskevich/govalidator"
)

// Position represents the coordinates of a contact
type Position struct {
	X float64
	Y float64

	Latitude  float64
	Longitude float64
}

// Address represents all the components of a contact's address
type Address struct {
	Position

	ID uint

	HouseNumber string
	Street      string
	PostalCode  string
	City        string
	County      string // Département
	State       string // Région
	Country     string
	Addition    string // Complément d'adresse

	PollingStation string // Code bureau de vote
}

// Contact represents all the components of a contact
type Contact struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Firstname   string     `sql:"not null" json:"firstname"`
	Surname     string     `sql:"not null" json:"surname"`
	MarriedName *string    `db:"married_name" json:"married_name,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	Birthdate   *time.Time `json:"birthdate,omitempty"`
	Mail        *string    `json:"mail,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Mobile      *string    `json:"mobile,omitempty"`
	Adress      *string    `json:"adress,omitempty"`

	Vote    string `json:"vote"`
	Support string `json:"support"`

	GroupID uint `sql:"not null" db:"group_id" json:"-"`

	Notes []Note `json:"notes,omitempty"`
	Tags  []Tag  `json:"tags,omitempty" gorm:"many2many:contact_tags;"`
}

// ContactArgs is used in the RPC communications between the gateway and Contacts
type ContactArgs struct {
	MissionID uint
	Contact   *Contact
}

// ContactReply is used in the RPC communications between the gateway and Contacts
type ContactReply struct {
	Contact  *Contact
	Contacts []Contact
}

// Validate checks if the contact is valid
func (c *Contact) Validate() map[string]string {
	var errs = make(map[string]string)

	// if c.Firstname == "" {
	// 	errs["firstname"] = "is required"
	// }

	// if c.Surname == "" {
	// 	errs["surname"] = "is required"
	// }

	// if c.Mail != nil && !govalidator.IsEmail(*c.Mail) {
	// 	errs["mail"] = "is not valid"
	// }

	return errs
}
