package datafeed

import (
	"encoding/csv"
	tester "go-backtesting-framework/internal/backtester"
	"os"
	"strconv"
)

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
