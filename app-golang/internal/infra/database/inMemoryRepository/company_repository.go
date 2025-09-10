package inmemoryrepository

import (
	"context"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepository struct {
	Companies map[string]*entities.Company
}

func (m *MockRepository) Create(c *entities.Company, ctx context.Context) (*entities.Company, error) {
	m.Companies[c.Document] = c
	return c, nil
}

func (m *MockRepository) FindByDocument(document string, ctx context.Context) (*entities.Company, error) {
	c, ok := m.Companies[document]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *MockRepository) FindMany(ctx context.Context) ([]*entities.Company, error) {
	result := make([]*entities.Company, 0, len(m.Companies))
	for _, c := range m.Companies {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockRepository) FindById(objID primitive.ObjectID, ctx context.Context) (*entities.Company, error) {
	for _, c := range m.Companies {
		if c.ID == objID {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) Save(objID primitive.ObjectID, c *entities.Company, ctx context.Context) error {
	for doc, company := range m.Companies {
		if company.ID == objID {
			c.ID = company.ID
			m.Companies[doc] = c
			return nil
		}
	}
	return nil
}

func (m *MockRepository) Delete(objID primitive.ObjectID, ctx context.Context) error {
	for key, company := range m.Companies {
		if company.ID == objID {
			delete(m.Companies, key)
			return nil
		}
	}
	return fails.ErrCompanyNotFound
}
