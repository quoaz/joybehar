package main

import (
	"fmt"
	"time"

	"github.com/ianmcmahon/joybehar/controls"
)

var dcs *dcsAgent

/*
	Need module support

	F-5E-3
	F-14B
*/

func main() {
	dcs = DCSAgent()

	warthog := controls.WarthogGroup()

	stick := warthog.Device("stick")

	warthog.ModeToggle(stick.Control("paddle"), controls.MODE_SHIFT, controls.MODE_NORM)

	stick.Control("trigger1").Action(controls.MODE_ALL, keyAction(K_T))
	stick.Control("trigger2").Action(controls.MODE_ALL, keyAction(K_SPACE))

	stick.Control("weaponrelease").Action(controls.MODE_NORM, dcsAction("WEAPON_RELEASE"))
	stick.Control("weaponrelease").Action(controls.MODE_SHIFT, dcsToggle("RWR_PWR"))

	stick.Control("nws").Action(controls.MODE_ALL, dcsAction("NWS"))

	stick.Control("index").Action(controls.MODE_NORM, dcsAction("FL_CHAFF_BT"))
	//stick.Control("index").Action(controls.MODE_SHIFT, tempo(keyPulse(0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0x5B, 0x5C), keyPulse(K_LSHIFT, K_R)))

	stick.Control("tms_up").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS1", "ARMPOS7"))
	stick.Control("tms_dn").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS4"))
	stick.Control("tms_rt").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS3", "ARMPOS5"))
	stick.Control("tms_lt").Action(controls.MODE_SHIFT, dcsToggle("ARMPOS2", "ARMPOS6"))

	stick.Control("cms_up").Action(controls.MODE_NORM, keyPulse(K_5))                           // dogfight missile
	stick.Control("cms_dn").Action(controls.MODE_NORM, keyPulse(K_6))                           // dogfight guns
	stick.Control("cms_dp").Action(controls.MODE_NORM, tempo(keyPulse(K_R), keyPulse(K_PAUSE))) // dogfight resume, pause
	stick.Control("cms_dp").Action(controls.MODE_SHIFT, tempo(keyPulse(K_ESC), keyPulse(K_F10)))

	throttle := warthog.Device("throttle")

	throttle.Control("slewpress").Action(controls.MODE_NORM, mouseToggle(LEFT))
	throttle.Control("slewpress").Action(controls.MODE_SHIFT, mouseToggle(RIGHT))

	throttle.Control("mic_up").Action(controls.MODE_NORM, keyPulse(K_LCONTROL, K_LALT, K_GRAVE))
	throttle.Control("mic_dn").Action(controls.MODE_NORM, keyPress(K_GRAVE))
	throttle.Control("mic_aft").Action(controls.MODE_SHIFT, keyPulse(K_BACKSLASH))

	throttle.Control("tdc_fwd").Action(controls.MODE_NORM, keyPulse(K_P))
	throttle.Control("tdc_aft").Action(controls.MODE_NORM, keyPulse(K_F1))
	throttle.Control("tdc_left").Action(controls.MODE_NORM, keyPulse(K_F2))
	throttle.Control("tdc_right").Action(controls.MODE_NORM, keyPulse(K_F6))

	throttle.Control("speedbrake").Action(controls.MODE_ALL, dcsAction("SPEED"))
	throttle.Control("boatswitch").Action(controls.MODE_ALL, dcsAction("A_FLAPS"))
	throttle.Control("chinahat").Action(controls.MODE_ALL, dcsAction("RADAR_RANGE").withVals("DEC", "", "INC"))

	throttle.Control("flaps").Action(controls.MODE_ALL, dcsAction("FLAPS"))
	throttle.Control("eacarm").Action(controls.MODE_ALL, dcsAction("LG_LEVER_SWITCH"))
	throttle.Control("apselect").Action(controls.MODE_ALL, masterArm())

	throttle.Control("leftidle").Action(controls.MODE_ALL, keyToggle(keyPulse(K_RALT, K_END), keyPulse(K_RALT, K_HOME)))
	throttle.Control("rightidle").Action(controls.MODE_ALL, keyToggle(keyPulse(K_RSHIFT, K_END), keyPulse(K_RSHIFT, K_HOME)))

	dcs.Register(&StringOutput{Addr: 0x0000, MaxLength: 16, Action: func(_ uint16, s string) {
		fmt.Printf("Airplane: %s\n", s)
	}})
	go dcs.Receive()

	for {
		time.Sleep(1 * time.Second)
	}
}
