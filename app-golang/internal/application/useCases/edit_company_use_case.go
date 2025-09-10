package usecases

import (
	"context"
	"fmt"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	valueobjects "github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/valueObjects"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EditCompany struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (ec *EditCompany) Handle(objId primitive.ObjectID, c *entities.Company, ctx context.Context) (bool, error) {
	existsCompany, err := ec.Repo.FindById(objId, ctx)

	if err != nil {
		return false, fmt.Errorf("fail to find by id: %w", err)
	}

	if existsCompany == nil {
		return false, fails.ErrCompanyNotFound
	}

	documentUsed, err := ec.Repo.FindByDocument(c.Document, ctx)

	if err != nil {
		return false, fmt.Errorf("fail to find by document: %w", err)
	}

	if documentUsed != nil {
		return false, fails.ErrCompanyDocumentIsAlreadyInUse
	}

	err = c.CheckPWDQuota(c.TotalEmployees, c.TotalEmployeesPwd)

	if err != nil {
		return false, fmt.Errorf("fail to consult pwd quota: %w", err)
	}

	c.PrepareForUpdate()

	err = ec.Repo.Save(objId, c, ctx)

	if err != nil {
		return false, fmt.Errorf("fail to updated: %w", err)
	}

	e := valueobjects.Event{
		Message: "company_edited",
		EventId: uuid.New().String(),
		Data:    existsCompany,
	}

	encodedEvent, err := e.ToJson()

	if err != nil {
		return false, fmt.Errorf("fail convert event to json: %w", err)
	}

	ec.Broker.Publish("businesshub-logger", encodedEvent)

	return true, nil
}
