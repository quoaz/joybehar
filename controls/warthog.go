package controls

const (
	MODE_NORM  Mode = 1 << 0
	MODE_SHIFT Mode = 1 << 1
	MODE_ALL   Mode = 0xFF
)

func WarthogGroup() *deviceGroup {
	group := DeviceGroup()
	group.AddDevice(WarthogThrottle())
	group.AddDevice(WarthogStick())

	return group
}

func WarthogThrottle() *device {
	d := NewDevice("throttle")

	d.AddControl("slewpress", Button(SLEW_PRESS))
	d.AddControl("speedbrake", Toggle3(SPEEDBRAKE_DEPLOY, SPEEDBRAKE_RETRACT))
	d.AddControl("boatswitch", Toggle3(BOAT_SWITCH_FWD, BOAT_SWITCH_AFT))
	d.AddControl("chinahat", Toggle3(CHINA_HAT_FWD, CHINA_HAT_AFT))
	d.AddControl("pinkyswitch", Toggle3(PINKY_SWITCH_FWD, PINKY_SWITCH_AFT))
	d.AddControl("pinkybutton", Button(PINKY_BUTTON))
	d.AddControl("overrideleft", Toggle(OVERRIDE_LEFT))
	d.AddControl("overrideright", Toggle(OVERRIDE_RIGHT))
	d.AddControl("engineleft", Toggle3(ENG_LEFT_IGN, ENG_LEFT_MOTOR))
	d.AddControl("engineright", Toggle3(ENG_RIGHT_IGN, ENG_RIGHT_MOTOR))
	d.AddControl("apu", Toggle(APU))
	d.AddControl("gearhornsilence", Button(GEAR_HORN_SILENCE))
	d.AddControl("flaps", Toggle3(FLAPS_UP, FLAPS_DOWN))
	d.AddControl("eacarm", Toggle(EAC_ARM))
	d.AddControl("radalt", Toggle(RADAR_ALTIMETER))
	d.AddControl("apengage", Button(AP_ENGAGE))
	d.AddControl("apselect", Toggle3(AP_PATH, AP_ALT))
	d.AddControl("leftidle", Toggle(LEFT_IDLE))
	d.AddControl("rightidle", Toggle(RIGHT_IDLE))

	d.AddControl("mic_up", Button(MIC_SWITCH_UP))
	d.AddControl("mic_fwd", Button(MIC_SWITCH_FWD))
	d.AddControl("mic_dn", Button(MIC_SWITCH_DOWN))
	d.AddControl("mic_aft", Button(MIC_SWITCH_AFT))
	d.AddControl("mic_dp", Button(MIC_SWITCH_PRESS))

	d.AddPOVControl("tdc_fwd", Button(TDC_FWD))
	d.AddPOVControl("tdc_aft", Button(TDC_AFT))
	d.AddPOVControl("tdc_left", Button(TDC_LEFT))
	d.AddPOVControl("tdc_right", Button(TDC_RIGHT))

	return d
}

func WarthogStick() *device {
	d := NewDevice("stick")

	d.AddControl("trigger1", Button(TRIGGER_FIRST))
	d.AddControl("trigger2", Button(TRIGGER_SECOND))
	d.AddControl("weaponrelease", Button(WEAPON_RELEASE))
	d.AddControl("nws", Button(NWS))
	d.AddControl("index", Button(INDEX_BUTTON))
	d.AddControl("paddle", Button(PINKY_PADDLE))

	/* I'm not sure fourways are worth the effort; if there's no four-state dcs objects to bind to
	d.AddControl("tms", FourWay(TMS_FWD, TMS_RIGHT, TMS_AFT, TMS_LEFT))
	d.AddControl("dms", FourWay(DMS_FWD, DMS_RIGHT, DMS_AFT, DMS_LEFT))
	d.AddControl("cms", FourWay(CMS_FWD, CMS_RIGHT, CMS_AFT, CMS_LEFT).Depress(CMS_DEPRESS))
	*/
	d.AddControl("tms_up", Button(TMS_FWD))
	d.AddControl("tms_rt", Button(TMS_RIGHT))
	d.AddControl("tms_dn", Button(TMS_AFT))
	d.AddControl("tms_lt", Button(TMS_LEFT))

	d.AddControl("dms_up", Button(DMS_FWD))
	d.AddControl("dms_rt", Button(DMS_RIGHT))
	d.AddControl("dms_dn", Button(DMS_AFT))
	d.AddControl("dms_lt", Button(DMS_LEFT))

	d.AddControl("cms_up", Button(CMS_FWD))
	d.AddControl("cms_rt", Button(CMS_RIGHT))
	d.AddControl("cms_dn", Button(CMS_AFT))
	d.AddControl("cms_lt", Button(CMS_LEFT))
	d.AddControl("cms_dp", Button(CMS_DEPRESS))

	return d
}

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
	TRIM_NOSE_DN   uint8 = 19
	TRIM_NOSE_UP   uint8 = 20
	TRIM_LWD       uint8 = 21
	TRIM_RWD       uint8 = 22
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
	TDC_AFT            uint8 = 32
	TDC_FWD            uint8 = 33
	TDC_LEFT           uint8 = 34
	TDC_RIGHT          uint8 = 35
)
