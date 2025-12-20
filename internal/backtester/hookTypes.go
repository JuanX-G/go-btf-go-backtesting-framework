/* 
* Hooks allow users to define their own logging functions 
* independatly of a logger class or similar
*/
package backtester

/* Hooks for different events in the testing loop*/
type Hooks struct {
	OnPositionClosed func(Position, SellTransactionInfo)
	OnPositionOpened func(Position, BuyTransactionInfo)
	OnOrderSubmitted func(Order)
	OnNext func(CurrentSimulationData)
	OnError func(error)
}

type SellTransactionInfo struct {
	Price float64
	Size float64
	CashInflow float64
	ComissiosSum float64
}

type BuyTransactionInfo struct {
	Price float64
	Size float64
	CashOutflow float64
	ComissiosSum float64
}

type CurrentSimulationData struct {
	Orders []Order
	Portfolio map[Symbol][]Position
	Cash float64
}
