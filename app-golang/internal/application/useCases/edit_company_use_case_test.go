package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	inmemoryrepository "github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/database/inMemoryRepository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestShouldEditCompanyIfExists(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &usecases.EditCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	mongoId := primitive.NewObjectID()

	companyDummie := &entities.Company{
		ID:                mongoId,
		Document:          "99862056000112",
		FantasyName:       "Old Company",
		SocialReason:      "Old Company LTDA",
		TotalEmployees:    10,
		TotalEmployeesPwd: 1,
		Address: entities.Address{
			Street:       "Rua teste",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	mockRepo.Create(companyDummie, context.TODO())

	companyDummieUpdate := &entities.Company{
		FantasyName:       "My Company",
		SocialReason:      "My Company LTDA",
		TotalEmployees:    150,
		TotalEmployeesPwd: 25,
		Address: entities.Address{
			Street:       "Rua teste 2",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	hasUpdated, err := useCase.Handle(mongoId, companyDummieUpdate, context.TODO())

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !hasUpdated {
		t.Fatalf("expected true, got false")
	}

	edited, _ := mockRepo.FindByDocument("99862056000112", context.TODO())
	if edited.FantasyName != "My Company" || edited.TotalEmployees != 150 {
		t.Fatalf("company not updated correctly")
	}
}

func TestShouldNotEditCompanyIfNotExists(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &usecases.EditCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	mongoId := primitive.NewObjectID()

	companyDummieUpdate := &entities.Company{
		FantasyName:       "My Company",
		SocialReason:      "My Company LTDA",
		TotalEmployees:    15,
		TotalEmployeesPwd: 4,
		Address: entities.Address{
			Street:       "Rua teste 2",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	hasUpdated, err := useCase.Handle(mongoId, companyDummieUpdate, context.TODO())

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "company not found" {
		t.Fatalf("expected 'company not found', got %v", err)
	}

	if hasUpdated {
		t.Fatalf("expected false, got true")
	}

	if len(mockRepo.Companies) != 0 {
		t.Fatalf("expected empty repository, got %d items", len(mockRepo.Companies))
	}
}

func TestShouldNotEditCompanyIfPWDQuotaInvalid(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &usecases.EditCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	mongoId := primitive.NewObjectID()

	companyDummie := &entities.Company{
		ID:                mongoId,
		Document:          "99862056000112",
		FantasyName:       "Old Company",
		SocialReason:      "Old Company LTDA",
		TotalEmployees:    100,
		TotalEmployeesPwd: 5,
		Address: entities.Address{
			Street:       "Rua teste",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	mockRepo.Create(companyDummie, context.TODO())

	update := &entities.Company{
		FantasyName:       "My Company",
		SocialReason:      "My Company LTDA",
		TotalEmployees:    120,
		TotalEmployeesPwd: 2,
		Address: entities.Address{
			Street:       "Rua teste 2",
			Complement:   "string",
			PostalCode:   "12345678",
			Neighborhood: "Jardins",
			City:         "Maua",
			State:        "SP",
		},
	}

	hasUpdated, err := useCase.Handle(mongoId, update, context.TODO())

	if err == nil {
		t.Fatalf("expected error due to PWD quota, got nil")
	}

	if hasUpdated {
		t.Fatalf("expected false, got true")
	}
}
