package usecases

import (
	"context"
	"fmt"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	valueobjects "github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/valueObjects"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteCompany struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (dc *DeleteCompany) Handle(objId primitive.ObjectID, ctx context.Context) error {
	existsCompany, err := dc.Repo.FindById(objId, ctx)

	if err != nil {
		return fmt.Errorf("fail to find by id: %w", err)
	}

	if existsCompany == nil {
		return fails.ErrCompanyNotFound
	}

	err = dc.Repo.Delete(objId, ctx)

	if err != nil {
		return fmt.Errorf("fail to delete: %w", err)
	}

	e := valueobjects.Event{
		Message: "company_deleted",
		EventId: uuid.New().String(),
		Data:    existsCompany,
	}

	encodedEvent, err := e.ToJson()

	if err != nil {
		return fmt.Errorf("fail to convert event to json: %w", err)
	}

	dc.Broker.Publish("businesshub-logger", encodedEvent)

	return err
}
