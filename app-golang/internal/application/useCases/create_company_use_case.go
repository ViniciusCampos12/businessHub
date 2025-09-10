package usecases

import (
	"context"
	"fmt"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/interfaces"
	valueobjects "github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/valueObjects"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"github.com/google/uuid"
)

type CreateCompany struct {
	Repo   interfaces.ICompanyRepository
	Broker interfaces.IMessageBroker
}

func (cc *CreateCompany) Handle(c *entities.Company, ctx context.Context) (*entities.Company, error) {
	existsCompany, err := cc.Repo.FindByDocument(c.Document, ctx)

	if err != nil {
		return nil, fmt.Errorf("fail to find by document: %w", err)
	}

	if existsCompany != nil {
		return nil, fails.ErrCompanyAlreadyExists
	}

	err = c.CheckPWDQuota(c.TotalEmployees, c.TotalEmployeesPwd)

	if err != nil {
		return nil, fmt.Errorf("fail to consult pwd quota: %w", err)
	}

	c.PrepareForCreate()
	newCompany, err := cc.Repo.Create(c, ctx)

	if err != nil {
		return nil, fmt.Errorf("fail to create new company: %w", err)
	}

	e := valueobjects.Event{
		Message: "company_created",
		EventId: uuid.New().String(),
		Data:    c,
	}
	encodedEvent, err := e.ToJson()

	if err != nil {
		return nil, fmt.Errorf("fail convert event to json : %w", err)
	}

	cc.Broker.Publish("businesshub-logger", encodedEvent)

	return newCompany, nil
}
