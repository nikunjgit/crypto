package indicator

import (
	"time"
	"fmt"
)

type Stat interface {
	Calculate(tp time.Duration) (float64, error)
}

type Generator struct {
	Stats []Stat
}

func (g *Generator) Start(duration time.Duration, interval time.Duration) (chan struct{}){
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				for _, stat := range g.Stats {
					val, err := stat.Calculate(duration)
					fmt.Println("Value", val, err)
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}
