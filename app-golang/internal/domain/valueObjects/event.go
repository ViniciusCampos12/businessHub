package valueobjects

import (
	"encoding/json"
	"fmt"
	"time"
)

type Unix = int64

type Event struct {
	Message   string      `json:"message"`
	EventId   string      `json:"event_id"`
	Timestamp Unix        `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func (e Event) ToJson() ([]byte, error) {
	e.Timestamp = time.Now().Unix()
	payload, err := json.Marshal(e)

	if err != nil {
		return nil, fmt.Errorf("fail to marshal payload rabbitmq: %w", err)
	}

	return payload, nil
}
