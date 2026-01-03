# Go-btf a go based backtesting engine

## This version 0.2
It is far from finished

## quick start
    - download the project
    - have a csv file with data of format
        - date (or anything else)| close price | price high | price low | price open | volume
    - Input the data to a datafeed using the datafeed module
    - Setup a strategy struct that fullfils the interface 'Strategy'
    - Get to testing
#### Examples are provided in examples 

## Architecture
    - A engine that has the event loop
    - The 'broker' struct containing: 
        - The information about the state of the backtest, like open positions etc.
       
