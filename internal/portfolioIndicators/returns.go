package portfolioIndicators

import (
	"math"
)


func(p *PortfolioData) NetProfit(cash float64) float64 {
	return cash - p.StartingCash
}

func(p *PortfolioData) PercentageReturn(netProfit float64) float64 {
	return netProfit / p.StartingCash
}

func(p *PortfolioData) CAGR(cash, years  float64) float64 {
	return math.Pow((cash / p.StartingCash), (1/years)) - 1
}

func(p *PortfolioData) WinRate() float64 {
	return float64(p.WinningTrades) / float64(p.Trades)
}

