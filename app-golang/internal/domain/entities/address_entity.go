package entities

import "strings"

// swagger:model Address
type Address struct {
	Street       string `json:"street" bson:"street" validate:"required" example:"Rua teste"`
	Complement   string `json:"complement,omitempty" bson:"complement,omitempty" example:"Perto do Mercado X"`
	PostalCode   string `json:"postal_code" bson:"postal_code" validate:"required,min=8,max=9" example:"123-45678"`
	State        string `json:"state" bson:"state" validate:"required,len=2" example:"SP"`
	City         string `json:"city" bson:"city" validate:"required" example:"Maua"`
	Neighborhood string `json:"neighborhood" bson:"neighborhood" validate:"required" example:"Jardins"`
}

func (a *Address) UnsmaskPostalCode() {
	postalCode := a.PostalCode
	postalCode = strings.ReplaceAll(postalCode, "-", "")

	a.PostalCode = postalCode
}
