# ITCH

This package implements the NASDAQ ITCH 5.0 protocol

https://www.nasdaqtrader.com/content/technicalsupport/specifications/dataproducts/NQTVITCHspecification.pdf

## Performance

A 12GB sample file `01302020.NASDAQ_ITCH50` was used, available from Nasdaq's FTP server.
The file contains 423,285,709 ITCH messages. Test was run on i5-8600k, 32GB RAM, SSD.

Loading the file into memory and then parsing using `itch.ParseBytes` took ~1.47s (286,419,921.07 messages/s), not including the time to read the file (which takes about ~5s itself)

Parsing from file using `itch.ParseFile` took ~49m. This is heavily bottlenecked by IO and should only be used
if you can't fit the whole file into memory. This can probably be improved in the future.

If you are parsing only the first few messages (e.g. <~100,000) of the file using `Configuration.MaxMessages`, then `itch.ParseFile` will be quicker. Parsing 100,000 messages took only 0.764s in this case.