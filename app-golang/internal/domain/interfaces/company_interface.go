package interfaces

import (
	"context"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICompanyRepository interface {
	Create(c *entities.Company, ctx context.Context) (*entities.Company, error)
	FindByDocument(document string, ctx context.Context) (*entities.Company, error)
	FindMany(ctx context.Context) ([]*entities.Company, error)
	FindById(id primitive.ObjectID, ctx context.Context) (*entities.Company, error)
	Save(id primitive.ObjectID, c *entities.Company, ctx context.Context) error
	Delete(id primitive.ObjectID, ctx context.Context) error
}
