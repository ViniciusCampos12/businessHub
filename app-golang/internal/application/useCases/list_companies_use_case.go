package usecases

import (
	"encoding/json"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	"github.com/google/uuid"
)

type ListCompanies struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (lc *ListCompanies) Handle() ([]*entities.Company, error) {
	companies, err := lc.Repo.FindMany()

	if err != nil {
		return nil, err
	}

	companiesEncoded, err := json.Marshal(companies)

	if err != nil {
		return nil, err
	}

	message := map[string]interface{}{
		"Message": "companies_readed",
		"EventId": uuid.New().String(),
		"Data": companiesEncoded,
	}

	payload, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	lc.Broker.Publish("businesshub-logger", payload)

	return companies, nil
}
