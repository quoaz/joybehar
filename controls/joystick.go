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

func (d *device) events(prev, next uint32, buttonCount, buttonIDOffset int) []ButtonEvent {
	out := []ButtonEvent{}
	mask := prev ^ next
	for i := uint8(0); i < uint8(buttonCount); i++ {
		if mask>>i&1 > 0 {
			event := ButtonEvent{d, i + uint8(buttonIDOffset), false, time.Now()}
			if next>>i&1 > 0 {
				event.pressed = true
			}
			out = append(out, event)
		}
	}
	return out
}

func (d *device) processEvents() {
	fmt.Printf("%s processing events\n", d.name)
	for v := range d.buttonEvents {
		fmt.Printf("handling: %v\n", v)
		fmt.Printf("handler: %v\n", d.buttonHandler[v.button])
		go d.buttonHandler[v.button].Handle(v)
	}
}

func (d *device) PollJoystick() {
	fmt.Printf("lets poll %s on joystick %d -- ", d.name, d.joyId)
	fmt.Printf("%d axes, %d buttons + %d pov buttons\n", d.joystick.AxisCount(), d.joystick.ButtonCount(), d.povButtonCount)

	go d.processEvents()

	var state joystick.State
	var povstate uint32 = 0
	for {
		time.Sleep(100 * time.Millisecond)
		newState, err := d.joystick.Read()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		//		fmt.Printf("Axes: %v\n", newState.AxisData)

		if state.Buttons != newState.Buttons {
			events := d.events(state.Buttons, newState.Buttons, d.joystick.ButtonCount(), 0)
			for _, v := range events {
				d.buttonEvents <- v
			}
		}

		if d.povButtonCount > 0 {
			var newpovstate uint32 = 0
			lon_axis := d.joystick.AxisCount() - 1
			lat_axis := lon_axis - 1
			if newState.AxisData[lat_axis] > 1<<10 {
				newpovstate |= 1 << 3
			}
			if newState.AxisData[lat_axis] < -1<<10 {
				newpovstate |= 1 << 2
			}
			if newState.AxisData[lon_axis] > 1<<10 {
				newpovstate |= 1 << 1
			}
			if newState.AxisData[lon_axis] < -1<<10 {
				newpovstate |= 1 << 0
			}

			if newpovstate != povstate {
				events := d.events(povstate, newpovstate, d.povButtonCount, d.joystick.ButtonCount())
				for _, v := range events {
					d.buttonEvents <- v
				}
			}
			povstate = newpovstate
		}

		state = newState
	}
}
