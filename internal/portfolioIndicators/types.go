package portfolioIndicators

type PortfolioData struct {
	StartingCash float64
	Trades int 
	WinningTrades int
	LosingTrades int
	PeakCash float64
	MaxDD float64
	EquitySnapshots []float64
}
