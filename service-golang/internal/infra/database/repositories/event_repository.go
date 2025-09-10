package repositories

import (
	"context"
	"fmt"

	"github.com/ViniciusCampos12/businessHub/service-golang/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository struct {
	Mongo *mongo.Database
	db    *mongo.Collection
}

func NewEventRepository(db *mongo.Database) (*EventRepository, error) {
	repo := &EventRepository{
		Mongo: db,
	}

	repo.db = repo.Mongo.Collection("outbox_events")

	if err := repo.ensureIndexes(); err != nil {
		return nil, fmt.Errorf("mongodb fail to ensure indexes: %w", err)
	}

	return repo, nil
}

func (e *EventRepository) ensureIndexes() error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"event_id": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := e.db.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (e *EventRepository) Create(c *entities.Event, ctx context.Context) error {
	_, err := e.db.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("mongodb fail to create: %w", err)
	}
	return nil
}
