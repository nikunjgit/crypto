package client

import (
	ws "github.com/gorilla/websocket"
	exchange "github.com/preichenberger/go-coinbase-exchange"
)

const (
	GDAX_SOCKET_URL = "wss://ws-feed.gdax.com"
	SUBSCRIBE_TYPE = "subscribe"
	MESSAGE_MATCH = "match"
)
type SocketClient struct {
	gdaxClient *exchange.Client
	wsConn     *ws.Conn
}

type SubscribeOptions struct {
	Product_id string
}

func NewSocketClient(gdaxClient *exchange.Client, options *SubscribeOptions) (*SocketClient, error) {
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial(GDAX_SOCKET_URL, nil)
	if err != nil {
		return nil, err
	}

	subscribe := map[string]string{
		"type":      SUBSCRIBE_TYPE,
		"product_id": options.Product_id,
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		return nil, err
	}

	return &SocketClient{gdaxClient, wsConn}, nil
}

func (sc SocketClient) Read(message *exchange.Message) error {
	if err := sc.wsConn.ReadJSON(message); err != nil {
		return err
	}

	return nil
}
