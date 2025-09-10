package entities

// swagger:model Address
type Address struct {
	Street       string `json:"street" bson:"street" binding:"required,min=2,max=120" example:"Rua teste"`
	Complement   string `json:"complement,omitempty" bson:"complement,omitempty" example:"Perto do Mercado X"`
	PostalCode   string `json:"postal_code" bson:"postal_code" binding:"required,len=8,numeric" example:"12345678"`
	State        string `json:"state" bson:"state" binding:"required,len=2,alpha" example:"SP"`
	City         string `json:"city" bson:"city" binding:"required,min=2,max=80" example:"Maua"`
	Neighborhood string `json:"neighborhood" bson:"neighborhood" binding:"required,min=2,max=80" example:"Jardins"`
}
