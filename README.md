# Gocognito

Gocognito calculates [cognitive complexities](https://www.sonarsource.com/docs/CognitiveComplexity.pdf) of functions in Go source files. Intutive difficulty for programmers is reflected in cognitive complexity.

## Installation

```shell
go install github.com/nrnrk/gocognito/cmd/gocognito@latest
```

## Usage

```
$ gocognito
gocognito: Calculate cognitive complexity of functions.

Usage: gocognito [-flag] [package]

The gocognito analysis reports functions whose complexity is over than the specified limit.
...
```

Examples:

```shell
gocognit .
gocognit main.go
gocognit -over 10 ./pkg/...
```

## Related project
- [Gocognit](https://github.com/uudashr/gocognit) where the codes are based on.