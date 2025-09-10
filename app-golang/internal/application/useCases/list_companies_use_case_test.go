package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	inmemoryrepository "github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/database/inMemoryRepository"
)

func TestReturnListOfCompanies(t *testing.T) {
	mockRepo := &inmemoryrepository.MockRepository{Companies: make(map[string]*entities.Company)}
	useCase := &usecases.ListCompanies{Repo: mockRepo, Broker: &adapters.MockPublisher{Fail: false}}

	companies, err := useCase.Handle(context.TODO())

	if err != nil {
		t.Fatalf("expect nil, got %v", err)
	}

	if companies == nil {
		t.Fatalf("expect list, got nil")
	}
}
