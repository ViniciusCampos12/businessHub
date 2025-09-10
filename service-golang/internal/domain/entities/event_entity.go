package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Unix = int64

type Event struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventId    string             `json:"event_id" bson:"event_id"`
	Message    string             `json:"message" bson:"message"`
	Timestamp  Unix               `json:"timestamp" bson:"timestamp"`
	ReceivedAt time.Time          `json:"received_at" bson:"received_at"`
	Data       interface{}        `json:"data" bson:"data"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

func (e *Event) PrepareToCreate() {
	e.ID = primitive.NewObjectID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	e.ReceivedAt = time.Now()
}
