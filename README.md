# Overview

A lightweight opinionated router to bridge function call and http api.
It helps to standardize the way to define input and output of a service.

## Features

- Auto generate openapi doc from code reflection
- Use function signature to define input and output
- Type safe without struct tags

## Usage

Hello world example: [main.go](lib/examples/hello-world/main.go).

RPC style example: [main.go](lib/examples/add/main.go).

Check the [examples folder](lib/examples/).

Read the tests for details.

## Benchmark

Without any optimization, goapi is about 7% slower than echo for the simplest usage.
This benchmark is only for avoiding drastic performance changes,
the real performance depends on the complexity of the service.

```text
go test -bench=. -benchmem ./lib/bench
goos: darwin
goarch: arm64
pkg: github.com/NaturalSelectionLabs/goapi/lib/bench
Benchmark_goapi-12         34472             33856 ns/op            8448 B/op        114 allocs/op
Benchmark_echo-12          36729             31175 ns/op            6776 B/op         82 allocs/op
PASS
ok      github.com/NaturalSelectionLabs/goapi/lib/bench 5.711s
```
