# mixpanel

Implements mixpanel API in Go.

## Usage
```
% mixpanel export -h
Usage: mixpanel export [options]

  Exports mixpanel data.

Options:

  -from=yesterday Start date to extract events.
  -to=yesterday   End date to extract events.
  -format=json    Choose export format between json/csv.
  -event=E        Extract data for only event E.
  -out=STDOUT     Decides where to write the data.
```

All the client methods, take an ```io.Writer```, and write results to it.
This makes it flexible enough to write to a buffer, file or plain STDOUT.

## Documentation

Full documentation is available on [Godoc](https://godoc.org/github.com/cskksc/mixpanel).
