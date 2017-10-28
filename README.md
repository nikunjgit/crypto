Go Coinbase Exchange [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/preichenberger/go-coinbase-exchange) [![Build Status](https://travis-ci.org/preichenberger/go-coinbase-exchange.svg?branch=master)](https://travis-ci.org/preichenberger/go-coinbase-exchange)
========

## Summary

Go trader for [GDAX](https://www.gdax.com)

## Installation

```sh
go get github.com/nikunjgit/crypto
```

## Documentation
For full details on functionality, see [GoDoc](http://godoc.org/github.com/nikunjgit/crypto) documentation.

### Setup
How to create a client:

```go

import (
  "os"
  exchange "github.com/nikunjgit/crypto"
)

```

### Testing
To test with Coinbase's public sandbox set the following environment variables:
  - TEST_COINBASE_SECRET
  - TEST_COINBASE_KEY
  - TEST_COINBASE_PASSPHRASE

Then run `go test`
```sh
TEST_COINBASE_SECRET=secret TEST_COINBASE_KEY=key TEST_COINBASE_PASSPHRASE=passphrase go test
```
