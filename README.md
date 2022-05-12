# Kid

Simple web framework written in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Tarocch1/kid.svg)](https://pkg.go.dev/github.com/Tarocch1/kid)

## Installation

```sh
go get -u github.com/Tarocch1/kid
```

## Quickstart

```go
package main

import (
    "github.com/Tarocch1/kid"
    "github.com/Tarocch1/kid/middlewares/recovery"
)

func main() {
    k := kid.New()
    k.Use(recovery.New())

    k.Get("/", func(c *kid.Ctx) error {
        return c.String("Hello, World 👋!")
    })

    k.Listen(":3000")
}
```
