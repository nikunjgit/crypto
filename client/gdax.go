package client

import (
	exchange "github.com/preichenberger/go-coinbase-exchange"
)

type GDAXClient struct {
	Client *exchange.Client
}

func NewGDAXClient(client *exchange.Client) GDAXClient{
	return GDAXClient{client}
}
