# Go Binance client

This application can be used to compare Withdrawal fees between Opendax and Binance platforms

### How to Build
```go build .```

### How to Run
#### Compare withdraw fees
```OPENDAX_BASE_URL=https://example.com BINANCE_API_KEY=*YOU_API_KEY* BINANCE_SECRET=*YOUR_API_SECRET* ./binance fees```

#### Compare Markets configuration
OPENDAX_API_KEY=*changeme* OPENDAX_API_SECRET=*changeme* OPENDAX_ENGINE_ID=4 ./binance markets
