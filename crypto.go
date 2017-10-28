package main

import (
	"fmt"
	"net/http"
	"os"

	exchange "github.com/preichenberger/go-coinbase-exchange"
	"github.com/nikunjgit/crypto/client"
	"github.com/nikunjgit/crypto/event"
	"github.com/nikunjgit/crypto/store"
	"time"
	"github.com/nikunjgit/crypto/indicator"
	"github.com/nikunjgit/crypto/view"
)

func dataHandler(state *event.State) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := state.DataString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func main() {
	gdaxClient := initExchange()
	socketClient, err := client.NewSocketClient(gdaxClient, &client.SubscribeOptions{
		Product_id: "ETH-USD",
	})
	if err != nil {
		panic(fmt.Errorf("unable to initialize socket client %v", err))
	}

	rClient, err := store.NewRedisClient(time.Hour)
	if err != nil {
		panic(fmt.Errorf("unable to initialize redis client %v", err))
	}

	bufferedStore := store.NewBufferedStorage(rClient, 5* time.Second, "data")
	state := event.NewState(bufferedStore)

	generator := event.NewEventGenerator(socketClient)
	generator.Register(state)
	generator.Start()

	timePeriod := view.TimePeriod{bufferedStore}
	stats := make([]indicator.Stat, 0, 1)
	stats = append(stats, indicator.NewMACD(12, 26,9 , timePeriod, 15 * time.Second))

	gen := indicator.Generator{Stats: stats}
	gen.Start(1 * time.Minute, 1 * time.Minute)
	fmt.Println("Starting server on port 8080")
	http.HandleFunc("/data", dataHandler(state))
	http.ListenAndServe(":8080", nil)
}

func initExchange() (*exchange.Client) {
	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")
	return exchange.NewClient(secret, key, passphrase)
}
