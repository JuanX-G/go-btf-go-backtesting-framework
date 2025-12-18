package backtester

import (
	"fmt"
	sliceUtils "go-backtesting-framework/internal/sliceUtils"
)

type SellTransactionInfo struct {
	price float64
	CashInflow float64
	ComissiosSum float64
}

type BuyTransactionInfo struct {
	CashOutflow float64
	ComissiosSum float64
}

type CurrentSimulationData struct {
	Orders []Order
	Portfolio map[Symbol][]Position
	Cash float64
}

type Hooks struct {
	OnPositionClosed func(Position, SellTransactionInfo)
	OnPositionOpened func(Position, BuyTransactionInfo)
	OnOrderSubmitted func(Order)
	OnNext func(CurrentSimulationData)
	OnError func(error)
}

type Broker struct {
	Orders []Order
	Portfolio map[Symbol][]Position
	Cash float64
	Commisions Comissions
	CurrentData map[Symbol]Candle
	Hooks Hooks
}

func (b Broker) GetPositions(sym Symbol) ([]Position, bool){
	retPos, f  := b.Portfolio[sym]
	if !f {
		return []Position{}, false
	}
	return retPos, true
}

func (b Broker) lookUpSymbolCandle(sym Symbol) (Candle, bool) {
	candle, f := b.CurrentData[sym]
	if !f {
		return Candle{}, false
	}
	return candle, true
}

func (b *Broker) SubmitOrder(order Order) error {
	var orderTypeValid bool
	for _, v := range POSSIBLE_ORDER_TYPES {
		if order.Type == v {
			orderTypeValid = true
		}
	}
	if orderTypeValid == false {
		return fmt.Errorf("error: invalid order type")
	}
	b.Orders = append(b.Orders, order)
	if b.Hooks.OnOrderSubmitted != nil {
		b.Hooks.OnOrderSubmitted(order)
	}
	return nil
}

func (b *Broker) openPosition(ord Order) error {
	pos := Position{ 
		Sym: ord.Sym,
		OpenPrice: ord.SubmissionPrice,
		TakeProfitPrice: ord.TakeProfitPrice,
		StopLossPrice: ord.StopLossPrice,
		Size: ord.Size,
	}
	cost := ord.Size * pos.OpenPrice
	comission := cost * b.Commisions.BuyComission
	cost = cost + comission
	if cost > b.Cash {
		return fmt.Errorf("not enough cash")
	}
	b.Cash -= cost
	b.Portfolio[ord.Sym] = append(b.Portfolio[ord.Sym], pos)
	b.Hooks.OnPositionOpened(pos, BuyTransactionInfo{
		ComissiosSum: comission,
		CashOutflow: cost,
	})
	return nil
}

func (b *Broker) closePosition(pos Position) {
	positions, ok := b.Portfolio[pos.Sym]
	if !ok {
		panic("Invalid symbol in a position")
	}

	candle, f := b.lookUpSymbolCandle(pos.Sym)
	if !f {
		panic("Invalid symbol in a position")
	}
	SellInflow := candle.Price * pos.Size
	comission := SellInflow * b.Commisions.SellComission
	b.Cash += SellInflow - comission

	for i, p := range positions {
		if p == pos {
			positions = sliceUtils.Remove(positions, i)
			break
		}
	}

	if len(positions) == 0 {
		delete(b.Portfolio, pos.Sym)
	} else {
		b.Portfolio[pos.Sym] = positions
	}
	if b.Hooks.OnPositionClosed != nil {
		b.Hooks.OnPositionClosed(pos, SellTransactionInfo{
			price: candle.Price,
			CashInflow: SellInflow,
			ComissiosSum: comission,
		})
	}
}

func(b *Broker) Next() {
	for s, positions := range b.Portfolio {
		for _, pos := range positions {
			if pos.StopLossPrice <= b.CurrentData[s].Price{
				b.closePosition(pos)
			} else if pos.TakeProfitPrice >= b.CurrentData[s].Price {
				b.closePosition(pos)
			}
		}
	}
	for ordIdx, ord := range b.Orders {
		if ord.BuyPrice >= b.CurrentData[ord.Sym].Price {
			err := b.openPosition(ord)
			if err != nil {
				b.Hooks.OnError(err)
			}
			b.Orders = sliceUtils.Remove(b.Orders, ordIdx)
		} else if ord.Type == "MarketBuy" {
			err := b.openPosition(ord)
			if err != nil {
				b.Hooks.OnError(err)
			}
			b.Orders = sliceUtils.Remove(b.Orders, ordIdx)
		} else {
			panic("invalid order type!")
		}
	}
	b.Hooks.OnNext(CurrentSimulationData{
		Portfolio: b.Portfolio,
		Cash: b.Cash,
		Orders: b.Orders,
	})
}

func(b *Broker) Shutdown() {
	for _, positions := range b.Portfolio {
		for _, pos := range positions {
				b.closePosition(pos)
		}
	}
}
