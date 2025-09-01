package inmemoryrepository

import (
	"errors"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
)

type MockRepository struct {
	Companies map[string]*entities.Company
}

func (m *MockRepository) Create(c *entities.Company) (*entities.Company, error) {
	m.Companies[c.Document] = c
	return c, nil
}

func (m *MockRepository) FindByDocument(document string) (*entities.Company, error) {
	c, ok := m.Companies[document]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *MockRepository) FindMany() ([]*entities.Company, error) {
	result := make([]*entities.Company, 0, len(m.Companies))
	for _, c := range m.Companies {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockRepository) FindById(id string) (*entities.Company, error) {
	for _, c := range m.Companies {
		if c.ID.Hex() == id {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) Save(id string, c *entities.Company) (bool, error) {
	for doc, company := range m.Companies {
		if company.ID.Hex() == id {
			c.ID = company.ID
			m.Companies[doc] = c
			return true, nil
		}
	}
	return false, nil
}

func (m *MockRepository) Delete(id string) error {
	for key, company := range m.Companies {
		if company.ID.Hex() == id {
			delete(m.Companies, key)
			return nil
		}
	}
	return errors.New("company not found")
}
