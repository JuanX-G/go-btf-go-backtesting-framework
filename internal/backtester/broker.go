package backtester

import (
	"fmt"
	"log"
	"strconv"
)

func Remove[S ~[]E, E any](s S, i int) S {
	if i >= 0 && i < len(s) {
		s = append(s[:i], s[i+1:]...)
	}
	return s
}
type Broker struct {
	Orders []Order
	Portfolio map[Symbol][]Position
	Cash float64
	Commisions Comissions
	CurrentData MarketData
}

func (b Broker) GetPositions(sym Symbol) ([]Position, bool){
	retPos, f  := b.Portfolio[sym]
	if !f {
		return []Position{}, false
	}
	return retPos, true
}

func (b Broker) lookUpSymbolPrice(sym Symbol) (float64, bool) {
	price, f := b.CurrentData.SymbolData[sym]
	if !f {
		return 0, false
	}
	return price, true
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
	if cost > b.Cash {
		return fmt.Errorf("not enough cash")
	}
	b.Cash -= cost
	b.Portfolio[ord.Sym] = append(b.Portfolio[ord.Sym], pos)
	return nil
}

func (b *Broker) closePosition(pos Position) {
		delete(b.Portfolio, pos.Sym)
		currentPrice, f := b.lookUpSymbolPrice(pos.Sym)
		if !f {
			panic("Invalid symbol in a position")
		}
		SellInflow := currentPrice * pos.Size
		comission := SellInflow * b.Commisions.SellComission
		b.Cash += SellInflow - comission
}

func(b *Broker) Next() {
	for s, positions := range b.Portfolio {
		for _, pos := range positions {
			if pos.StopLossPrice <= b.CurrentData.SymbolData[s] {
				b.closePosition(pos)
				log.Println("Position: ", pos, " closed")
			} else if pos.TakeProfitPrice > b.CurrentData.SymbolData[s] {
				b.closePosition(pos)
				log.Println("Position: ", pos, " closed")
			}
		}
	}
	for ordIdx, ord := range b.Orders {
		if ord.BuyPrice >= b.CurrentData.SymbolData[ord.Sym] {
			err := b.openPosition(ord)
			if err != nil {
				subPriceStr := strconv.FormatFloat(ord.SubmissionPrice, 'f', 2, 64)
				log.Printf("not enoguh cash for order: \nSymbol: %s\nPrice: %s\n", ord.Sym.Name, subPriceStr)
			}
			log.Println("Position from order: ", ord, " opened")
			log.Println("order idx: ", ordIdx)
			b.Orders = Remove(b.Orders, ordIdx)
		} else if ord.Type == "MarketBuy" {
			err := b.openPosition(ord)
			if err != nil {
				subPriceStr := strconv.FormatFloat(ord.SubmissionPrice, 'f', 2, 64)
				log.Printf("not enoguh cash for order: \nSymbol: %s\nPrice: %s\n", ord.Sym.Name, subPriceStr)
			}
			log.Println("Position from order: ", ord, " opened")
			log.Println("order idx: ", ordIdx)
			b.Orders = Remove(b.Orders, ordIdx)
		} else {
			panic("invalid order type!")
		}
	}
}

func(b *Broker) Shutdown() {
	for _, positions := range b.Portfolio {
		for _, pos := range positions {
				b.closePosition(pos)
				log.Println("Position: ", pos, " closed")
		}
	}
}
