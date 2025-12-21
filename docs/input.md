# Input
## Formats
### Backtester 
It takes input in the form of a map of candles to their symbols for the current state. \
The backtester does not have knowledge of data before or after, that is up to the user (examples of how to do it are provided) \
    `// What 'broker' holds \
    map[Symbol]Candle \
    ` \
    In short, it is up to the user to manage the data and feed it into the broker correctly; that is usually done by incrementing an index into an array of candles that is mapped to each symbol.

### Datafeed
The datafeed accepts data in the following CSV format. \
    `date (or anything else)| close price | price high | price low | price open | volume`
Then parsing it into an array of candles \


