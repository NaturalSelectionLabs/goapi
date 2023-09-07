# Overview

A lightweight opinionated router to bridge function call and http api.
It helps to standardize the way to define input and output of a service.

## Features

- Auto generate openapi doc from code reflection
- Use function signature to define input and output
- Type safe without struct tags

## Usage

Check the [examples folder](lib/examples/).

Read the tests for details.

## Benchmark

Without any optimization, goapi is about 10% slower than echo for the simplest usage.
This benchmark is only for avoiding drastic performance changes,
the real performance depends on the complexity of the service.

```text
go test -bench=. -benchmem -cpuprofile profile.out ./lib/bench
goos: darwin
goarch: arm64
pkg: github.com/NaturalSelectionLabs/goapi/lib/bench
Benchmark_goapi-12         35328             32968 ns/op            8412 B/op        113 allocs/op
Benchmark_echo-12          38246             30866 ns/op            6778 B/op         82 allocs/op
PASS
ok      github.com/NaturalSelectionLabs/goapi/lib/bench 3.283s
```
