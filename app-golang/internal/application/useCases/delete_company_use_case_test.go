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

func TestShouldDeleteCompanyIfExists(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &usecases.DeleteCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

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

	err := useCase.Handle(mongoId, context.TODO())

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, foundErr := mockRepo.FindById(mongoId, context.TODO())
	if foundErr != nil {
		t.Fatalf("unexpected error: %v", foundErr)
	}

	if len(mockRepo.Companies) != 0 {
		t.Fatalf("expected repository to be empty, got %d items", len(mockRepo.Companies))
	}
}

func TestShouldNotDeleteCompanyIfNotExists(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &usecases.DeleteCompany{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	mongoId := primitive.NewObjectID()

	err := useCase.Handle(mongoId, context.TODO())

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "company not found" {
		t.Fatalf("expected 'company not found', got %v", err)
	}
}
