package adapters

import "errors"

type MockPublisher struct {
	Fail bool
}

func (m *MockPublisher) Publish(queue string, body []byte) error {
	if m.Fail {
		return errors.New("failed to publish: connection/channel error")
	}

	return nil
}
