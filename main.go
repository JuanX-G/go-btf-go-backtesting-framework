package main

import (
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
	"os"
	"go-backtesting-framework/internal/datafeed"
	defaultHooks"go-backtesting-framework/internal/defaultLoggingHooks"
	"go-backtesting-framework/internal/indicators"
)

const STARTING_CASH = 10000.0
const DAYS_TO_TEST = 5
var msftSymbol = tester.Symbol{
	Name: "MSFT",
}

type myStrategy struct {
	broker tester.Broker
	data map[tester.Symbol][]tester.Candle
	dataLen int
	iteration int
	hooks tester.Hooks
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
	s.broker.Commisions.BuyComission = 0.012
	s.broker.Commisions.SellComission = 0.005
	s.data = myData
	curMap := make(map[tester.Symbol]tester.Candle)
	curMap[msftSymbol] = s.data[msftSymbol][0]
	s.broker.CurrentData = curMap
	s.dataLen = len(myData[msftSymbol])
}

func(s *myStrategy) Shutdown() {
	s.broker.Shutdown()
	fmt.Println("cash: ", s.broker.Cash, "after: ", s.iteration, " iterations")
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
		ema16, err := indicators.EMA(s.data[symName], 16, s.iteration)
		if err != nil {
			continue
		}
		ema120, err := indicators.EMA(s.data[symName], 120, s.iteration)
		if err != nil {
			continue
		}
		macd, err := indicators.MACD(s.data[symName], 26, 12, 9, s.iteration)
		if err != nil {
			continue
		}
		fmt.Println(macd.PrettyString())
		if ema16 > ema120 * 1.04 {
			positionsForSym, _ := s.broker.Portfolio[symName] 
			fmt.Println("LEN POSITIONS:", len(positionsForSym))
			if len(positionsForSym) > 2 {
				continue
			}
			ord := tester.Order {
				Sym: symName,
				Size: 0.75,
				StopLossPrice: candle.Price * 0.5,
				TakeProfitPrice: candle.Price * 1.65,
				Type: "MarketBuy",
				BuyPrice: candle.Price,
				SubmissionPrice: candle.Price,
			}
			s.broker.SubmitOrder(ord)
		}
	}
}



func main() {
	msftCandles, err := datafeed.LoadCandleData("./data/3mo_1h_MSFT")
	if err != nil {
		panic(err)
	}
	dtMap := make(map[tester.Symbol][]tester.Candle)
	dtMap[msftSymbol] = msftCandles
	myData = dtMap 
	var strat myStrategy
	tester.TestLoopStart(&strat)
}
