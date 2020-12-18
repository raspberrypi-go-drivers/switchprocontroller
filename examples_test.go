package switchprocontroller_test

import (
	"fmt"
	"time"

	"github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller"
)

func Example_nonBlocking() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		// display button A state
		aState, _ := controller.GetButtonState("a")
		fmt.Printf("A:%d\n", aState)
		// display left stick position
		leftStick, _ := controller.GetStick("left")
		fmt.Printf("x:%f - y:%f\n", leftStick.X, leftStick.Y)
		time.Sleep(100 * time.Millisecond)
	}
	// Output:
	// A:0
	// x:0.000000 - y:0.000000
	// [...]
}

func Example_blocking() {
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

func ExampleNewSwitchProController() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
}

func ExampleSwitchProController_GetStick() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	leftStick, _ := controller.GetStick("left")
	fmt.Printf("x:%f - y:%f\n", leftStick.X, leftStick.Y)
	// Output: x:0.000000 - y:0.000000
}

func ExampleSwitchProController_StartListener() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
}

func ExampleSwitchProController_GetButton_released() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	aButton, _ := controller.GetButton("a")
	fmt.Printf("name:%s - state:%d\n", aButton.Name, aButton.State)
	// Output: name:a - state:0
}

func ExampleSwitchProController_GetButton_pressed() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	aButton, _ := controller.GetButton("a")
	fmt.Printf("name:%s - state:%d\n", aButton.Name, aButton.State)
	// Output: name:a - state:1
}

func ExampleSwitchProController_GetButtonState_released() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	aState, _ := controller.GetButtonState("a")
	fmt.Printf("state:%d\n", aState)
	// Output: state:0
}

func ExampleSwitchProController_GetButtonState_pressed() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	aState, _ := controller.GetButtonState("a")
	fmt.Printf("state:%d\n", aState)
	// Output: state:1
}
