package usecases

import (
	"testing"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	inmemoryrepository "github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/database/inMemoryRepository"
)

func TestShouldCreateCompanyIfNotExists(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &CreateCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	input := &entities.Company{
		Document:          "99862056000112",
		FantasyName:       "My Company",
		SocialReason:      "My Company LTDA",
		TotalEmployees:    10,
		TotalEmployeesPwd: 3,
		Address: entities.Address{
			Street:       "Rua teste",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	result, err := useCase.Handle(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected a company, got nil")
	}

	if mockRepo.Companies[result.Document] == nil {
		t.Fatal("expected a company, got nil")
	}

}

func TestShouldNotCreateCompanyIfAlreadyExists(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &CreateCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	input := &entities.Company{
		Document:          "99862056000112",
		FantasyName:       "My Company",
		SocialReason:      "My Company LTDA",
		TotalEmployees:    100,
		TotalEmployeesPwd: 2,
		Address: entities.Address{
			Street:       "Rua teste",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	mockRepo.Create(input)

	existingCompany, err := useCase.Handle(input)

	if existingCompany != nil {
		t.Fatalf("expected nil, got %v", existingCompany)
	}

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldNotCreateCompanyIfPWDQuotaIsInvalid(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &CreateCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	input := &entities.Company{
		Document:          "99862056000112",
		FantasyName:       "My Company",
		SocialReason:      "My Company LTDA",
		TotalEmployees:    100,
		TotalEmployeesPwd: 0,
		Address: entities.Address{
			Street:       "Rua teste",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	existingCompany, err := useCase.Handle(input)

	if existingCompany != nil {
		t.Fatalf("expected nil, got %v", existingCompany)
	}

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
