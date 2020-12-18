# Nintendo Switch Pro Controller

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bbayszczak/raspberrypi-go-drivers)](https://pkg.go.dev/github.com/raspberrypi-go-drivers/switchprocontroller)
![golangci-lint](https://github.com/raspberrypi-go-drivers/switchprocontroller/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/raspberrypi-go-drivers/switchprocontroller)](https://goreportcard.com/report/github.com/raspberrypi-go-drivers/switchprocontroller)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This drivers allows to read inputs from a Nintendo Switch Pro Controller
connected using Bluetooth

## Requirements

- a Nintendo Switch Pro Controller connected to a Raspberry Pi using Bluetooth

## Documentation

For full documentation, please visit [![PkgGoDev](https://pkg.go.dev/badge/github.com/bbayszczak/raspberrypi-go-drivers)](https://pkg.go.dev/github.com/raspberrypi-go-drivers/switchprocontroller)

## Quick start

```go
import (
	"fmt"

	"github.com/raspberrypi-go-drivers/switchprocontroller"
)

func main() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		select {
		case event := <-controller.Events:
			if event.Button != nil {
				fmt.Printf("%s:%d\n", event.Button.Name, event.Button.State)
			} else if event.Stick != nil {
				fmt.Printf("%s - Y:%f X:%f\n", event.Stick.Name, event.Stick.Y, event.Stick.X)
			}
		case <-time.After(60 * time.Second):
			fmt.Println("timeout")
			return
		}
	}
}
```
