package main

import (
	"time"

	"github.com/ianmcmahon/joybehar/controls"
)

func F5EMap(group *controls.DeviceGroup) {
	m := group.ModuleMap("F-5E-3")

	m.ModeToggle("stick", "paddle", controls.MODE_SHIFT, controls.MODE_NORM)
	m.Control("stick", "trigger1").Action(controls.MODE_ALL, keyAction(K_T))
	m.Control("stick", "trigger2").Action(controls.MODE_ALL, keyAction(K_SPACE))

	m.Control("stick", "weaponrelease").Action(controls.MODE_NORM, dcsAction("WEAPON_RELEASE"))
	m.Control("stick", "weaponrelease").Action(controls.MODE_SHIFT, dcsToggle("RWR_PWR"))

	m.Control("stick", "nws").Action(controls.MODE_ALL, dcsAction("NWS"))

	m.Control("stick", "index").Action(controls.MODE_NORM, dcsAction("FL_CHAFF_BT"))
	//m.Control("index").Action(controls.MODE_SHIFT, tempo(keyPulse(0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0x5B, 0x5C), keyPulse(K_LSHIFT, K_R)))

	m.Control("stick", "tms_up").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS1", "ARMPOS7"))
	m.Control("stick", "tms_dn").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS4"))
	m.Control("stick", "tms_rt").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS3", "ARMPOS5"))
	m.Control("stick", "tms_lt").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS2", "ARMPOS6"))

	m.Control("stick", "cms_up").Action(controls.MODE_NORM, keyPulse(K_5))                           // dogfight missile
	m.Control("stick", "cms_dn").Action(controls.MODE_NORM, keyPulse(K_6))                           // dogfight guns
	m.Control("stick", "cms_dp").Action(controls.MODE_NORM, tempo(keyPulse(K_R), keyPulse(K_PAUSE))) // dogfight resume, pause
	m.Control("stick", "cms_dp").Action(controls.MODE_SHIFT, tempo(keyPulse(K_ESC), keyPulse(K_F10)))

	m.Control("stick", "dms_up").Action(controls.MODE_NORM, keyAction(K_KPPLUS)) // discord PTT

	m.Control("throttle", "slewpress").Action(controls.MODE_NORM, mouseToggle(LEFT))
	m.Control("throttle", "slewpress").Action(controls.MODE_SHIFT, mouseToggle(RIGHT))

	m.Control("throttle", "mic_up").Action(controls.MODE_NORM, keyPulse(K_LCONTROL, K_LALT, K_GRAVE))
	m.Control("throttle", "mic_dn").Action(controls.MODE_NORM, keyPress(K_GRAVE))
	m.Control("throttle", "mic_aft").Action(controls.MODE_SHIFT, keyPulse(K_BACKSLASH))

	m.Control("throttle", "mic_dp").Action(controls.MODE_SHIFT, keyPulse(K_RSHIFT, K_K))

	m.Control("throttle", "tdc_fwd").Action(controls.MODE_NORM, keyPulse(K_P))
	m.Control("throttle", "tdc_aft").Action(controls.MODE_NORM, keyPulse(K_F1))
	m.Control("throttle", "tdc_left").Action(controls.MODE_NORM, keyPulse(K_F2))
	m.Control("throttle", "tdc_right").Action(controls.MODE_NORM, keyPulse(K_F6))

	m.Control("throttle", "tdc_fwd").Action(controls.MODE_SHIFT, mouseScroll(-300))
	m.Control("throttle", "tdc_aft").Action(controls.MODE_SHIFT, mouseScroll(300))
	m.Control("throttle", "tdc_left").Action(controls.MODE_SHIFT, keyPulse(K_LEFTBRACE))
	m.Control("throttle", "tdc_right").Action(controls.MODE_SHIFT, keyPulse(K_RIGHTBRACE))

	m.Control("throttle", "speedbrake").Action(controls.MODE_ALL, dcsAction("SPEED"))
	m.Control("throttle", "boatswitch").Action(controls.MODE_ALL, dcsAction("A_FLAPS"))
	m.Control("throttle", "chinahat").Action(controls.MODE_ALL, dcsAction("RADAR_RANGE").withVals("DEC", "", "INC"))

	m.Control("throttle", "flaps").Action(controls.MODE_ALL, dcsAction("FLAPS"))
	m.Control("throttle", "eacarm").Action(controls.MODE_ALL, dcsAction("LG_LEVER_SWITCH"))
	m.Control("throttle", "apselect").Action(controls.MODE_ALL, masterArm())
	m.Control("throttle", "pinkybutton").Action(controls.MODE_NORM, dcsAction("MISSILE_UNCAGE"))
	m.Control("throttle", "pinkybutton").Action(controls.MODE_SHIFT, keyPulse(K_LCONTROL, K_E))

	m.Control("throttle", "leftidle").Action(controls.MODE_ALL, keyToggle(keyPulse(K_RALT, K_HOME), keyPulse(K_RALT, K_END)))
	m.Control("throttle", "rightidle").Action(controls.MODE_ALL, keyToggle(keyPulse(K_RSHIFT, K_HOME), keyPulse(K_RSHIFT, K_END)))

	m.Intercept("HSI_HDG_KNOB", controls.MODE_SHIFT, makeIncremental("TACAN_10"))
	m.Intercept("HSI_CRS_KNOB", controls.MODE_SHIFT, makeIncremental("TACAN_1"))
}

type _masterArm struct{}

func masterArm() _masterArm {
	return _masterArm{}
}

func (a _masterArm) HandleEvent(_ controls.Control, state controls.State) {
	switch state {
	case controls.STATE_HI:
		dcsAgent.Send("MASTER_ARM_GUARD", "1")
		time.Sleep(50 * time.Millisecond)
		dcsAgent.Send("MASTER_ARM", "2")
	case controls.STATE_OFF:
		dcsAgent.Send("MASTER_ARM", "1")
		//dcsAgent.Send("MASTER_ARM_GUARD", "0")
	case controls.STATE_LOW:
		//dcsAgent.Send("MASTER_ARM_GUARD", "1")
		dcsAgent.Send("MASTER_ARM", "0")
	}
}
