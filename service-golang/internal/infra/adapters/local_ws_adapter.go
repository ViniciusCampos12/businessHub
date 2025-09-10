package adapters

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type LocalWs struct{}

var ErrWebsocketConnectionFail = errors.New("websocket connection fail")
var ErrWebsocketFailedToConnect = errors.New("websocket failed to connect")

func (lw *LocalWs) SendMessage(url string, payload []byte) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWebsocketFailedToConnect, err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, payload)

	if err != nil {
		return fmt.Errorf("fail to send message: %w", err)
	}

	log.Info("message send successfully")

	return nil
}
