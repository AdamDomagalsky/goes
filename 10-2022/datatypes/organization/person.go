package organization

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Identifiable interface {
	ID() string
}
type Citizen interface {
	Identifiable
	Country() string
}

type Name struct {
	first string
	last  string
}

type socialSecurityNumber string

func NewSocialSecurityNumber(value string) Citizen {
	return socialSecurityNumber(value)
}

func (ssn socialSecurityNumber) ID() string {
	return string(ssn)
}

func (ssn socialSecurityNumber) Country() string {
	return "United States of America"

}

type europeanUnionIdentifier struct {
	id      string
	country string
}

// interface{} // antyhing - not very safe
func NewEuropeanUnionIdentifier(id interface{}, country string) Citizen {
	switch v := id.(type) {
	case string:
		return europeanUnionIdentifier{
			id:      v,
			country: country,
		}
	case int:
		return europeanUnionIdentifier{
			id:      strconv.Itoa(v),
			country: country,
		}
	case europeanUnionIdentifier:
		return v
	case Person:
		euID, _ := v.Citizen.(europeanUnionIdentifier)
		return euID
	default:
		panic("using invalid type to initialize EU Identifier")
	}
}

func (eui europeanUnionIdentifier) ID() string {
	return eui.id
}

func (eui europeanUnionIdentifier) Country() string {
	return eui.country
}

type Person struct {
	Name
	twitterHandler TwitterHandler
	Citizen
}

type Employee struct {
	Name
}

type TwitterHandler2 = string // type Alias
type TwitterHandler string    // type Definition

func (th TwitterHandler) RedirectURL() string {
	cleanHandler := strings.TrimPrefix(string(th), "@")
	return fmt.Sprintf("https://www.twitter.com/%s", cleanHandler)
}

func NewPerson(firstName, lastName string, citizen Citizen) Person {
	return Person{
		Name: Name{
			first: firstName,
			last:  lastName,
		},
		Citizen: citizen,
	}
}

func (p *Person) SetTwitterHandler(handler TwitterHandler) error {
	if len(handler) == 0 {
		p.twitterHandler = handler
	} else if !strings.HasPrefix(string(handler), "@") {
		return errors.New("twitter handler must start with an @ symbol")
	}
	p.twitterHandler = handler
	return nil
}
func (p *Person) FullName() string {
	return fmt.Sprintf("%s %s", p.first, p.last)
}

func (p *Person) TwitterHandler() TwitterHandler {
	return p.twitterHandler
}

// @Overwrtielike
func (p Person) ID() string {
	return fmt.Sprintf("Wrapper of ID %v", p.Citizen.ID())
}
