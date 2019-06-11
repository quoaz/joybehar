package controls

import (
	"fmt"

	"github.com/ianmcmahon/joybehar/alert"
	"github.com/ianmcmahon/joybehar/dcs"
)

// on error, drops the message entirely
type InterceptAction func(dcs.DCSMsg) (dcs.DCSMsg, error)

type ControlMap interface {
	AddControl(string, Control)
	AddPOVControl(string, Control)
}

type moduleMap struct {
	name        string
	deviceGroup Modal

	// map["device_name"]map["control_name"]Control
	controls   map[string]map[string]Control
	intercepts map[string]map[Mode]InterceptAction
}

func (m *moduleMap) ModeToggle(device, control string, on, off Mode) {
	c := m.Control(device, control)
	c.Action(MODE_ALL, modeAction{m.deviceGroup, on, off})
}

func (m *moduleMap) Control(device, control string) Control {
	if d, ok := m.controls[device]; ok {
		if c, ok := d[control]; ok {
			return c
		}
	}
	return nil
}

func (c *moduleMap) ReLabel(label string, mode Mode, alt string) {
	c.Intercept(label, mode, func(m dcs.DCSMsg) (dcs.DCSMsg, error) {
		m.Message = alt
		fmt.Printf("relabeling %s: %v\n", label, m)
		return m, nil
	})
}

func (c *moduleMap) Intercept(label string, mode Mode, action InterceptAction) {
	if _, ok := c.intercepts[label]; !ok {
		c.intercepts[label] = make(map[Mode]InterceptAction, 0)
	}
	c.intercepts[label][mode] = action
}

func (dg *DeviceGroup) Intercept(msg dcs.DCSMsg) (dcs.DCSMsg, error) {
	curMap := dg.CurrentMap()
	if m, ok := curMap.intercepts[msg.Message]; ok {
		if a, ok := m[dg.Mode()]; ok {
			return a(msg)
		}
	}

	return msg, nil
}

func (dg *DeviceGroup) HasModule(name string) bool {
	m, ok := dg.moduleMaps[name]
	return ok && m != nil
}

func (dg *DeviceGroup) SetModule(name string) {
	dg.currentModule = name
}

func (dg *DeviceGroup) CurrentMap() *moduleMap {
	return dg.moduleMaps[dg.currentModule]
}

func (dg *DeviceGroup) ModuleMap(name string) *moduleMap {
	alert.Sayf("creating module map '%s'", name)
	if m, ok := dg.moduleMaps[name]; ok {
		return m
	}

	module := &moduleMap{
		name:        name,
		deviceGroup: dg,
		controls:    make(map[string]map[string]Control, 0),
		intercepts:  make(map[string]map[Mode]InterceptAction, 0),
	}

	for _, dev := range dg.devices {
		module.controls[dev.name] = make(map[string]Control, 0)
		mapper := controlMapper{dev.name, module}
		mapControls(mapper.device, mapper)
	}

	dg.moduleMaps[name] = module

	return module
}

type controlMapper struct {
	device    string
	moduleMap *moduleMap
}

func (m controlMapper) AddControl(name string, control Control) {
	control.setParent(m.moduleMap.deviceGroup)
	m.moduleMap.controls[m.device][name] = control
}

func (m controlMapper) AddPOVControl(name string, control Control) {
	m.AddControl(name, control)
}
