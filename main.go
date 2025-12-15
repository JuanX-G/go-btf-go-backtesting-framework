package main

import (
	"fmt"
	tester "go-trading-strategy/internal/backtester"
	"math/rand/v2"
	"os"
	"strconv"
)

const STARTING_CASH = 10000.0
const DAYS_TO_TEST = 5

func GenerateRandomStringOfInts(slen int) string {
	outString := ""
	for i := 0; i < slen; i++ {
		n := rand.Uint32N(9)
		sn := strconv.Itoa(int(n))
		outString = fmt.Sprint(outString, sn)
	}
	return outString
}

func generateMarketData(lenOfArr, lenOfMap int) []tester.MarketData {
	var retArray []tester.MarketData
	var symArray []tester.Symbol
	for i := 0; i < lenOfMap; i++ {
		sym := tester.Symbol{Name: GenerateRandomStringOfInts(32)}
		symArray = append(symArray, sym)
	}

	for i := 0; i < lenOfArr; i++ {
		sdMap := make(map[tester.Symbol]float64)
		for _, sym := range symArray {
			priceCoeff := rand.UintN(100) + 10
			price := rand.Float64() * float64(priceCoeff)
			sdMap[sym] = price 
			dataToAppend := tester.MarketData{SymbolData: sdMap}
			retArray = append(retArray, dataToAppend)
		}
	}
	return retArray
}

type myStrategy struct {
	broker tester.Broker
	data []tester.MarketData
	dataLen int
	iteration int
}

func(s *myStrategy) Initialize() {
	s.iteration = 0
	s.broker.Portfolio = make(map[tester.Symbol][]tester.Position)
	s.broker.Cash = STARTING_CASH
	s.data = generateMarketData(DAYS_TO_TEST, 2)
	s.broker.CurrentData = s.data[0]
	s.dataLen = DAYS_TO_TEST
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
	for symName, symPrice := range s.broker.CurrentData.SymbolData {
		if symPrice < 1000 {
			positionsForSym, _ := s.broker.Portfolio[symName] 
			fmt.Println("LEN POSITIONS:", len(positionsForSym))
			if len(positionsForSym) > 2 {
				continue
			}
			ord := tester.Order {
				Sym: symName,
				Size: 0.75,
				StopLossPrice: symPrice * 0.5,
				TakeProfitPrice: symPrice * 1.65,
				Type: "MarketBuy",
				BuyPrice: symPrice,
				SubmissionPrice: symPrice,
			}
			s.broker.SubmitOrder(ord)
		}
	}
}



func main() {
	fmt.Println("Hello!")
	var strat myStrategy
	tester.TestLoopStart(&strat)
}
