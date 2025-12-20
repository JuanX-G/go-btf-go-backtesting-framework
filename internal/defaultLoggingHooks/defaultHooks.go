/* Default hooks for logging diffrent events */
package loggingHooks 

import (
	"fmt"
	tester "go-backtesting-framework/internal/backtester"
)

func DefOnOpenHook(pos tester.Position, tInfo tester.BuyTransactionInfo) {
	fmt.Printf("Position on symbol: %s opened\nOf size: %f \nAt price: %f\ntransaction sum is: %f\nComission: %f\n", pos.Sym.Name, pos.OpenPrice, pos.Size, tInfo.CashOutflow, tInfo.ComissiosSum)
}

func DefOnClosedHook(pos tester.Position, tInfo tester.SellTransactionInfo) {
	fmt.Printf("Position on symbol: %s closed\nOf size: %f \nAt price: %f\ntransaction sum is: %f\nSell Price: %f\nComission: %f\n", pos.Sym.Name, pos.Size, pos.OpenPrice, tInfo.CashInflow, tInfo.Price, tInfo.ComissiosSum)
}

func DefOnNextHook(currDt tester.CurrentSimulationData) {
	fmt.Printf("Current state:\nCash: %f\nCurrent positions: %#v\nCurrent orders: %#v\n", currDt.Cash, currDt.Portfolio, currDt.Orders)
}

func DefOnOrderHook(ord tester.Order) {
	fmt.Printf("Order for symbol: %s\nof size: %f\nwith buy price: %f\nAnd submission price of: %f\n", ord.Sym.Name, ord.Size, ord.BuyPrice, ord.SubmissionPrice)
}

func DefOnErrorHook(err error) {
	fmt.Printf("Error: %s occured\n", err.Error())
}
