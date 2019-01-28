package main

type F5E struct{}

func (f F5E) Stick() *device {
	stick := Device("stick")

	stick.DCSToggle(Press(NWS), "NWS", "0", "1")
	stick.DCSToggle(LongPress(NWS), "LNWS", "0", "1")
	stick.DCSToggle(Press(NWS).Mode(MODE_SHIFT), "SNWS", "0", "1")
	stick.DCSToggle(LongPress(NWS).Mode(MODE_SHIFT), "SLNWS", "0", "1")

	stick.DCSToggle(Press(WEAPON_RELEASE), "WEAPON_RELEASE", "0", "1")
	stick.DCSToggle(Press(INDEX_BUTTON), "FL_CHAFF_BT", "0", "1")
	// MISSILE_UNCAGE not sure where to put this, maybe on trigger first detent
	// dogfight resume and trim don't seem to be expose

	return stick
}

func (f F5E) Throttle() *device {
	throttle := Device("throttle")

	throttle.DCS3Pos(Press(BOAT_SWITCH_FWD), Press(BOAT_SWITCH_AFT), "A_FLAPS", "0", "1", "2")
	throttle.DCS3Pos(Press(FLAPS_UP), Press(FLAPS_DOWN), "FLAPS", "0", "1", "2")

	throttle.DCS3Pos(Press(SPEEDBRAKE_DEPLOY), Press(SPEEDBRAKE_RETRACT), "SPEED", "0", "1", "2")

	throttle.DCSToggle(Press(EAC_ARM), "LG_LEVER_SWITCH", "0", "1")

	// 2pos battery SW_BATTERY 0-1
	// eng start R_START/L_START 0-1 or "TOGGLE"
	// boost pumps L_BOOSTPUMP 0-1
	// nose strut NS_STRUCT 0-1
	// 2pos landing light LG_LIGHT 0-1
	// radar elevation RADAR_ELEVATION 0-65535
	// RADAR_MODE 0..3
	// RADAR_RANGE 0..3 or INC/DEC (for china hat)
	// 3pos MASTER_ARM 0..2
	// MASTER_ARM_GUARD 0..1 send with logic
	//

	return throttle
}

/*
// other useful stuff for extension panels

*/
