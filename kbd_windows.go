package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/win"
)

var dll = syscall.NewLazyDLL("user32.dll")
var procKeybdEvent = dll.NewProc("keybd_event")
var procMapVirtualKey = dll.NewProc("MapVirtualKeyExW")
var procGetKeyboardLayout = dll.NewProc("GetKeyboardLayout")

var klOnce sync.Once
var keyboardLayout uintptr

func getKeyboardLayout() uintptr {
	klOnce.Do(func() {
		ret, _, _ := procGetKeyboardLayout.Call(0)
		keyboardLayout = ret
	})
	return keyboardLayout
}

func mapVirtualKey(vk uint16) uint16 {
	ret, _, _ := procMapVirtualKey.Call(uintptr(vk), 0, getKeyboardLayout())
	return uint16(ret)
}

func sendPause() {
	time.Sleep(50 * time.Millisecond)
}

func downKey(key uint16) {
	if key == K_PAUSE {
		procKeybdEvent.Call(uintptr(K_PAUSE&0xFF), 0, 0, 0)
	}

	var flags uint32 = win.KEYEVENTF_SCANCODE
	if key&0xFF00 >= 0xE000 {
		flags |= win.KEYEVENTF_EXTENDEDKEY
	}

	inputs := []win.KEYBD_INPUT{{
		Type: win.INPUT_KEYBOARD,
		Ki: win.KEYBDINPUT{
			WVk:         0,
			WScan:       key,
			DwFlags:     flags,
			Time:        0,
			DwExtraInfo: 0,
		},
	}}

	ret := win.SendInput(uint32(len(inputs)), unsafe.Pointer(&inputs[0]), int32(unsafe.Sizeof(inputs[0])))
	fmt.Printf("downKey %x send returned %d\n", key, ret)
}

func upKey(key uint16) {
	if key == K_PAUSE {
		procKeybdEvent.Call(uintptr(K_PAUSE&0xFF), 0, win.KEYEVENTF_KEYUP, 0)
		return
	}

	var flags uint32 = win.KEYEVENTF_SCANCODE | win.KEYEVENTF_KEYUP
	if key&0xFF00 >= 0xE000 {
		flags |= win.KEYEVENTF_EXTENDEDKEY
	}

	inputs := []win.KEYBD_INPUT{{
		Type: win.INPUT_KEYBOARD,
		Ki: win.KEYBDINPUT{
			WVk:         0,
			WScan:       key,
			DwFlags:     flags,
			Time:        0,
			DwExtraInfo: 0,
		},
	}}

	if key == K_PAUSE {
		inputs = append(inputs, win.KEYBD_INPUT{
			Type: win.INPUT_KEYBOARD,
			Ki: win.KEYBDINPUT{
				WVk:         0,
				WScan:       0x45,
				DwFlags:     flags ^ win.KEYEVENTF_EXTENDEDKEY,
				Time:        0,
				DwExtraInfo: 0,
			},
		})
	}

	ret := win.SendInput(uint32(len(inputs)), unsafe.Pointer(&inputs[0]), int32(unsafe.Sizeof(inputs[0])))
	fmt.Printf("upKey %x send returned %d\n", key, ret)
}

const (
	LEFT   uint8 = 0
	MIDDLE       = 1
	RIGHT        = 2
	DOWN   bool  = true
	UP           = false
)

func mouse(button uint8, down bool) {
	var flags uint32
	switch down {
	case DOWN:
		switch button {
		case LEFT:
			flags = win.MOUSEEVENTF_LEFTDOWN
		case MIDDLE:
			flags = win.MOUSEEVENTF_MIDDLEDOWN
		case RIGHT:
			flags = win.MOUSEEVENTF_RIGHTDOWN
		}
	case UP:
		switch button {
		case LEFT:
			flags = win.MOUSEEVENTF_LEFTUP
		case MIDDLE:
			flags = win.MOUSEEVENTF_MIDDLEUP
		case RIGHT:
			flags = win.MOUSEEVENTF_RIGHTUP
		}
	}
	input := win.MOUSE_INPUT{
		Type: win.INPUT_MOUSE,
		Mi: win.MOUSEINPUT{
			DwFlags: flags,
		},
	}
	ret := win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	fmt.Printf("mouse event %v %v send returned %d\n", button, down, ret)
}

const (
	K_ESC        uint16 = 0x01
	K_1                 = 0x02
	K_2                 = 0x03
	K_3                 = 0x04
	K_4                 = 0x05
	K_5                 = 0x06
	K_6                 = 0x07
	K_7                 = 0x08
	K_8                 = 0x09
	K_9                 = 0x0A
	K_0                 = 0x0B
	K_MINUS             = 0x0C
	K_EQUAL             = 0x0D
	K_BACKSPACE         = 0x0E
	K_TAB               = 0x0F
	K_Q                 = 0x10
	K_W                 = 0x11
	K_E                 = 0x12
	K_R                 = 0x13
	K_T                 = 0x14
	K_Y                 = 0x15
	K_U                 = 0x16
	K_I                 = 0x17
	K_O                 = 0x18
	K_P                 = 0x19
	K_LEFTBRACE         = 0x1A
	K_RIGHTBRACE        = 0x1B
	K_ENTER             = 0x1C
	K_LCONTROL          = 0x1D
	K_A                 = 0x1E
	K_S                 = 0x1F
	K_D                 = 0x20
	K_F                 = 0x21
	K_G                 = 0x22
	K_H                 = 0x23
	K_J                 = 0x24
	K_K                 = 0x25
	K_L                 = 0x26
	K_SEMICOLON         = 0x27
	K_QUOTE             = 0x28
	K_GRAVE             = 0x29
	K_LSHIFT            = 0x2A
	K_BACKSLASH         = 0x2B
	K_Z                 = 0x2C
	K_X                 = 0x2D
	K_C                 = 0x2E
	K_V                 = 0x2F
	K_B                 = 0x30
	K_N                 = 0x31
	K_M                 = 0x32
	K_COMMA             = 0x33
	K_PERIOD            = 0x34
	K_SLASH             = 0x35
	K_RSHIFT            = 0x36
	K_KPASTERISK        = 0x37
	K_LALT              = 0x38
	K_SPACE             = 0x39
	K_CAPSLOCK          = 0x3A
	K_F1                = 0x3B
	K_F2                = 0x3C
	K_F3                = 0x3D
	K_F4                = 0x3E
	K_F5                = 0x3F
	K_F6                = 0x40
	K_F7                = 0x41
	K_F8                = 0x42
	K_F9                = 0x43
	K_F10               = 0x44
	K_NUMLOCK           = 0x45
	K_SCROLLOCK         = 0x46
	K_KP7               = 0x47
	K_KP8               = 0x48
	K_KP9               = 0x49
	K_KPMINUS           = 0x4A
	K_KP4               = 0x4B
	K_KP5               = 0x4C
	K_KP6               = 0x4D
	K_KPPLUS            = 0x4E
	K_KP1               = 0x4F
	K_KP2               = 0x50
	K_KP3               = 0x51
	K_KP0               = 0x52
	K_KPDOT             = 0x53
	K_F11               = 0x57
	K_F12               = 0x58
	K_KPENTER           = 0xe01C
	K_RCONTROL          = 0xe01D
	K_KPSLASH           = 0xe02A
	K_RALT              = 0xe038
	K_HOME              = 0xe047
	K_UP                = 0xe048
	K_PGUP              = 0xe049
	K_LEFT              = 0xe04B
	K_RIGHT             = 0xe04D
	K_END               = 0xe04F
	K_DOWN              = 0xe050
	K_PGDOWN            = 0xe051
	K_INSERT            = 0xe052
	K_DELETE            = 0xe053

	K_PAUSE = 0xF013
)

/*
const (
	K_LBUTTON             uint16 = 0x01
	K_RBUTTON                    = 0x02
	K_CANCEL                     = 0x03
	K_MBUTTON                    = 0x04
	K_XBUTTON1                   = 0x05
	K_XBUTTON2                   = 0x06
	K_BACK                       = 0x08
	K_TAB                        = 0x09
	K_CLEAR                      = 0x0C
	K_ENTER                      = 0x0D
	K_SHIFT                      = 0x10
	K_CTRL                       = 0x11
	K_ALT                        = 0x12
	K_CAPITAL                    = 0x14
	K_KANA                       = 0x15
	K_HANGUEL                    = 0x15
	K_HANGUL                     = 0x15
	K_JUNJA                      = 0x17
	K_FINAL                      = 0x18
	K_HANJA                      = 0x19
	K_KANJI                      = 0x19
	K_ESC                        = 0x1B
	K_CONVERT                    = 0x1C
	K_NONCONVERT                 = 0x1D
	K_ACCEPT                     = 0x1E
	K_MODECHANGE                 = 0x1F
	K_SPACE                      = 0x20
	K_PRIOR                      = 0x21
	K_NEXT                       = 0x22
	K_END                        = 0x23
	K_HOME                       = 0x24
	K_LEFT                       = 0x25
	K_UP                         = 0x26
	K_RIGHT                      = 0x27
	K_DOWN                       = 0x28
	K_SELECT                     = 0x29
	K_PRINT                      = 0x2A
	K_EXECUTE                    = 0x2B
	K_SNAPSHOT                   = 0x2C
	K_INSERT                     = 0x2D
	K_DELETE                     = 0x2E
	K_HELP                       = 0x2F
	K_0                          = 0x30
	K_1                          = 0x31
	K_2                          = 0x32
	K_3                          = 0x33
	K_4                          = 0x34
	K_5                          = 0x35
	K_6                          = 0x36
	K_7                          = 0x37
	K_8                          = 0x38
	K_9                          = 0x39
	K_A                          = 0x41
	K_B                          = 0x42
	K_C                          = 0x43
	K_D                          = 0x44
	K_E                          = 0x45
	K_F                          = 0x46
	K_G                          = 0x47
	K_H                          = 0x48
	K_I                          = 0x49
	K_J                          = 0x4A
	K_K                          = 0x4B
	K_L                          = 0x4C
	K_M                          = 0x4D
	K_N                          = 0x4E
	K_O                          = 0x4F
	K_P                          = 0x50
	K_Q                          = 0x51
	K_R                          = 0x52
	K_S                          = 0x53
	K_T                          = 0x54
	K_U                          = 0x55
	K_V                          = 0x56
	K_W                          = 0x57
	K_X                          = 0x58
	K_Y                          = 0x59
	K_Z                          = 0x5A
	K_LWIN                       = 0x5B
	K_RWIN                       = 0x5C
	K_APPS                       = 0x5D
	K_SLEEP                      = 0x5F
	K_NUMPAD0                    = 0x60
	K_NUMPAD1                    = 0x61
	K_NUMPAD2                    = 0x62
	K_NUMPAD3                    = 0x63
	K_NUMPAD4                    = 0x64
	K_NUMPAD5                    = 0x65
	K_NUMPAD6                    = 0x66
	K_NUMPAD7                    = 0x67
	K_NUMPAD8                    = 0x68
	K_NUMPAD9                    = 0x69
	K_MULTIPLY                   = 0x6A
	K_ADD                        = 0x6B
	K_SEPARATOR                  = 0x6C
	K_SUBTRACT                   = 0x6D
	K_DECIMAL                    = 0x6E
	K_DIVIDE                     = 0x6F
	K_F1                         = 0x70
	K_F2                         = 0x71
	K_F3                         = 0x72
	K_F4                         = 0x73
	K_F5                         = 0x74
	K_F6                         = 0x75
	K_F7                         = 0x76
	K_F8                         = 0x77
	K_F9                         = 0x78
	K_F10                        = 0x79
	K_F11                        = 0x7A
	K_F12                        = 0x7B
	K_F13                        = 0x7C
	K_F14                        = 0x7D
	K_F15                        = 0x7E
	K_F16                        = 0x7F
	K_F17                        = 0x80
	K_F18                        = 0x81
	K_F19                        = 0x82
	K_F20                        = 0x83
	K_F21                        = 0x84
	K_F22                        = 0x85
	K_F23                        = 0x86
	K_F24                        = 0x87
	K_NUMLOCK                    = 0x90
	K_SCROLL                     = 0x91
	K_LSHIFT                     = 0xA0
	K_RSHIFT                     = 0xA1
	K_LCONTROL                   = 0xA2
	K_RCONTROL                   = 0xA3
	K_LALT                       = 0xA4
	K_RALT                       = 0xA5
	K_BROWSER_BACK               = 0xA6
	K_BROWSER_FORWARD            = 0xA7
	K_BROWSER_REFRESH            = 0xA8
	K_BROWSER_STOP               = 0xA9
	K_BROWSER_SEARCH             = 0xAA
	K_BROWSER_FAVORITES          = 0xAB
	K_BROWSER_HOME               = 0xAC
	K_VOLUME_MUTE                = 0xAD
	K_VOLUME_DOWN                = 0xAE
	K_VOLUME_UP                  = 0xAF
	K_MEDIA_NEXT_TRACK           = 0xB0
	K_MEDIA_PREV_TRACK           = 0xB1
	K_MEDIA_STOP                 = 0xB2
	K_MEDIA_PLAY_PAUSE           = 0xB3
	K_LAUNCH_MAIL                = 0xB4
	K_LAUNCH_MEDIA_SELECT        = 0xB5
	K_LAUNCH_APP1                = 0xB6
	K_LAUNCH_APP2                = 0xB7
	K_OEM_1                      = 0xBA
	K_OEM_PLUS                   = 0xBB
	K_OEM_COMMA                  = 0xBC
	K_OEM_MINUS                  = 0xBD
	K_OEM_PERIOD                 = 0xBE
	K_OEM_2                      = 0xBF
	K_OEM_3                      = 0xC0
	K_GRAVE                      = 0xC0
	K_OEM_4                      = 0xDB
	K_OEM_5                      = 0xDC
	K_BACKSLASH                  = 0xDC
	K_OEM_6                      = 0xDD
	K_OEM_7                      = 0xDE
	K_OEM_8                      = 0xDF
	K_OEM_102                    = 0xE2
	K_PROCESSKEY                 = 0xE5
	K_PACKET                     = 0xE7
	K_ATTN                       = 0xF6
	K_CRSEL                      = 0xF7
	K_EXSEL                      = 0xF8
	K_EREOF                      = 0xF9
	K_PLAY                       = 0xFA
	K_ZOOM                       = 0xFB
	K_NONAME                     = 0xFC
	K_PA1                        = 0xFD
	K_OEM_CLEAR                  = 0xFE
)
*/
