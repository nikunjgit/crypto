package event

import (
	exchange "github.com/preichenberger/go-coinbase-exchange"
	"encoding/json"
	"time"
)

type Message struct {
	Type          string  `json:"type"`
	ProductId     string  `json:"product_id"`
	TradeId       int     `json:"trade_id,number"`
	OrderId       string  `json:"order_id"`
	Sequence      int64   `json:"sequence,number"`
	Time          time.Time    `json:"time,string"`
	Size          float64 `json:"size,string"`
	Price         float64 `json:"price,string"`
}

type Messages []*Message

func (s Messages) Len() int {
	return len(s)
}
func (s Messages) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Messages) Less(i, j int) bool {
	return s[i].Time.Before(s[j].Time)
}

func ToMessage(message *exchange.Message) *Message {
	return &Message{
		Type: message.Type,
		ProductId: message.ProductId,
		TradeId: message.TradeId,
		OrderId: message.OrderId,
		Sequence: message.Sequence,
		Time: message.Time.Time(),
		Size: message.Size,
		Price: message.Price,
	}
}

func FromMessage(message *Message) *exchange.Message {
	return &exchange.Message{
		Type: message.Type,
		ProductId: message.ProductId,
		TradeId: message.TradeId,
		OrderId: message.OrderId,
		Sequence: message.Sequence,
		Time: exchange.Time(message.Time),
		Size: message.Size,
		Price: message.Price,
	}

}
type Storage interface {
	Get(start time.Time, end time.Time) (Messages, error)
	Set(*Message) error
}

type State struct {
	store Storage
}

func NewState(store Storage) *State {
	return &State{store}
}

func (m *State) OnData(message *exchange.Message) {
	if message.Type == "match" {
		m.store.Set(ToMessage(message))
	}
}

func (m *State) DataString() ([]byte, error) {
	data, err := m.store.Get(time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}
