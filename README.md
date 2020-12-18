# Nintendo Switch Pro Controller

[![Go Reference](https://pkg.go.dev/badge/github.com/raspberrypi-go-drivers/switchprocontroller.svg)](https://pkg.go.dev/github.com/raspberrypi-go-drivers/switchprocontroller)
![golangci-lint](https://github.com/raspberrypi-go-drivers/switchprocontroller/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/raspberrypi-go-drivers/switchprocontroller)](https://goreportcard.com/report/github.com/raspberrypi-go-drivers/switchprocontroller)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This drivers allows to read inputs from a Nintendo Switch Pro Controller
connected using Bluetooth

## Requirements

- a Nintendo Switch Pro Controller connected to a Raspberry Pi using Bluetooth

## Documentation

For full documentation, please visit [![Go Reference](https://pkg.go.dev/badge/github.com/raspberrypi-go-drivers/switchprocontroller.svg)](https://pkg.go.dev/github.com/raspberrypi-go-drivers/switchprocontroller)

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

## Raspberry Pi compatibility

This driver has has only been tested on an Raspberry Pi Zero WH using integrated bluetooth but may work well on other Raspberry Pi having integrated Bluetooth

## License

MIT License