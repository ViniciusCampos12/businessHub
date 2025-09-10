package interfaces

import (
	"context"

	"github.com/ViniciusCampos12/businessHub/service-golang/internal/domain/entities"
)

type IEventRepository interface {
	Create(e *entities.Event, ctx context.Context) error
}
