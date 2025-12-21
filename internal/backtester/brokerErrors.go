package backtester

import "fmt"


/*	 		Error types for backtesting 		*/
type InvalidOrderTypeError struct {
	Type string
}
func(i InvalidOrderTypeError) Error() string {
	return fmt.Sprintf("Order type: %s is invalid", i.Type)
}

type NotEnoughCashError struct {
	Cost float64
}
func(n NotEnoughCashError) Error() string {
	return fmt.Sprintf("Not enough cash to pay: %f", n.Cost)
}

type InvalidSymbolError struct {
	SymbolGiven Symbol
}
func(i InvalidSymbolError) Error() string {
	return fmt.Sprintf("Symbol: %s is invalid", i.SymbolGiven.Name)
}
