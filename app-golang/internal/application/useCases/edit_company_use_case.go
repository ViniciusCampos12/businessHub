package usecases

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	"github.com/google/uuid"
)

type EditCompany struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (ec *EditCompany) Handle(id string, c *entities.Company) (bool, error) {
	existsCompany, err := ec.Repo.FindById(id)

	if err != nil {
		return false, err
	}

	if existsCompany == nil {
		return false, errors.New("company not found")
	}

	err = c.CheckPWDQuota(c.TotalEmployees, c.TotalEmployeesPwd)

	if err != nil {
		return false, err
	}

	c.Document = existsCompany.Document
	c.Address.UnsmaskPostalCode()
	c.UpdatedAt = time.Now()

	hasUpdated, err := ec.Repo.Save(id, c)

	if err != nil {
		return false, err
	}

	if !hasUpdated {
		return false, errors.New("update failed")
	}

	message := map[string]interface{}{
		"Message": "company_edited",
		"EventId": uuid.New().String(),
		"Data": existsCompany,
	}

	payload, err := json.Marshal(message)

	if err != nil {
		return false, err
	}

	ec.Broker.Publish("businesshub-logger", payload)

	return true, nil
}
