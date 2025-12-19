package main

import (
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
	_ "go-backtesting-framework/internal/datafeed"
	defaultHooks "go-backtesting-framework/internal/defaultLoggingHooks"
	_ "go-backtesting-framework/internal/indicators"
	"os"
)


var msftSymbol = tester.Symbol{
	Name: "MSFT",
}

const STARTING_CASH = 10000.0

type myStrategy struct {
	broker tester.Broker
	data map[tester.Symbol][]tester.Candle
	dataLen int
	iteration int
	hooks tester.Hooks
	positionsOpened int 
}

var myData map[tester.Symbol][]tester.Candle

func(s *myStrategy) Initialize() {
	s.hooks.OnPositionOpened = defaultHooks.DefOnOpenHook
	s.hooks.OnPositionClosed = defaultHooks.DefOnClosedHook
	s.hooks.OnNext = defaultHooks.DefOnNextHook
	s.hooks.OnOrderSubmitted = defaultHooks.DefOnOrderHook
	s.hooks.OnError = defaultHooks.DefOnErrorHook
	s.iteration = 0
	s.broker.Portfolio = make(map[tester.Symbol][]tester.Position)
	s.broker.Cash = STARTING_CASH
	s.broker.Hooks = s.hooks
	s.broker.Commisions.BuyComission = 0.005
	s.broker.Commisions.SellComission = 0.0005
	s.data = myData
	curMap := make(map[tester.Symbol]tester.Candle)
	curMap[msftSymbol] = s.data[msftSymbol][0]
	s.broker.CurrentData = curMap
	s.dataLen = len(myData[msftSymbol])
}

func(s *myStrategy) Shutdown() {
	s.broker.Shutdown()
	fmt.Println("cash: ", s.broker.Cash, "after: ", s.iteration, " iterations", 
	"with: ", s.positionsOpened, " positions being opened in total")
	os.Exit(0)
}

func(s *myStrategy) BrokerNext() {
	s.iteration++
	if s.iteration >= s.dataLen {
		s.Shutdown()
	}
	curMap := make(map[tester.Symbol]tester.Candle)
	curMap[msftSymbol] = s.data[msftSymbol][s.iteration]
	s.broker.CurrentData = curMap
	s.broker.Next()
}


func(s *myStrategy) Eval() {
	currDtIndex := s.iteration
	if currDtIndex < 2 {
		return
	}
	for symName, candle := range s.broker.CurrentData {
			ord := tester.Order {
				Sym: symName,
				Size: 1,
				StopLossPrice: candle.Price * 0.7,
				TakeProfitPrice: candle.Price * 1.5,
				Type: "MarketBuy",
				BuyPrice: candle.Price,
				SubmissionPrice: candle.Price,
			}
			s.broker.SubmitOrder(ord)
		}
	
}



func main() {
	strategy := myStrategy{}
	tester.TestLoopStart(&strategy)
}
