package backtester


func TestLoopStart(strat Strategy) {
	strat.Initialize()
	defer strat.Shutdown()
	var s bool
	for s {
		strat.Eval()
		s = strat.BrokerNext()
	}
}

/* 
* - 'Initialize' should prepare maps etc. 
* - 'Shutdown' should be used for cleanup and shutdown for calculating end metrics and other data. 
*     Shutdown also, MUST call the shutdown of 'broker'.
* - 'Eval' is called every bar, there you can open and close positions based on your signals
* - 'Broker next' must call 'Next' on broker. It should return 'true' if testing should continue, 'false' to stop it, 
*     like when the end of data is reached
*/
type Strategy interface {
	Initialize()
	Shutdown()
	Eval()
	BrokerNext() bool
}

type Comissions struct {
	BuyComission float64
	SellComission float64
}

var POSSIBLE_ORDER_TYPES = []string{"Buy", "Sell", "MarketBuy", "MarketSell"}

type Position struct {
	Sym Symbol
	OpenPrice float64
	StopLossPrice float64
	TakeProfitPrice float64
	Size float64
	ID int
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

