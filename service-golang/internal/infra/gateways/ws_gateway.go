package gateways

type Websocket interface {
	SendMessage(url string, payload []byte) error
}	
