package main

import (
	"fmt"
	"time"

	"github.com/ianmcmahon/joybehar/controls"
)

type _dcsAction struct {
	msg  []string
	vals []string
}

func dcsAction(msgs ...string) _dcsAction {
	return _dcsAction{
		msg:  msgs,
		vals: []string{"0", "1", "2"}, // low, off, high
	}
}

func reverse(s []string) []string {
	out := []string{}
	for _, v := range s {
		out = append(out, v)
	}
	return out
}

func (a _dcsAction) invert() _dcsAction {
	a.vals = reverse(a.vals)
	return a
}

func (a _dcsAction) withVals(s ...string) _dcsAction {
	a.vals = s
	return a
}

func (a _dcsAction) HandleEvent(control controls.Control, state controls.State) {
	val := ""
	switch len(control.ButtonIDs()) {
	case 1:
		switch state {
		case controls.STATE_OFF:
			val = a.vals[0]
		case controls.STATE_ON:
			val = a.vals[1]
		}
	case 2:
		switch state {
		case controls.STATE_LOW:
			val = a.vals[0]
		case controls.STATE_OFF:
			val = a.vals[1]
		case controls.STATE_HI:
			val = a.vals[2]
		}
	default:
		panic(fmt.Errorf("don't know how to handle event for %+v\n", control))
	}
	for _, msg := range a.msg {
		if val != "" {
			dcs.Send(msg, val)
		}
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
	keys []uint16
}

func keyAction(keys ...uint16) _keyAction {
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
	keys []uint16
}

func keyPulse(keys ...uint16) _keyPulse {
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

type _keyPress struct {
	keys []uint16
}

func keyPress(keys ...uint16) _keyPress {
	return _keyPress{keys}
}

func (a _keyPress) HandleEvent(control controls.Control, state controls.State) {
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

type _keyToggle struct {
	on, off controls.Action
}

func keyToggle(on, off controls.Action) _keyToggle {
	return _keyToggle{on, off}
}

func (a _keyToggle) HandleEvent(control controls.Control, state controls.State) {
	switch state {
	case controls.STATE_ON:
		a.on.HandleEvent(control, controls.STATE_ON)
	case controls.STATE_OFF:
		a.off.HandleEvent(control, controls.STATE_ON)
	}
}

type _mouseToggle struct {
	button uint8
}

func mouseToggle(button uint8) _mouseToggle {
	return _mouseToggle{button}
}

func (a _mouseToggle) HandleEvent(control controls.Control, state controls.State) {
	switch state {
	case controls.STATE_ON:
		mouse(a.button, DOWN)
	case controls.STATE_OFF:
		mouse(a.button, UP)
	}
}

// tempo only sends "short press" events, ie it sends a 50ms on/off pulse to the downstream actions
// only suitable for buttons and momentary toggles
type _tempo struct {
	short     controls.Action
	long      controls.Action
	threshold time.Duration
	pressed   *time.Time
}

func tempo(shortPress, longPress controls.Action) *_tempo {
	return &_tempo{
		short:     shortPress,
		long:      longPress,
		threshold: 500 * time.Millisecond,
	}
}

func (a *_tempo) HandleEvent(control controls.Control, state controls.State) {
	switch state {
	case controls.STATE_ON:
		t := time.Now()
		a.pressed = &t
		go func() {
			fmt.Printf("pressin and sleepin: %s\n", a.pressed)
			time.Sleep(a.threshold)
			if a.pressed != nil {
				a.pressed = nil
				fmt.Printf("timed out, sending long press\n")
				a.long.HandleEvent(control, controls.STATE_ON)
				time.Sleep(50 * time.Millisecond)
				a.long.HandleEvent(control, controls.STATE_OFF)
			}
		}()
	case controls.STATE_OFF:
		if a.pressed != nil {
			fmt.Printf("natural release, sending short press\n")
			a.pressed = nil
			go func() {
				a.short.HandleEvent(control, controls.STATE_ON)
				time.Sleep(50 * time.Millisecond)
				a.short.HandleEvent(control, controls.STATE_OFF)
			}()
		}
	}
}
