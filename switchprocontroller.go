package switchprocontroller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xcafed00d/joystick"
	log "github.com/sirupsen/logrus"
)

const (
	stickPeakValue float32       = 20000
	fetchDelta     time.Duration = 10 * time.Millisecond
)

// SwitchProController represent the physical controller
//
type SwitchProController struct {
	fetchDelta time.Duration
	// each time a new event is received, true is sent to this channel
	Events chan *Event `json:"-"`
	// Sticks list
	Sticks []*Stick
	// Buttons list
	Buttons []*Button
}

// Event is an event
type Event struct {
	Stick  *Stick
	Button *Button
}

// Stick represent a physical stick
//
// It contains a value (in %) for each axis (x and y), this value is contained between -100 and 100
//
// 0 is the default value when the stick is in the default position
type Stick struct {
	Name string
	X    float32
	Y    float32
	xMin float32
	xMax float32
	yMin float32
	yMax float32
}

// Button represent a physical button
//
// There's two possible State
// 1: button pressed
// 0: button released
type Button struct {
	Name  string
	State int
	code  uint32
}

// GetStick returns a pointer to a Stick instance with specified name
//
// err not nil if stick name not found
func (controller *SwitchProController) GetStick(name string) (*Stick, error) {
	for _, stick := range controller.Sticks {
		if stick.Name == name {
			return stick, nil
		}
	}
	return nil, fmt.Errorf("impossible to find stick")
}

// GetButton returns a pointer to a Button instance with specified name
//
// err not nil if button name not found
func (controller *SwitchProController) GetButton(name string) (*Button, error) {
	for _, button := range controller.Buttons {
		if button.Name == name {
			return button, nil
		}
	}
	return nil, fmt.Errorf("impossible to find button")
}

// GetButtonState returns the state of the button with specified name
//
// returns:
// 0: button released
// 1: button pressed
//
// err not nil if button name not found
func (controller *SwitchProController) GetButtonState(name string) (int, error) {
	for _, button := range controller.Buttons {
		if button.Name == name {
			return button.State, nil
		}
	}
	return 0, fmt.Errorf("impossible to find button")
}

func (controller *SwitchProController) updateStick(name string, x float32, y float32) {
	stick, err := controller.GetStick(name)
	valueChanged := false
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  name,
		}).Error("stick not found")
	}
	if x < stick.xMin {
		stick.xMin = x
		log.WithField("stick.xMin", x).Info("stick peak value changed")
	} else if x > stick.xMax {
		stick.xMax = x
		log.WithField("stick.xMax", x).Info("stick peak value changed")
	}
	if y < stick.yMin {
		stick.yMin = y
		log.WithField("stick.yMin", y).Info("stick peak value changed")
	} else if y > stick.yMax {
		stick.yMax = y
		log.WithField("stick.yMax", y).Info("stick peak value changed")
	}
	var newX float32
	var newY float32
	if x > 0 {
		newX = (100.0 * x) / stick.xMax
	} else if x < 0 {
		newX = (-100.0 * x) / stick.xMin
	} else if x == 0 {
		newX = 0
	}
	if newX != stick.X {
		stick.X = newX
		valueChanged = true
	}
	if y > 0 {
		newY = (-100.0 * y) / stick.yMax
	} else if y < 0 {
		newY = (100.0 * y) / stick.yMin
	} else if y == 0 {
		newY = 0
	}
	if newY != stick.Y {
		stick.Y = newY
		valueChanged = true
	}
	if valueChanged {
		controller.eventChangeStick(stick)
	}
}

func (controller *SwitchProController) updateSticks(axisData []int) {
	controller.updateStick("left", float32(axisData[0]), float32(axisData[1]))
	controller.updateStick("right", float32(axisData[2]), float32(axisData[3]))
	controller.updateStick("pad", float32(axisData[4]), float32(axisData[5]))
}

func (controller *SwitchProController) updateButtons(buttonsValue uint32) {
	for _, button := range controller.Buttons {
		previousState := button.State
		if buttonsValue >= button.code {
			button.State = 1
			buttonsValue = buttonsValue - button.code
		} else {
			button.State = 0
		}
		if button.State != previousState {
			controller.eventChangeButton(button)
		}
	}
}

func initStick(name string) *Stick {
	return &Stick{
		Name: name,
		X:    0,
		Y:    0,
		xMin: -stickPeakValue,
		xMax: stickPeakValue,
		yMin: -stickPeakValue,
		yMax: stickPeakValue,
	}
}

func initButton(name string, code uint32) *Button {
	return &Button{
		Name:  name,
		State: 0,
		code:  code,
	}
}

// NewSwitchProController creates and initialize a SwitchProController instance
func NewSwitchProController() *SwitchProController {
	log.Info("creating new SwitchProController")
	controller := SwitchProController{
		fetchDelta: fetchDelta,
		Events:     make(chan *Event, 100),
		Sticks: []*Stick{
			initStick("left"),
			initStick("right"),
			initStick("pad"),
		},
		Buttons: []*Button{
			initButton("capture", 8192),
			initButton("home", 4096),
			initButton("rs", 2048), // right stick
			initButton("ls", 1024), // left stick
			initButton("+", 512),
			initButton("-", 256),
			initButton("zr", 128),
			initButton("zl", 64),
			initButton("r", 32),
			initButton("l", 16),
			initButton("x", 8),
			initButton("y", 4),
			initButton("a", 2),
			initButton("b", 1),
		},
	}
	return &controller
}

// Display pprint the current controller status
func (controller *SwitchProController) Display() {
	marshalled, err := json.MarshalIndent(*controller, "", "  ")
	if err != nil {
		log.WithField("error", err).Error("impossible to marshal controller")
		fmt.Println(err)
	}
	fmt.Print(string(marshalled))
}

// StartListener starts listening for controller inputs and
// keep the controller instance up to date
func (controller *SwitchProController) StartListener(joystickID int) error {
	js, err := joystick.Open(joystickID)
	if err != nil {
		return err
	}
	go func() {
		defer js.Close()
		for {
			state, err := js.Read()
			if err != nil {
				panic(err)
			}
			controller.updateSticks(state.AxisData)
			controller.updateButtons(state.Buttons)
			time.Sleep(fetchDelta)
		}
	}()
	return nil
}

func (controller *SwitchProController) eventChangeStick(stick *Stick) {
	controller.Events <- &Event{Stick: stick}
}

func (controller *SwitchProController) eventChangeButton(button *Button) {
	controller.Events <- &Event{Button: button}
}
