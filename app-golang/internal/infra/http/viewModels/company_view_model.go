package viewmodels

import "time"

type CompanyOutput struct {
	Id                time.Time     `json:"id" example:"68ba25fe8cd931d0437112ad"`
	Document          string        `json:"document" example:"19862056000223"`
	FantasyName       string        `json:"fantasy_name" example:"My Company"`
	SocialReason      string        `json:"social_reason" example:"My Company LTDA"`
	Address           AddressOutput `json:"address"`
	TotalEmployees    int           `json:"total_employees" example:"100"`
	TotalEmployeesPwd int           `json:"total_employees_pwd" example:"1"`
	CreatedAt         time.Time     `json:"created_at" example:"2025-09-04T23:51:26.881469869Z"`
	UpdatedAt         time.Time     `json:"updated_at" example:"2025-09-04T23:51:26.881469869Z"`
}

type AddressOutput struct {
	Street       string `json:"street" example:"Rua teste"`
	Complement   string `json:"complement,omitempty" example:"Perto do Mercado X"`
	PostalCode   string `json:"postal_code" example:"123-45678"`
	State        string `json:"state" example:"SP"`
	City         string `json:"city" example:"Maua"`
	Neighborhood string `json:"neighborhood" example:"Jardins"`
}

type CompanyCreatedResponse struct {
	Success bool          `json:"success"`
	Data    CompanyOutput `json:"data"`
}

type CompanyUpdatedResponse struct {
	Success bool          `json:"success"`
	Data    CompanyOutput `json:"data"`
}

type CompaniesListResponse struct {
	Success bool            `json:"success"`
	Data    []CompanyOutput `json:"data"`
}

type CompanyNotFoundResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error" example:"company not found"`
}

type CompanyConflitResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error" example:"company already exists"`
}

type CompanyBadRequestResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error" example:"Insufficient quota: company must have 10 PWD(s), but has 5 (required by Brazilian Law nº 8.213/91, art. 93)"`
}
