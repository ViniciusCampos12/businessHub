package usecases

import (
	"encoding/json"
	"errors"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	"github.com/google/uuid"
)

type DeleteCompany struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (dc *DeleteCompany) Handle(id string) error {
	existsCompany, err := dc.Repo.FindById(id)

	if err != nil {
		return err
	}

	if existsCompany == nil {
		return errors.New("company not found")
	}

	err = dc.Repo.Delete(id)

	if err != nil {
		return err
	}

	message := map[string]interface{}{
		"Message": "company_deleted",
		"EventId": uuid.New().String(),
		"Data": existsCompany,
	}

	payload, err := json.Marshal(message)

	if err != nil {
		return err
	}

	dc.Broker.Publish("businesshub-logger", payload)

	return err
}
