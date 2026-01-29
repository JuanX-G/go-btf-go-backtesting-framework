package datafeed

import (
	"encoding/csv"
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
	"os"
	"strconv"
)

type BasicDatafeed struct {
	data map[tester.Symbol][]tester.Candle
	idx int
}

func NewDatafeedFromFiles(fileNames []string, symNames []string) (*BasicDatafeed, error) {
	datafeed := BasicDatafeed{idx: 0}
	if len(fileNames) != len(symNames) {
		return nil, fmt.Errorf("mismatched filenames length with symNames length")
	}
	for idx, fname := range fileNames {
		dt, err := LoadCandleData(fname)
		if err != nil {
			return nil, err
		}
		currSym := tester.Symbol{Name: symNames[idx]}
		datafeed.data[currSym] = dt
	}
	return &datafeed, nil
}

func(b *BasicDatafeed) Next() (map[tester.Symbol]tester.Candle, bool) {
	b.idx++
	retMap := make(map[tester.Symbol]tester.Candle)
	for sym, candles := range b.data {
		if b.idx >= len(candles) {
			return nil, false
		}
		retMap[sym] = candles[b.idx]
	}
	return retMap, true
}

func(b *BasicDatafeed) Reset() {
	b.idx = 0
}
func LoadCandleData(fileName string) ([]tester.Candle, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return []tester.Candle{}, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return []tester.Candle{}, err
	}
	var retArr []tester.Candle
	for i, rec := range records { 
		if i < 3 {
			continue
		}
		if len(rec) == 6 {
			price, err := strconv.ParseFloat(rec[1], 64)
			if err != nil {
				return []tester.Candle{}, err
			}
			closeP, err := strconv.ParseFloat(rec[2], 64)
			if err != nil {
				return []tester.Candle{}, err
			}
			high, err := strconv.ParseFloat(rec[3], 64)
			if err != nil {
				return []tester.Candle{}, err
			}
			low, err := strconv.ParseFloat(rec[4], 64)
			if err != nil {
				return []tester.Candle{}, err
			}
			vol, err := strconv.ParseFloat(rec[5], 64)
			if err != nil {
				return []tester.Candle{}, err
			}

			entry := tester.Candle{
				Price: price,
				Close: closeP,
				High: high,
				Low: low,
				Volume: vol,
			}
			retArr = append(retArr, entry)
		}
	}
	return retArr, nil
}
