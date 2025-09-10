package jobs

import (
	"encoding/json"
	"errors"

	usecases "github.com/ViniciusCampos12/businessHub/service-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/domain/interfaces"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/infra/database/repositories"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartCompanyConsumer(queue string, c interfaces.IMessageBroker, mongoDB *mongo.Database) {
	messages, err := c.Consume(queue)

	if err != nil {
		log.Fatalf("error to start rabbitmq consumer: %v", err)
	}

	go consumorRabbitmqJob(messages, mongoDB)
}

func consumorRabbitmqJob(messages <-chan amqp091.Delivery, mongoDB *mongo.Database) {
	log.Info("started job")

	var event *entities.Event
	for msg := range messages {
		err := json.Unmarshal(msg.Body, &event)

		if err != nil {
			log.Errorf("error in dispatcher job: %v", err)
			continue
		}

		go dispatchMessageToWebsocketJob(event, mongoDB)
		log.Info("dispatching job to websocket")
	}
}

func dispatchMessageToWebsocketJob(e *entities.Event, mongoDB *mongo.Database) {
	repo, err := repositories.NewEventRepository(mongoDB)

	if err != nil {
		log.Errorf("erro start repo in job: %v", err)
	}

	sendEventWs := &usecases.SendEventWs{
		Repo:      repo,
		WsGateway: &adapters.LocalWs{},
	}

	err = sendEventWs.Handle(e)

	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrJobFailToCreate):
			log.Errorf("job create event: %v", err)
		case errors.Is(err, usecases.ErrJobFailToMarshal):
			log.Errorf("job marshal event: %v", err)
		case errors.Is(err, adapters.ErrWebsocketConnectionFail):
			log.Errorf("job send message websocket: %v", err)
		case errors.Is(err, adapters.ErrWebsocketFailedToConnect):
			log.Errorf("job send message websocket: %v", err)
		case errors.Is(err, usecases.ErrWebsocketEnvNotFound):
			log.Errorf("job get env websocket url: %v", err)
		default:
			log.Errorf("job fail: %v", err)
		}
	}
}
