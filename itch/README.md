# ITCH

This package implements the Nasdaq ITCH 5.0 protocol

https://www.nasdaqtrader.com/content/technicalsupport/specifications/dataproducts/NQTVITCHspecification.pdf

## Performance

A 12GB sample file `01302020.NASDAQ_ITCH50` was used, available from Nasdaq's FTP server. The file contains 423,285,709 ITCH messages. A test was run on i5-8600k, 32GB RAM, SSD.

Loading the whole file into memory and then parsing using `itch.ParseBytes` while filtering for only `S` message types takes ~2.4s (176,041,695.88 messages/s), not including the time to read the file into memory (which takes about ~5s itself).

The function `itch.ParseFile` creates a buffered reader using bufio. The parsing speed depends on the buffer size up to a certain point. A 1GB buffer takes ~8s to parse, whilst the default 4KB buffer takes ~18s.