/*
* 	The 'Broker' type is the main part of the engine, keep track of most things
*	state like positions, orderds, and more; also managing the metrics of the portfolio
 */
package backtester

import (
	"fmt"
	"go-backtesting-framework/internal/portfolioIndicators"
	sliceUtils "go-backtesting-framework/internal/sliceUtils"
	"math"
)

type Broker struct {
	Orders []Order
	Portfolio map[Symbol][]Position
	Cash float64
	Commisions Comissions
	CurrentData map[Symbol]Candle // Maps any sumbol the its current candle
	Hooks Hooks // Struct of function pointers to be called upon certain events
	nextPositionID int
	PortfolioData portfolioIndicators.PortfolioData
}

/* Look up all open positions for a given symbol  */
func (b Broker) GetPositions(sym Symbol) ([]Position, bool){
	retPos, f  := b.Portfolio[sym]
	if !f {
		return []Position{}, false
	}
	return retPos, true
}

/* Look up all candles for a given symbol  */
func (b Broker) lookUpSymbolCandle(sym Symbol) (Candle, bool) {
	candle, f := b.CurrentData[sym]
	if !f {
		return Candle{}, false
	}
	return candle, true
}
/* User facing function, used to send orders from the 'Eval' function */
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
	b.nextPositionID++
	pos := Position{ 
		Sym: ord.Sym,
		TakeProfitPrice: ord.TakeProfitPrice,
		StopLossPrice: ord.StopLossPrice,
		ID: b.nextPositionID,
	}
	price := b.CurrentData[ord.Sym].Price
	cost := 0.0
	if ord.Type == "MarketBuy" || ord.Type == "Buy" {
		pos.Size = ord.Size
	} else {
		pos.Size = -ord.Size
	}
	cost = math.Abs(ord.Size) * price

	comission := 0.0

	if pos.Size > 0 {
		comission = cost * b.Commisions.BuyComission
	} else {
		comission = cost * b.Commisions.SellComission
	}
	cost = cost + comission
	pos.OpenPrice = price
	if cost > b.Cash {
		return NotEnoughCashError{Cost: cost}
	}

	b.Cash -= cost
	b.Portfolio[ord.Sym] = append(b.Portfolio[ord.Sym], pos)
	if b.Hooks.OnPositionOpened == nil {
		return nil
	}
	b.Hooks.OnPositionOpened(pos, BuyTransactionInfo{
		ComissiosSum: comission,
		CashOutflow: cost,
	})
	return nil
}

func (b *Broker) closePosition(pos Position) {
	positions, ok := b.Portfolio[pos.Sym]
	if !ok {
		if b.Hooks.OnError != nil {
			b.Hooks.OnError(InvalidSymbolError{SymbolGiven: pos.Sym})
		}
		return
	}

	candle, f := b.lookUpSymbolCandle(pos.Sym)
	if !f {
		if b.Hooks.OnError != nil {
			b.Hooks.OnError(InvalidSymbolError{SymbolGiven: pos.Sym})
		}
		return
	}
	SellInflow := candle.Price * math.Abs(pos.Size)
	comission := 0.0
	if pos.Size > 0 {
		comission = SellInflow * b.Commisions.SellComission
	} else {
		comission = SellInflow * b.Commisions.BuyComission
	}
	b.Cash += SellInflow - comission

	for i, p := range positions {
		if p.ID == pos.ID {
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
			Price: candle.Price,
			CashInflow: SellInflow,
			ComissiosSum: comission,
		})
	}

	entryValue := pos.OpenPrice * math.Abs(pos.Size)
	exitValue  := candle.Price * math.Abs(pos.Size)

	sgn := 0.0
	if pos.Size > 0 {
		sgn = 1
	} else {
		sgn = -1
	}
	pnl := (exitValue - entryValue) * sgn
	pnl -= (comission)

	if pnl > 0 { 
		b.PortfolioData.WinningTrades++
	} else {
		b.PortfolioData.LosingTrades++
	}
}

func(b *Broker) Next() {
	equity := 0.0
	for s, positions := range b.Portfolio {
		positions := append([]Position(nil), positions...)
		for _, pos := range positions {
			isLong := pos.Size > 0
			//TODO: fix 0 tp/sl
			price := b.CurrentData[s].Price

			if isLong {
				if price <= pos.StopLossPrice || price >= pos.TakeProfitPrice {
					b.closePosition(pos)
				}
			} else {
				if price >= pos.StopLossPrice || price <= pos.TakeProfitPrice {
					b.closePosition(pos)
				}
			}
		}
		for _, pos := range positions {
			candle, f := b.lookUpSymbolCandle(pos.Sym)
			if !f {
				if b.Hooks.OnError != nil {
					b.Hooks.OnError(InvalidSymbolError{SymbolGiven: pos.Sym})
				}
			}
			//TODO decide if we want negatives here
			equity += pos.Size * candle.Price
		}
	}
	for i := len(b.Orders)-1; i >= 0; i-- {
		ord := b.Orders[i]
		executedShort := false
		executedLong := false
		if ord.Type == "MarketBuy" {
			executedLong = true
		} else if ord.Type == "Buy" && ord.BuyPrice >= b.CurrentData[ord.Sym].Price {
			executedLong = true
		} else if ord.Type == "MarketSell"   {
			executedShort = true
		} else if ord.Type == "Sell" && ord.BuyPrice <= b.CurrentData[ord.Sym].Price  {
			executedShort = true
		}

		if executedLong {
			if err := b.openPosition(ord); err != nil {
				if b.Hooks.OnError != nil {
					b.Hooks.OnError(err)
				}
			}
			b.Orders = sliceUtils.Remove(b.Orders, i)
		} else if executedShort {
			if err := b.openPosition(ord); err != nil {
				if b.Hooks.OnError != nil {
					b.Hooks.OnError(err)
				}
			}
			b.Orders = sliceUtils.Remove(b.Orders, i)
		}
	}
	equity += b.Cash
	if b.Hooks.OnNext != nil {
		b.Hooks.OnNext(CurrentSimulationData{
			Portfolio: b.Portfolio,
			Cash: b.Cash,
			Orders: b.Orders,
		})
	}
	b.PortfolioData.PeakEquity = max(b.PortfolioData.PeakEquity, equity)
	drawdown := (b.PortfolioData.PeakEquity - equity) / b.PortfolioData.PeakEquity
	b.PortfolioData.MaxDD  = max(drawdown, b.PortfolioData.MaxDD)
}

func(b *Broker) Shutdown() {
	for _, v := range b.Portfolio {
		positions := append([]Position(nil), v...)
		for _, pos := range positions {
			b.closePosition(pos)
		}
	}
}
