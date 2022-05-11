# Kid

Simple web framework written in Go

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
