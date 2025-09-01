package interfaces

import "github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"

type ICompanyRepository interface {
	Create(c *entities.Company) (*entities.Company, error)
	FindByDocument(document string) (*entities.Company, error)
	FindMany() ([]*entities.Company, error)
	FindById(id string) (*entities.Company, error)
	Save(id string, c *entities.Company) (bool, error)
	Delete(id string) error
}
