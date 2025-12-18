package backtester


func TestLoopStart(strat Strategy) {
	strat.Initialize()
	defer strat.Shutdown()
	for {
		strat.Eval()
		strat.BrokerNext()
	}
}

type Strategy interface {
	Initialize()
	Shutdown()
	Eval()
	BrokerNext()
}

var POSSIBLE_ORDER_TYPES = []string{"Buy", "Sell", "MarketBuy", "MarketSell"}

type Position struct {
	Sym Symbol
	OpenPrice float64
	StopLossPrice float64
	TakeProfitPrice float64
	Size float64
}

type Symbol struct {
	Name string
}

type Order struct {
	Sym Symbol
	Size float64
	StopLossPrice float64
	TakeProfitPrice float64
	Type string
	BuyPrice  float64
	SubmissionPrice float64
}

type Comissions struct {
	BuyComission float64
	SellComission float64
}

type Candle struct {
	Price float64
	Close  float64
	High float64
	Low float64
	Volume float64
}

type MarketData struct {
	Candles []Candle
}

