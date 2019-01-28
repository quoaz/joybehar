package main

import "time"

// stick buttons
const (
	TRIGGER_FIRST  uint8 = 0
	WEAPON_RELEASE uint8 = 1
	NWS            uint8 = 2
	PINKY_PADDLE   uint8 = 3
	INDEX_BUTTON   uint8 = 4
	TRIGGER_SECOND uint8 = 5
	TMS_FWD        uint8 = 6
	TMS_RIGHT      uint8 = 7
	TMS_AFT        uint8 = 8
	TMS_LEFT       uint8 = 9
	DMS_FWD        uint8 = 10
	DMS_RIGHT      uint8 = 11
	DMS_AFT        uint8 = 12
	DMS_LEFT       uint8 = 13
	CMS_FWD        uint8 = 14
	CMS_RIGHT      uint8 = 15
	CMS_AFT        uint8 = 16
	CMS_LEFT       uint8 = 17
	CMS_DEPRESS    uint8 = 18
)

// throttle buttons
const (
	SLEW_PRESS         uint8 = 0
	MIC_SWITCH_PRESS   uint8 = 1
	MIC_SWITCH_UP      uint8 = 2
	MIC_SWITCH_FWD     uint8 = 3
	MIC_SWITCH_DOWN    uint8 = 4
	MIC_SWITCH_AFT     uint8 = 5
	SPEEDBRAKE_DEPLOY  uint8 = 6
	SPEEDBRAKE_RETRACT uint8 = 7
	BOAT_SWITCH_FWD    uint8 = 8
	BOAT_SWITCH_AFT    uint8 = 9
	CHINA_HAT_FWD      uint8 = 10
	CHINA_HAT_AFT      uint8 = 11
	PINKY_SWITCH_AFT   uint8 = 12
	PINKY_SWITCH_FWD   uint8 = 13
	PINKY_BUTTON       uint8 = 14
	OVERRIDE_LEFT      uint8 = 15
	OVERRIDE_RIGHT     uint8 = 16
	ENG_LEFT_MOTOR     uint8 = 17
	ENG_RIGHT_MOTOR    uint8 = 18
	APU                uint8 = 19
	GEAR_HORN_SILENCE  uint8 = 20
	FLAPS_UP           uint8 = 21
	FLAPS_DOWN         uint8 = 22
	EAC_ARM            uint8 = 23
	RADAR_ALTIMETER    uint8 = 24
	AP_ENGAGE          uint8 = 25
	AP_PATH            uint8 = 26
	AP_ALT             uint8 = 27
	RIGHT_IDLE         uint8 = 28
	LEFT_IDLE          uint8 = 29
	ENG_LEFT_IGN       uint8 = 30
	ENG_RIGHT_IGN      uint8 = 31
)

type Mode uint8

const MODE_SHIFT Mode = 1

type press struct {
	id   uint8
	mode Mode

	longPress bool
}

func (p press) Mode(m Mode) press {
	p.mode = p.mode & m
	return p
}

func Press(id uint8) press {
	return press{
		id:        id,
		mode:      Mode(0),
		longPress: false,
	}
}

func LongPress(id uint8) press {
	return press{
		id:        id,
		mode:      Mode(0),
		longPress: false,
	}
}

// receives the raw button events and
// passes along modified events such as tempo and shifted events
type buttonHandler struct {
	Device *device
	Button uint8

	Actions map[bool]map[Mode]func(rawEvent)

	delay time.Duration
}

func (d *device) ButtonHandler(id uint8) *buttonHandler {
	handler := &buttonHandler{
		Device:  d,
		Button:  id,
		Actions: make(map[bool]map[Mode]func(rawEvent), 0),
		delay:   500 * time.Millisecond,
	}
	handler.Actions[false] = make(map[Mode]func(rawEvent), 0)
	handler.Actions[true] = make(map[Mode]func(rawEvent), 0)

	return handler
}

func (h *buttonHandler) Handle(mode Mode, event rawEvent) {
	if h.Actions[true][mode] == nil {

	} else { // handle long press

	}
}
