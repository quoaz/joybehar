package controls

import (
	"fmt"

	"github.com/simulatedsimian/joystick"
)

type Mode uint8

type device struct {
	name  string
	group *deviceGroup

	controls      map[string]Control
	buttonHandler map[uint8]Control

	joyId        int
	joystick     joystick.Joystick
	buttonEvents chan ButtonEvent
}

func NewDevice(name string) *device {
	return &device{
		name:          name,
		controls:      make(map[string]Control, 0),
		buttonHandler: make(map[uint8]Control, 0),
		buttonEvents:  make(chan ButtonEvent, 0),
	}
}

func (d *device) Control(name string) Control {
	return d.controls[name]
}

func (d *device) AddControl(name string, c Control) {
	c.setParent(d)
	d.controls[name] = c
	for _, id := range c.ButtonIDs() {
		d.buttonHandler[id] = c
	}
}

func (d *device) buttonCount() int {
	cnt := 0
	for range d.buttonHandler {
		cnt++
	}
	return cnt
}

type deviceGroup struct {
	mode    Mode
	devices map[string]*device
}

func DeviceGroup() *deviceGroup {
	return &deviceGroup{
		mode:    MODE_NORM,
		devices: make(map[string]*device, 0),
	}
}

func (dg *deviceGroup) Device(name string) *device {
	return dg.devices[name]
}

func (dg *deviceGroup) AddDevice(device *device) {
	dg.devices[device.name] = device
	device.group = dg

	//fmt.Printf("device %s has %d buttons\n", device.name, device.buttonCount())

	for jsId := 0; jsId < 10; jsId++ {
		js, err := joystick.Open(jsId)
		if err != nil {
			continue
		}
		//fmt.Printf("Found device %d with Axes: %d, Buttons: %d\n", js.Name(), js.AxisCount(), js.ButtonCount())
		if js.ButtonCount() == device.buttonCount() {
			device.joyId = jsId
			device.joystick = js
			fmt.Printf("using buttoncount %d device for %s\n", js.ButtonCount(), device.name)

			go device.PollJoystick()
			return
		}
	}
}

type State uint8

const (
	STATE_OFF State = 0
	STATE_HI  State = 1
	STATE_ON  State = 1
	STATE_LOW State = 2
)

type Action interface {
	HandleEvent(control Control, state State)
}

type modeAction struct {
	dg      *deviceGroup
	on, off Mode
}

func (a modeAction) HandleEvent(_ Control, state State) {
	if state == STATE_HI {
		a.dg.mode = a.dg.mode&(0xFF^a.off) | a.on
	} else {
		a.dg.mode = a.dg.mode&(0xFF^a.on) | a.off
	}
	//fmt.Printf("MODE: %b\n", a.dg.mode)
}

func (dg *deviceGroup) ModeToggle(c Control, on, off Mode) {
	c.Action(MODE_ALL, modeAction{dg, on, off})
}

type Value string

type Control interface {
	ButtonIDs() []uint8
	Handle(ButtonEvent)
	Action(Mode, Action)
	setParent(*device)
}

type button struct {
	parent   *device
	buttonId uint8
	actions  map[Mode]Action
}

func Button(b uint8) *button {
	return &button{
		buttonId: b,
		actions:  make(map[Mode]Action, 0),
	}
}

func (c *button) ButtonIDs() []uint8 {
	return []uint8{c.buttonId}
}

func (c *button) Handle(ev ButtonEvent) {
	state := STATE_OFF
	if ev.pressed {
		state = STATE_ON
	}
	for mode, action := range c.actions {
		if c.parent.group.mode == 0 || c.parent.group.mode&mode > 0 {
			action.HandleEvent(c, state)
		}
	}
}

func (c *button) Action(mode Mode, action Action) {
	//fmt.Printf("adding action for mode %b\n", mode)
	c.actions[mode] = action
}

func (c *button) setParent(dg *device) {
	c.parent = dg
}

type toggle struct {
	parent   *device
	buttonId uint8
	actions  map[Mode]Action
}

func Toggle(b uint8) *toggle {
	return &toggle{
		buttonId: b,
		actions:  make(map[Mode]Action, 0),
	}
}

func (c *toggle) ButtonIDs() []uint8 {
	return []uint8{c.buttonId}
}

func (c *toggle) Handle(ev ButtonEvent) {
	state := STATE_OFF
	if ev.pressed {
		state = STATE_ON
	}
	for mode, action := range c.actions {
		if c.parent.group.mode == 0 || c.parent.group.mode&mode > 0 {
			action.HandleEvent(c, state)
		}
	}
}

func (c *toggle) Action(mode Mode, action Action) {
	c.actions[mode] = action
}

func (c *toggle) setParent(dg *device) {
	c.parent = dg
}

type toggle3 struct {
	parent  *device
	upId    uint8
	downId  uint8
	actions map[Mode]Action
}

func Toggle3(u, d uint8) *toggle3 {
	return &toggle3{
		upId:    u,
		downId:  d,
		actions: make(map[Mode]Action, 0),
	}
}

func (c *toggle3) ButtonIDs() []uint8 {
	return []uint8{c.upId, c.downId}
}

func (c *toggle3) Handle(ev ButtonEvent) {
	state := STATE_OFF
	if ev.pressed {
		switch ev.button {
		case c.upId:
			state = STATE_HI
		case c.downId:
			state = STATE_LOW
		}
	}
	for mode, action := range c.actions {
		if c.parent.group.mode == 0 || c.parent.group.mode&mode > 0 {
			action.HandleEvent(c, state)
		}
	}
}

func (c *toggle3) Action(mode Mode, action Action) {
	c.actions[mode] = action
}

func (c *toggle3) setParent(dg *device) {
	c.parent = dg
}
