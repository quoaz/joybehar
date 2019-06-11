package controls

import (
	"github.com/ianmcmahon/joybehar/alert"
	"github.com/ianmcmahon/joybehar/dcs"
	"github.com/simulatedsimian/joystick"
)

type Mode uint8

type device struct {
	name  string
	group *DeviceGroup

	controls      map[string]Control
	buttonHandler map[uint8]string

	povButtonCount int

	joyId        int
	joystick     joystick.Joystick
	buttonEvents chan ButtonEvent
}

func NewDevice(name string) *device {
	d := &device{
		name:          name,
		controls:      make(map[string]Control, 0),
		buttonHandler: make(map[uint8]string, 0),
		buttonEvents:  make(chan ButtonEvent, 0),
	}

	mapControls(name, d)

	return d
}

func (d *device) Control(name string) Control {
	return d.controls[name]
}

func (d *device) AddPOVControl(name string, c Control) {
	d.povButtonCount++
	d.AddControl(name, c)
}

func (d *device) AddControl(name string, c Control) {
	c.setParent(d.group)
	d.controls[name] = c
	for _, id := range c.ButtonIDs() {
		d.buttonHandler[id] = name
	}
}

func (d *device) HandleEvent(ev ButtonEvent) {
	ctlName := d.buttonHandler[ev.button]
	moduleMap := d.group.CurrentMap()
	if moduleMap == nil {
		alert.Sayf("selected module: '%s' no module map found", d.group.currentModule)
		return
	}
	if ctl, ok := moduleMap.controls[d.name][ctlName]; ok {
		ctl.Handle(ev)
	}
}

func (d *device) buttonCount() int {
	return len(d.buttonHandler) - d.povButtonCount
}

type Modal interface {
	Mode() Mode
	SetMode(Mode)
}

type Interceptor interface {
	Modal
	Intercept(msg dcs.DCSMsg) (dcs.DCSMsg, error)
}

type DeviceGroup struct {
	mode          Mode
	devices       map[string]*device
	moduleMaps    map[string]*moduleMap
	currentModule string
}

func NewDeviceGroup() *DeviceGroup {
	return &DeviceGroup{
		mode:       MODE_NORM,
		devices:    make(map[string]*device, 0),
		moduleMaps: make(map[string]*moduleMap, 0),
	}
}

func (dg *DeviceGroup) Device(name string) *device {
	return dg.devices[name]
}

func (dg *DeviceGroup) Mode() Mode {
	return dg.mode
}

func (dg *DeviceGroup) SetMode(mode Mode) {
	dg.mode = mode
}

func (dg *DeviceGroup) AddDevice(device *device) {
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
			//fmt.Printf("using buttoncount %d device for %s\n", js.ButtonCount(), device.name)

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
	dg      Modal
	on, off Mode
}

func (a modeAction) HandleEvent(_ Control, state State) {
	if state == STATE_HI {
		a.dg.SetMode(a.dg.Mode()&(0xFF^a.off) | a.on)
	} else {
		a.dg.SetMode(a.dg.Mode()&(0xFF^a.on) | a.off)
	}
}

type Value string

type Control interface {
	ButtonIDs() []uint8
	Handle(ButtonEvent)
	Action(Mode, Action)
	setParent(Modal)
}

type button struct {
	parent   Modal
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
		if c.parent.Mode() == 0 || c.parent.Mode()&mode > 0 {
			action.HandleEvent(c, state)
		}
	}
}

func (c *button) Action(mode Mode, action Action) {
	//fmt.Printf("adding action for mode %b\n", mode)
	c.actions[mode] = action
}

func (c *button) setParent(dg Modal) {
	c.parent = dg
}

type toggle struct {
	parent   Modal
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
		if c.parent.Mode() == 0 || c.parent.Mode()&mode > 0 {
			action.HandleEvent(c, state)
		}
	}
}

func (c *toggle) Action(mode Mode, action Action) {
	c.actions[mode] = action
}

func (c *toggle) setParent(dg Modal) {
	c.parent = dg
}

type toggle3 struct {
	parent  Modal
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
		if c.parent.Mode() == 0 || c.parent.Mode()&mode > 0 {
			action.HandleEvent(c, state)
		}
	}
}

func (c *toggle3) Action(mode Mode, action Action) {
	c.actions[mode] = action
}

func (c *toggle3) setParent(dg Modal) {
	c.parent = dg
}
