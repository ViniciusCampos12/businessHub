package usecases

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCompany struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (cc *CreateCompany) Handle(c *entities.Company) (*entities.Company, error) {
	c.UnsmaskDocument()
	existsCompany, err := cc.Repo.FindByDocument(c.Document)

	if err != nil {
		return nil, err
	}

	if existsCompany != nil {
		return nil, errors.New("company already exists")
	}

	err = c.CheckPWDQuota(c.TotalEmployees, c.TotalEmployeesPwd)

	if err != nil {
		return nil, err
	}

	c.Address.UnsmaskPostalCode()

	c.ID = primitive.NewObjectID()
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Now()

	newCompany, err := cc.Repo.Create(c)

	if err != nil {
		return nil, err
	}

	message := map[string]interface{}{
		"Message": "company_created",
		"EventId": uuid.New().String(),
		"Data": c,
	}

	payload, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	cc.Broker.Publish("businesshub-logger", payload)

	return newCompany, nil

}
