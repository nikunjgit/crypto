package view

import (
	"time"
	"github.com/nikunjgit/crypto/event"
	"sort"
	"fmt"
)

type TimePeriod struct {
	Storage event.Storage
}

type Values struct {
	Values     []float64
	resolution time.Duration
}

type Series struct {
	startTime         time.Time
	consolidationFunc ConsolidationFunc
	Values            Values
}

func (t *TimePeriod) Get(tp time.Duration, resolution time.Duration, maxPoints int, consolidationFunc ConsolidationFunc) (*Series, error) {
	current := time.Now()
	prev := current.Add(-1 * tp * time.Duration(maxPoints))
	messages, err := t.Storage.Get(prev, current)
	if err != nil {
		return nil, err
	}

	fmt.Println("messages from storage", len(messages))
	dp := messages
	total := len(messages)
	if total > maxPoints {
		dp = messages[total - maxPoints:]
	}
	fmt.Println("datapoints", len(dp))

	series := t.CreateSeries(dp, resolution, consolidationFunc)
	return series, nil
}
func (t *TimePeriod) CreateSeries(datapoints event.Messages, resolution time.Duration, consolidationFunc ConsolidationFunc) *Series {
	sort.Sort(datapoints)
	values := make([]float64, 0, 10)
	if len(datapoints) == 0 {
		return nil
	}

	start := datapoints[0].Time
	for i := 0; i < len(datapoints); {
		messages := make(event.Messages, 0, 10)
		messages = append(messages, datapoints[i])
		duration := time.Duration(0)
		if i > 0 {
			duration = datapoints[i].Time.Sub(datapoints[i-1].Time)
		}
		j := i + 1
		for ; duration < resolution && j < len(datapoints); j++ {
			messages = append(messages, datapoints[j])
			duration = duration + datapoints[j].Time.Sub(datapoints[j-1].Time)
		}

		i = j
		if duration >= resolution {
			values = append(values, consolidationFunc(messages))
		}
	}

	return &Series{start, consolidationFunc, Values{values, resolution}}

}

// A ConsolidationFunc consolidates values at a given point in time.  It takes the current consolidated
// value, the new value to add to the consolidation, and a count of the number of values that have
// already been consolidated.
type ConsolidationFunc func(messages event.Messages) float64
