package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	valueobjects "github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/valueObjects"
	"github.com/google/uuid"
)

type ListCompanies struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (lc *ListCompanies) Handle(ctx context.Context) ([]*entities.Company, error) {
	companies, err := lc.Repo.FindMany(ctx)

	if err != nil {
		return nil, fmt.Errorf("fail to find many: %w", err)
	}

	companiesEncoded, err := json.Marshal(companies)

	if err != nil {
		return nil, fmt.Errorf("fail to convert companies to json : %w", err)
	}

	e := valueobjects.Event{
		Message: "companies_readed",
		EventId: uuid.New().String(),
		Data:    companiesEncoded,
	}

	encodedEvent, err := e.ToJson()

	if err != nil {
		return nil, fmt.Errorf("fail convert event to json : %w", err)
	}

	lc.Broker.Publish("businesshub-logger", encodedEvent)

	return companies, nil
}
