package portfolioIndicators

import (
	"math"
)

func(p *PortfolioData) SortinoRatio(riskFree float64) float64 {
	equity := p.EquitySnapshots
	if len(equity) < 2 {
		return 0
	}
	var returns []float64
	for i := 1; i < len(equity); i++ {
		r := (equity[i] - equity[i-1]) / equity[i-1]
		returns = append(returns, r)
	}
	var avg float64
	for _, r := range returns {
		avg += r
	}
	avg /= float64(len(returns))
	var downsideSum float64
	for _, r := range returns {
		if r < riskFree {
			d := r - riskFree
			downsideSum += d * d
		}
	}
	downDev := math.Sqrt(downsideSum / float64(len(returns)))
	if downDev == 0 {
		return math.Inf(1)
	}
	return (avg - riskFree) / downDev
}

