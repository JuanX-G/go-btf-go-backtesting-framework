package backtester

type DataFeed interface {
	Next() (map[Symbol]Candle, bool)
	Reset()
}
