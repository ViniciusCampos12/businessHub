package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ViniciusCampos12/businessHub/service-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/domain/interfaces"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/infra/gateways"
)

type SendEventWs struct {
	Repo      interfaces.IEventRepository
	WsGateway gateways.Websocket
}

var ErrJobFailToCreate = errors.New("fail to create event in db")
var ErrJobFailToMarshal = errors.New("fail to marshal event")
var ErrWebsocketEnvNotFound = errors.New("fail to get env websocket url")

func (uc *SendEventWs) Handle(e *entities.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	e.PrepareToCreate()
	err := uc.Repo.Create(e, ctx)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrJobFailToCreate, err)
	}

	encoded, err := json.Marshal(e)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrJobFailToMarshal, err)
	}

	url := os.Getenv("WEBSOCKET_URL")

	if url == "" {
		return ErrWebsocketEnvNotFound
	}

	err = uc.WsGateway.SendMessage(url, encoded)

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
