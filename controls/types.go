package controls

import (
	"fmt"

	"github.com/simulatedsimian/joystick"
)

type Mode uint8

type device struct {
	name  string
	group *deviceGroup

	controls      map[string]control
	buttonHandler map[uint8]control

	joyId        int
	joystick     joystick.Joystick
	buttonEvents chan ButtonEvent
}

func NewDevice(name string) *device {
	return &device{
		name:          name,
		controls:      make(map[string]control, 0),
		buttonHandler: make(map[uint8]control, 0),
		buttonEvents:  make(chan ButtonEvent, 0),
	}
}

func (d *device) Control(name string) control {
	return d.controls[name]
}

func (d *device) AddControl(name string, c control) {
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

	fmt.Printf("device %s has %d buttons\n", device.name, device.buttonCount())

	for jsId := 0; jsId < 10; jsId++ {
		js, err := joystick.Open(jsId)
		if err != nil {
			continue
		}
		fmt.Printf("Found device %d with Axes: %d, Buttons: %d\n", js.Name(), js.AxisCount(), js.ButtonCount())
		if js.ButtonCount() == device.buttonCount() {
			device.joyId = jsId
			device.joystick = js
			fmt.Printf("using buttoncount %d device for %s\n", js.ButtonCount(), device.name)

			go device.PollJoystick()
			return
		}
	}
}

func (dg *deviceGroup) ModeToggle(c control, on, off Mode) {
	c.ModeAction(MODE_ALL, func(ev ButtonEvent) {
		if ev.pressed {
			dg.mode = dg.mode&(0xFF^off) | on
		} else {
			dg.mode = dg.mode&(0xFF^on) | off
		}
		fmt.Printf("MODE: %b\n", dg.mode)
	})
}

type Action func(ButtonEvent)

func OnPress(a func()) Action {
	return func(ev ButtonEvent) {
		if ev.pressed {
			a()
		}
	}
}

func ToggleAction(on, off func()) Action {
	return func(ev ButtonEvent) {
		if ev.pressed {
			on()
		} else {
			off()
		}
	}
}

type control interface {
	ButtonIDs() []uint8
	Handle(ButtonEvent)
	ModeAction(Mode, Action)
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
	for mode, action := range c.actions {
		if c.parent.group.mode == 0 || c.parent.group.mode&mode > 0 {
			action(ev)
		}
	}
}

func (c *button) ModeAction(mode Mode, action Action) {
	fmt.Printf("adding action for mode %b\n", mode)
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
	fmt.Printf("Toggle is unimplemented %v\n", ev)
}
func (c *toggle) ModeAction(mode Mode, action Action) {
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
	fmt.Printf("Toggle3 is unimplemented %v\n", ev)
}
func (c *toggle3) ModeAction(mode Mode, action Action) {
	c.actions[mode] = action
}

func (c *toggle3) setParent(dg *device) {
	c.parent = dg
}

type fourWay struct {
	parent  *device
	up      uint8
	right   uint8
	down    uint8
	left    uint8
	depress uint8

	pressable bool

	actions map[Mode]Action
}

func FourWay(u, r, d, l uint8) *fourWay {
	return &fourWay{
		up:      u,
		right:   r,
		down:    d,
		left:    l,
		actions: make(map[Mode]Action, 0),
	}
}

func (c *fourWay) ButtonIDs() []uint8 {
	out := []uint8{c.up, c.right, c.down, c.left}
	if c.pressable {
		out = append(out, c.depress)
	}
	return out
}

func (c *fourWay) Handle(ev ButtonEvent) {
	fmt.Printf("FourWay is unimplemented %v\n", ev)
}

func (f *fourWay) Depress(b uint8) *fourWay {
	f.depress = b
	f.pressable = true

	return f
}

func (c *fourWay) ModeAction(mode Mode, action Action) {
	c.actions[mode] = action
}

func (c *fourWay) setParent(dg *device) {
	c.parent = dg
}
