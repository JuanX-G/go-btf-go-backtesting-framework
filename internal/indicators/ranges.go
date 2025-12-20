package indicators

import (
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
	"math"
)
//TODO: add error types
func ATR(dt []tester.Candle, period, currIdx int) (float64, error) {
	if period - 1 > len(dt) {
		return 0, fmt.Errorf("period too big")
	}
	if currIdx < period {
		return 0, fmt.Errorf("period too big")
	}
	var firstTRs []float64
	for i := currIdx - period; i <= currIdx;i++ {
		if i == currIdx {
			var firstTRsSum float64
			for _, v := range firstTRs {
				firstTRsSum += v
			}
			firstATRsAvg := firstTRsSum / float64(period  - 1)
			ATR := (firstATRsAvg * float64(period - 1)) + firstTRs[len(firstTRs) - 1]
			ATR = ATR / float64(period)
			return ATR, nil 
		}
		currHL := dt[i].High - dt[i].Low
		currHprevC := math.Abs(dt[i].High - dt[i - 1].Close)
		currLprevC := math.Abs(dt[i].Low - dt[i - 1].Close)
		maxFinal := max(currHL, currHprevC, currLprevC)
		firstTRs = append(firstTRs, maxFinal)
	}
	return 0, fmt.Errorf("Error occured in ATR")
}
