# go-finproto

[![Go Reference](https://pkg.go.dev/badge/github.com/markwinter/go-finproto.svg)](https://pkg.go.dev/github.com/markwinter/go-finproto)
[![Go Report Card](https://goreportcard.com/badge/github.com/markwinter/go-finproto)](https://goreportcard.com/report/github.com/markwinter/go-finproto)


go-finproto is a collection of finance-related protocols implemented in Golang.

Protocols are contained within top level packages in this repo, and examples of each protocol can be found in the `cmd` directory.

## Protocols

### Nasdaq ITCH 5.0

The `itch` directory contains an implementation of Nasdaq's ITCH 5.0 protocol.

### Nasdaq SoupBinTCP 4.1

The `soupbintcp` directory contains an implementation of Nasdaq's SoupBinTCP 4.1 protocol.