# pulsefetch
This CLI program will fetch pulsescan API and returns token balances of address in the Pulsechain blockchain.

# Usage
`./pulsefetch <address>`

## Output
```
| Token Name                      |         Balance |
|---------------------------------|-----------------|
| PulseX                          |           5,680 |
| Communis                        | 164,604,880,983 |
| HEXFIREIO                       |           5,555 |
| Incentive                       |               8 |
| Hedron                          |           8,405 |
| HEX                             |           8,047 |
| Wrapped Ether from Ethereum     |               0 |
```

# Build
```
go mod init pulsefetch
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o pulsefetch
upx --best pulsefetch
```
