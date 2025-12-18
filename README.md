# Go-btf a go based backtesting engine

## Architecture
    - A engine that has the event loop
    - The 'broker' struct containing: 
        - The information about the state of the backtest, like open positions etc.
        - The array of input data taken from a .csv file
