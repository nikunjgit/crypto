package indicator

import (
	"time"
	"github.com/markcheno/go-talib"
	"github.com/nikunjgit/crypto/view"
	"github.com/nikunjgit/crypto/event"
	"fmt"
)

type MACD struct {
	optInFastPeriod   int
	optInSlowPeriod   int
	optInSignalPeriod int
	timeperiod        view.TimePeriod
	resolution        time.Duration
}

func NewMACD(optInFastPeriod int, optInSlowPeriod int, optInSignalPeriod int, timeperiod view.TimePeriod, resolution time.Duration) *MACD {
	return &MACD{optInFastPeriod, optInSlowPeriod, optInSignalPeriod, timeperiod, resolution}
}

func consolidate(messages event.Messages) float64 {
	var total float64
	for _, v := range messages {
		total = total + v.Price
	}

	return total / float64(len(messages))
}


func (macd *MACD) Calculate(tp time.Duration) (float64, error) {
	series, err := macd.timeperiod.Get(tp, macd.resolution, 100000, consolidate)
	if err != nil {
		return 0, err
	}
	points := series.Values.Values
	fmt.Println(points)
	if len(points) < macd.optInSlowPeriod {
		return 0, fmt.Errorf("too few points for macd %d", len(points))
	}

	_, _, hist := talib.Macd(points, macd.optInFastPeriod, macd.optInSlowPeriod, macd.optInSignalPeriod)
	fmt.Println(hist)
	return hist[len(hist)-1], nil
}
