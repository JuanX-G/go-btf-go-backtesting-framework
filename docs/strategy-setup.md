# Strategy-setup
## Strategy Object
A strategy object must fulfill the interface 'strategy'.
#### The 'eval' Method 
This is where you keep the "buisness"/"entry" logic. \
You There you call the 'submitOrder' method of the broker (which you should keep in strategy)

#### The 'BrokerNext' Method 
Here you feed new data into the broker by setting its 'CurrentData' field \
You also check if you have reached end-of-data, if so you return false to stop the testing loop \ 
Otherwise true should be returned to continue to the next "bar"

#### The 'Initialize' Method
This is were you initialize maps etc. 

#### The 'Shutdown' Method
Here you cleanup everything that needs to be cleanup, specifically, here you should log the final data into stdout or a file 





