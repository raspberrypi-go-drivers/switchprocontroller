# Nintendo Switch Pro Controller

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller)](https://pkg.go.dev/github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller)

This drivers allows to read inputs from a Nintendo Switch Pro Controller
connected using Bluetooth

## Requirements

- a Nintendo Switch Pro Controller connected using Bluetooth

## Documentation

For full documentation, please visit [![PkgGoDev](https://pkg.go.dev/badge/github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller)](https://pkg.go.dev/github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller)

## Quick start

```go
import (
	"fmt"

	"github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller"
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
