package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ianmcmahon/joybehar/controls"
	kbd "github.com/micmonay/keybd_event"

	"github.com/simulatedsimian/joystick"
)

var throttle, stick *device

const (
	PRESSED  bool = true
	RELEASED bool = false
)

func dcsToggle(msg string) controls.Action {
	return dcsToggleWithVals(msg, "1", "0")
}

func dcsToggleWithVals(msg, on, off string) controls.Action {
	return controls.ToggleAction(dcsSend(msg, on), dcsSend(msg, off))
}

func dcsSend(msg, val string) func() {
	return func() {
		dcs.Send(msg, val)
	}
}

var dcs *dcsAgent

func main() {
	_, err := kbd.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	dcs = DCSAgent()

	warthog := controls.WarthogGroup()

	stick := warthog.Device("stick")

	warthog.ModeToggle(stick.Control("paddle"), controls.MODE_SHIFT, controls.MODE_NORM)

	//stick.Control("trigger1").ModeAction(controls.MODE_ALL, pulse(kb, kbd.VK_SPACE))
	//stick.Control("trigger2").ModeAction(controls.MODE_ALL, pulse(kb, kbd.VK_T))
	stick.Control("weaponrelease").ModeAction(controls.MODE_NORM, dcsToggle("WEAPON_RELEASE"))
	stick.Control("weaponrelease").ModeAction(controls.MODE_SHIFT, controls.OnPress(dcsSend("RWR_PWR", "TOGGLE")))

	dcs.Register(&StringOutput{Addr: 0x0000, MaxLength: 16, Action: func(_ uint16, s string) {
		fmt.Printf("Airplane: %s\n", s)
	}})
	dcs.Register(&IntegerOutput{Addr: 0x0408, Mask: 0xffff, Action: func(_, val uint16) {
		fmt.Printf("Altitude: %d' MSL\n", val)
	}})
	dcs.Register(&IntegerOutput{Addr: 0x765a, Mask: 0x2000, Action: func(_, val uint16) {
		fmt.Printf("Weapon Release: %d\n", val)
	}})
	go dcs.Receive()

	for {
		time.Sleep(1 * time.Second)
	}
}

type device struct {
	name     string
	joyId    int
	joystick joystick.Joystick

	rawChan     chan rawEvent
	monitorAxis []bool
	rawAxisChan chan rawAxisEvent

	mode           Mode
	buttonHandlers map[uint8]*buttonHandler
}

func Device(name string) *device {
	d := &device{
		name:           name,
		rawChan:        make(chan rawEvent, 0),
		rawAxisChan:    make(chan rawAxisEvent, 0),
		monitorAxis:    make([]bool, 10),
		mode:           Mode(0),
		buttonHandlers: make(map[uint8]*buttonHandler, 0),
	}

	go d.processEvents()
	go d.processAxisEvents()

	return d
}

func (d *device) processEvents() {
	for v := range d.rawChan {
		//fmt.Printf("got raw event: %v\n", v)

		handler, ok := d.buttonHandlers[v.button]
		if !ok {
			continue // unmapped
		}

		handler.Handle(d.mode, v)
	}
}

func (d *device) processAxisEvents() {
	for v := range d.rawAxisChan {
		fmt.Printf("got axis event: %v\n", v)
	}
}

type rawEvent struct {
	device    *device
	button    uint8
	pressed   bool
	timestamp time.Time
}

type rawAxisEvent struct {
	device    *device
	axes      []int
	timestamp time.Time
}

func (r rawEvent) String() string {
	pressed := "released"
	if r.pressed {
		pressed = "pressed"
	}
	return fmt.Sprintf("device %d: %d axis %d button: button %d %s\n", r.device.joyId, r.device.joystick.AxisCount(), r.device.joystick.ButtonCount(), r.button, pressed)
}

func (d *device) events(prev, next joystick.State) []rawEvent {
	out := []rawEvent{}
	mask := prev.Buttons ^ next.Buttons
	for i := uint8(0); i < uint8(d.joystick.ButtonCount()); i++ {
		if mask>>i&1 > 0 {
			event := rawEvent{d, i, false, time.Now()}
			if next.Buttons>>i&1 > 0 {
				event.pressed = true
			}
			out = append(out, event)
		}
	}
	return out
}

func pollDevice(id int) {
	js, err := joystick.Open(id)
	if err != nil {
		if err.Error() != fmt.Sprintf("Failed to read Joystick %d", id) {
			fmt.Printf("err: %v\n", err)
		}
		return
	}
	defer js.Close()
	fmt.Printf("Name: '%s'  Axes: %d, Buttons: %d\n", js.Name(), js.AxisCount(), js.ButtonCount())

	device := Device(js.Name())

	if js.AxisCount() == 7 && js.ButtonCount() == 32 {
		device = throttle
	}
	if js.AxisCount() == 4 && js.ButtonCount() == 19 {
		device = stick
	}
	device.joyId = id
	device.joystick = js

	pollAxes := false
	for _, mon := range device.monitorAxis {
		if mon {
			pollAxes = true
		}
	}

	var state joystick.State
	for {
		time.Sleep(100 * time.Millisecond)
		newState, err := js.Read()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		if pollAxes {
			for i := 0; i < js.AxisCount(); i++ {
				if !device.monitorAxis[i] {
					newState.AxisData[i] = 0
				}
			}
			fmt.Printf("%+v -> %+v\n", state, newState)
			if !reflect.DeepEqual(state.AxisData, newState.AxisData) {
				fmt.Printf("Axes: %#v\n", newState.AxisData)
			}
		}

		if state.Buttons != newState.Buttons {
			events := device.events(state, newState)
			for _, v := range events {
				device.rawChan <- v
			}
		}

		state = newState
	}
}
