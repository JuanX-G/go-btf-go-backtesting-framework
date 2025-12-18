package indicators

import (
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
)

func SMA(dt []tester.Candle, period, currIdx int) (float64, error) {
	 if period - 1 > len(dt) {
		return 0, fmt.Errorf("period too big")
	}
	runningTotal := 0.0
	for i := 0; i <= period; i++ {
		runningTotal += dt[currIdx - i].Price
	}
	avg := runningTotal / float64(period)
	return avg, nil
}

func RawSMA(dt []float64,) (float64, error) {
	runningTotal := 0.0
	for _, v := range dt {
		runningTotal += v
	}
	avg := runningTotal / float64(len(dt))
	return avg, nil

}

func EMA(dt []tester.Candle, period, currIdx int) (float64, error){
	if period - 1 > len(dt) {
		return 0, fmt.Errorf("period too big")
	}
	if currIdx < period {
		return 0, fmt.Errorf("period too big")
	}
	multiplier := 2/(period + 1)
	var prevEma float64
	itrc := 0

	for i := currIdx - period; i <= currIdx; i++ {
		if itrc == 0 { 
		     startingPoint, err := SMA(dt, period, currIdx)
		     if err != nil {
		     	return 0, err
		     }
		     currPrice := dt[period].Price
		     prevEma = (currPrice - startingPoint) * float64(multiplier) + startingPoint
		}
		currPrice := dt[i].Price 
		prevEma = (currPrice - prevEma) * float64(multiplier) + prevEma
		itrc++
	}
	return prevEma, nil
}

func RawEMA(dt []float64) (float64, error) {
	multiplier := 2/(len(dt)+ 1)
	var prevEma float64
	itrc := 0
	for _, v := range dt {
		if itrc == 0 { 
		     startingPoint, err := RawSMA(dt)
		     if err != nil {
		     	return 0, err
		     }
		     currPrice := dt[0]
		     prevEma = (currPrice - startingPoint) * float64(multiplier) + startingPoint
		}
		currPrice := v
		prevEma = (currPrice - prevEma) * float64(multiplier) + prevEma
		itrc++
	}
	return prevEma, nil
}

/* Moving Average Divergance Convergance */
type MACDdata struct  {
	ShortEma float64
	LongEma float64
	MACDval float64
	Signal float64
}

func(m MACDdata) PrettyString() string {
	return fmt.Sprintf("Long EMA: %f\nShort EMA: %f\nSignal EMA: %f\nMACD: %f", m.LongEma, m.ShortEma, m.Signal, m.MACDval)
}

func MACD(dt []tester.Candle, shortPeriod, longPeriod, signalPeriod, currIdx int) (MACDdata, error) { 
	var macds []float64
	var shortEma float64
	var longEma float64
	var macd float64
	for i := signalPeriod; i >= 0; i-- {
		shortEma, err := EMA(dt, shortPeriod, currIdx - i)
		if err != nil {
			return MACDdata{}, err 
		}
		longEma, err = EMA(dt, longPeriod, currIdx - i)
		if err != nil {
			return MACDdata{}, err 
		}
		macd = shortEma - longEma
		macds = append(macds, macd)
	}
	signal, err := RawEMA(macds)
	if err != nil {
		return MACDdata{}, err 
	}
	return MACDdata{ShortEma: shortEma, LongEma: longEma, MACDval: macd, Signal: signal}, nil
}
