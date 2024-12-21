# go-finproto

[![Go Reference](https://pkg.go.dev/badge/github.com/markwinter/go-finproto.svg)](https://pkg.go.dev/github.com/markwinter/go-finproto)
[![Go Report Card](https://goreportcard.com/badge/github.com/markwinter/go-finproto)](https://goreportcard.com/report/github.com/markwinter/go-finproto)


go-finproto is a collection of finance-related protocols implemented in Golang.

Protocols are contained within top level packages in this repo, and examples of each protocol can be found in the `cmd` directory.

## Protocols

### Nasdaq ITCH 5.0

The `itch` directory contains a parser for Nasdaq's ITCH 5.0 protocol.

### Nasdaq SoupBinTCP 4.1

The `soupbintcp` directory contains a server and client implementation of Nasdaq's SoupBinTCP 4.1 protocol.


## Suggestions or Contributions

Have a suggestion for a protocol to be implemented? Open an issue and I'll take a look. If I think it will be fun, I'll implement it. I also welcome implementations from others if you create a pull request.
