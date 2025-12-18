package main

import (
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
	"os"
	"go-backtesting-framework/internal/datafeed"
)

const STARTING_CASH = 10000.0
const DAYS_TO_TEST = 5


type myStrategy struct {
	broker tester.Broker
	data []tester.MarketData
	dataLen int
	iteration int
	hooks tester.Hooks
}

func myOnOpenHook(pos tester.Position, tInfo tester.BuyTransactionInfo) {
	// fmt.Printf("Position on symbol: %s opened\nOf size: %f \nAt price: %f\ntransaction sum is: %f\nComission: %f\n", pos.Sym.Name, pos.OpenPrice, pos.Size, tInfo.CashOutflow, tInfo.ComissiosSum)
}
var myData []tester.MarketData;
func myOnClosedHook(pos tester.Position, tInfo tester.SellTransactionInfo) {
	// fmt.Printf("Position on symbol: %s opened\nOf size: %f \nAt price: %f\ntransaction sum is: %f\nComission: %f\n", pos.Sym.Name, pos.OpenPrice, pos.Size, tInfo.CashInflow, tInfo.ComissiosSum)
}
func myOnNextHook(currDt tester.CurrentSimulationData) {
	// fmt.Printf("Current state:\nCash: %f\nCurrent positions: %#v\nCurrent orders: %#v\n", currDt.Cash, currDt.Portfolio, currDt.Orders)
}

func myOnOrderHook(ord tester.Order) {
	// fmt.Printf("Order for symbol: %s\nof size: %f\nwith buy price: %f\nAnd submission price of: %f\n", ord.Sym.Name, ord.Size, ord.BuyPrice, ord.SubmissionPrice)
}

func myOnErrorHook(err error) {
	// fmt.Printf("Error: %s occured\n", err.Error())
}

func(s *myStrategy) Initialize() {
	s.hooks.OnPositionOpened = myOnOpenHook
	s.hooks.OnPositionClosed = myOnClosedHook
	s.hooks.OnNext = myOnNextHook
	s.hooks.OnOrderSubmitted = myOnOrderHook
	s.hooks.OnError = myOnErrorHook
	s.iteration = 0
	s.broker.Portfolio = make(map[tester.Symbol][]tester.Position)
	s.broker.Cash = STARTING_CASH
	s.broker.Hooks = s.hooks
	s.broker.Commisions.BuyComission = 0.012
	s.broker.Commisions.SellComission = 0.005
	s.data = myData
	s.broker.CurrentData = s.data[0]
	s.dataLen = len(myData)
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
	s.broker.CurrentData = s.data[s.iteration]
	s.broker.Next()
}

func(s *myStrategy) Eval() {
	currDtIndex := s.iteration
	if currDtIndex < 2 {
		return
	}
	for symName, candle := range s.broker.CurrentData.SymbolData {
		avg := (candle.Price + s.data[s.iteration - 1].SymbolData[symName].Price) / 2
		if candle.Price < 480 || avg < 500 {

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
	msftSymbol := tester.Symbol{
		Name: "MSFT",
	}
	mData := []tester.MarketData{}
	for _, v := range msftCandles {
		dataToAdd := tester.MarketData{
			SymbolData: make(map[tester.Symbol]tester.Candle),
		}
		dataToAdd.SymbolData[msftSymbol] = v 
		mData = append(mData, dataToAdd)
	}
	myData = mData
	fmt.Println("Hello!")
	var strat myStrategy
	tester.TestLoopStart(&strat)
}
