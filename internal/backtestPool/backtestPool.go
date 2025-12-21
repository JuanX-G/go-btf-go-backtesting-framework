package backtestpool

import (
	"fmt"
	"sync"

	tester "go-backtesting-framework/internal/backtester"
)

func loopWarpper(wg *sync.WaitGroup, s tester.Strategy ) {
	defer wg.Done()
	tester.TestLoopStart(s)
}

type TestPool struct {
	strategies []tester.Strategy
}

func(t TestPool) Run() {
	var wg sync.WaitGroup
	for _, s := range t.strategies {
		wg.Add(1)
		go loopWarpper(&wg, s)
	}
	wg.Wait()
	fmt.Println("All testes finished")
}


