package main

import (
	"fmt"
	"github.com/markcheno/go-talib"
	"math"
)

func main() {
	fmt.Println(talib.Sin([]float64{0, math.Pi / 2}))
}
