package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ianmcmahon/joybehar/controls"

	"github.com/simulatedsimian/joystick"
)

var throttle, stick *device

type _dcsAction struct {
	msg []string
}

func dcsAction(msgs ...string) _dcsAction {
	return _dcsAction{msgs}
}

func (a _dcsAction) HandleEvent(control controls.Control, state controls.State) {
	for _, msg := range a.msg {
		dcs.Send(msg, string(control.Value(state)))
	}
}

type _dcsToggle struct {
	msg []string
}

func dcsToggle(msgs ...string) _dcsToggle {
	return _dcsToggle{msgs}
}

func (a _dcsToggle) HandleEvent(_ controls.Control, state controls.State) {
	if state == controls.STATE_ON {
		for _, msg := range a.msg {
			dcs.Send(msg, "TOGGLE")
		}
	}
}

type _keyAction struct {
	keys []int
}

func keyAction(keys ...int) _keyAction {
	return _keyAction{keys}
}

func (a _keyAction) HandleEvent(control controls.Control, state controls.State) {
	switch state {
	case controls.STATE_ON:
		for _, k := range a.keys {
			downKey(k)
		}
	case controls.STATE_OFF:
		for _, k := range a.keys {
			upKey(k)
		}
	}
}

type _keyPulse struct {
	keys []int
}

func keyPulse(keys ...int) _keyPulse {
	return _keyPulse{keys}
}

func (a _keyPulse) HandleEvent(control controls.Control, state controls.State) {
	if state == controls.STATE_ON {
		for _, k := range a.keys {
			downKey(k)
		}
		time.Sleep(50 * time.Millisecond)
		for _, k := range a.keys {
			upKey(k)
		}
	}
}

var dcs *dcsAgent

func main() {
	dcs = DCSAgent()

	warthog := controls.WarthogGroup()

	stick := warthog.Device("stick")

	warthog.ModeToggle(stick.Control("paddle"), controls.MODE_SHIFT, controls.MODE_NORM)

	stick.Control("trigger1").Action(controls.MODE_ALL, keyAction(K_SPACE))
	stick.Control("trigger2").Action(controls.MODE_ALL, keyAction(K_T))

	stick.Control("weaponrelease").Action(controls.MODE_NORM, dcsAction("WEAPON_RELEASE"))
	stick.Control("weaponrelease").Action(controls.MODE_SHIFT, dcsToggle("RWR_PWR"))

	stick.Control("nws").Action(controls.MODE_ALL, dcsAction("NWS"))

	stick.Control("index").Action(controls.MODE_NORM, dcsAction("FL_CHAFF_BT"))
	stick.Control("index").Action(controls.MODE_SHIFT, keyPulse(0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0x5B, 0x5C))

	stick.Control("tms_up").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS0", "ARMPOS6"))
	stick.Control("tms_dn").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS3"))
	stick.Control("tms_rt").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS2", "ARMPOS4"))
	stick.Control("tms_lt").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS1", "ARMPOS5"))

	/*
		stick.Control("cms_up").Action(controls.MODE_NORM, keyPulse("dogfight guns"))
		stick.Control("cms_dn").Action(controls.MODE_NORM, keyPulse("dogfight missile"))
		stick.Control("cms_dp").Action(controls.MODE_NORM, keyPulse("dogfight resume"))
		stick.Control("cms_dp").Action(controls.MODE_SHIFT, keyPulse("escape (long press should be map)"))
	*/

	throttle := warthog.Device("throttle")

	//throttle.Control("mic_up").Action(controls.MODE_NORM, keyPulse("vr recenter"))
	//throttle.Control("mic_dn").Action(controls.MODE_NORM, keyPulse("vr zoom"))

	throttle.Control("speedbrake").Action(controls.MODE_ALL, dcsAction("SPEED"))
	throttle.Control("boatswitch").Action(controls.MODE_ALL, dcsAction("A_FLAPS"))
	throttle.Control("flaps").Action(controls.MODE_ALL, dcsAction("FLAPS"))
	throttle.Control("eacarm").Action(controls.MODE_ALL, dcsAction("GEAR"))

	dcs.Register(&StringOutput{Addr: 0x0000, MaxLength: 16, Action: func(_ uint16, s string) {
		fmt.Printf("Airplane: %s\n", s)
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
