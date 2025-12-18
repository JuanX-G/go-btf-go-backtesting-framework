package testDataGen

import (
	"math/rand/v2"
	"strconv"
	tester "go-backtesting-framework/internal/backtester"
	"fmt"
)
func GenerateRandomStringOfInts(slen int) string {
	outString := ""
	for i := 0; i < slen; i++ {
		n := rand.Uint32N(9)
		sn := strconv.Itoa(int(n))
		outString = fmt.Sprint(outString, sn)
	}
	return outString
}

func GenerateMarketData(lenOfArr, lenOfMap int) []tester.MarketData {
	var retArray []tester.MarketData
	var symArray []tester.Symbol
	for i := 0; i < lenOfMap; i++ {
		sym := tester.Symbol{Name: GenerateRandomStringOfInts(32)}
		symArray = append(symArray, sym)
	}

	for i := 0; i < lenOfArr; i++ {
		sdMap := make(map[tester.Symbol]tester.Candle)
		for _, sym := range symArray {
			priceCoeff := rand.UintN(100) + 10
			price := rand.Float64() * float64(priceCoeff)
			sdMap[sym] = tester.Candle{Price: price} 
			dataToAppend := tester.MarketData{SymbolData: sdMap}
			retArray = append(retArray, dataToAppend)
		}
	}
	return retArray
}
