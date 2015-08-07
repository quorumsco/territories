package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Position struct {
	X float64
	Y float64

	Latitude  float64
	Longitude float64
}

type Address struct {
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

	Position
}

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

	UserID uint `sql:"not null" db:"user_id" json:"-"`

	Address *Address `json:"address,omitempty"`
	Notes   []Note   `json:"notes,omitempty"`
	Tags    []Tag    `json:"tags,omitempty" gorm:"many2many:user_tags;"`
}

func (c *Contact) Validate() map[string]string {
	var errs = make(map[string]string)

	if c.Firstname == "" {
		errs["firstname"] = "is required"
	}

	if c.Surname == "" {
		errs["surname"] = "is required"
	}

	if c.Mail != nil && !govalidator.IsEmail(*c.Mail) {
		errs["mail"] = "is not valid"
	}

	return errs
}
