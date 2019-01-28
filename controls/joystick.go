package controls

import (
	"fmt"
	"time"

	"github.com/simulatedsimian/joystick"
)

type ButtonEvent struct {
	device    *device
	button    uint8
	pressed   bool
	timestamp time.Time
}

func (r ButtonEvent) String() string {
	pressed := "released"
	if r.pressed {
		pressed = "pressed"
	}
	return fmt.Sprintf("device %d: %d axis %d button: button %d %s\n", r.device.joyId, r.device.joystick.AxisCount(), r.device.joystick.ButtonCount(), r.button, pressed)
}

func (d *device) events(prev, next joystick.State) []ButtonEvent {
	out := []ButtonEvent{}
	mask := prev.Buttons ^ next.Buttons
	for i := uint8(0); i < uint8(d.joystick.ButtonCount()); i++ {
		if mask>>i&1 > 0 {
			event := ButtonEvent{d, i, false, time.Now()}
			if next.Buttons>>i&1 > 0 {
				event.pressed = true
			}
			out = append(out, event)
		}
	}
	return out
}

func (d *device) processEvents() {
	for v := range d.buttonEvents {
		d.buttonHandler[v.button].Handle(v)
	}
}

func (d *device) PollJoystick() {
	fmt.Printf("lets poll %s on joystick %d\n", d.name, d.joyId)

	go d.processEvents()

	var state joystick.State
	for {
		time.Sleep(100 * time.Millisecond)
		newState, err := d.joystick.Read()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		if state.Buttons != newState.Buttons {
			events := d.events(state, newState)
			for _, v := range events {
				d.buttonEvents <- v
			}
		}

		state = newState
	}
}
