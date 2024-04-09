# Ethereum Parser

Service that observes ethereum blockchain and ethereum addresses.

## Requirement

1. [GoLang 1.22.1 or higher](https://go.dev/dl/)

## Running

Run the command below:

    go run ./cmd/app

This will run an HTTP Server on port `8080`. If this port is busy in your machine,
you can set up another port via flag.

Also, as requested, JSON RPC server is default hosted by `cloudflare-eth.com` but 
it is also possible to change it via flag if wanted. 

    go run ./cmd/app --port=<YOUR_PORT> --host=<ETH_RPC_HOST>

Example 


    go run ./cmd/app --port=8080 --host=https://cloudflare-eth.com

## Usage

Examples below are going to use port `8080`, but replace it if you decided to change it via flag when running the server.

### Getting current block

**Path**: `/currentBlock`

**Method**: GET

**Required query parameter**: none

**Required body request**: none

Example:

    curl --location 'http://localhost:8080/currentBlock'

### Subscribe address

**Path**: `/subscribe`

**Method**: POST

**Required query parameter**: none

**Required body request**: `{"address": "<ETHEREUM_HEX_ADDRESS>"}`

Example:

```shell
curl --location 'http://localhost:8080/subscribe' \
--header 'Content-Type: application/json' \
--data '{"address":"0x21a31ee1afc51d94c2efccaa2092ad1028285549"}'
```

### Getting address' transactions

**Path**: `/currentBlock`

**Method**: GET

**Required query parameter**: `address` (ethereum hex address)

**Required body request**: none

Example:

    curl --location 'http://localhost:8081/transactions?address=0x21a31ee1afc51d94c2efccaa2092ad1028285549'

## Tests

To run tests, run the command below

    go test -race ./...
