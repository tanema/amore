// +build js

package ui

const (
	//mouse buttons
	BUTTON_LEFT   MouseButton = 0
	BUTTON_MIDDLE MouseButton = 1
	BUTTON_RIGHT  MouseButton = 2
	BUTTON_X1     MouseButton = 3
	BUTTON_X2     MouseButton = 4

	//keyboard keys
	K_UNKNOWN      Keycode = ""
	K_RETURN       Keycode = "Enter"
	K_ESCAPE       Keycode = "Escape"
	K_BACKSPACE    Keycode = "Backspace"
	K_TAB          Keycode = "Tab"
	K_SPACE        Keycode = "Space"
	K_COMMA        Keycode = "Comma"
	K_MINUS        Keycode = "Minus"
	K_PERIOD       Keycode = "Period"
	K_SLASH        Keycode = "Slash"
	K_0            Keycode = "Digit0"
	K_1            Keycode = "Digit1"
	K_2            Keycode = "Digit2"
	K_3            Keycode = "Digit3"
	K_4            Keycode = "Digit4"
	K_5            Keycode = "Digit5"
	K_6            Keycode = "Digit6"
	K_7            Keycode = "Digit7"
	K_8            Keycode = "Digit8"
	K_9            Keycode = "Digit9"
	K_SEMICOLON    Keycode = "Semicolon"
	K_EQUALS       Keycode = "Equal"
	K_LEFTBRACKET  Keycode = "BracketLeft"
	K_BACKSLASH    Keycode = "BackSlash"
	K_RIGHTBRACKET Keycode = "BracketRight"
	K_a            Keycode = "KeyA"
	K_b            Keycode = "KeyB"
	K_c            Keycode = "KeyC"
	K_d            Keycode = "KeyD"
	K_e            Keycode = "KeyE"
	K_f            Keycode = "KeyF"
	K_g            Keycode = "KeyG"
	K_h            Keycode = "KeyH"
	K_i            Keycode = "KeyI"
	K_j            Keycode = "KeyJ"
	K_k            Keycode = "KeyK"
	K_l            Keycode = "KeyL"
	K_m            Keycode = "KeyM"
	K_n            Keycode = "KeyN"
	K_o            Keycode = "KeyO"
	K_p            Keycode = "KeyP"
	K_q            Keycode = "KeyQ"
	K_r            Keycode = "KeyR"
	K_s            Keycode = "KeyS"
	K_t            Keycode = "KeyT"
	K_u            Keycode = "KeyU"
	K_v            Keycode = "KeyV"
	K_w            Keycode = "KeyW"
	K_x            Keycode = "KeyX"
	K_y            Keycode = "KeyY"
	K_z            Keycode = "KeyZ"
	K_CAPSLOCK     Keycode = "CapsLock"
	K_F1           Keycode = "F1"
	K_F2           Keycode = "F2"
	K_F3           Keycode = "F3"
	K_F4           Keycode = "F4"
	K_F5           Keycode = "F5"
	K_F6           Keycode = "F6"
	K_F7           Keycode = "F7"
	K_F8           Keycode = "F8"
	K_F9           Keycode = "F9"
	K_F10          Keycode = "F10"
	K_F11          Keycode = "F11"
	K_F12          Keycode = "F12"
	K_PRINTSCREEN  Keycode = ""
	K_SCROLLLOCK   Keycode = ""
	K_PAUSE        Keycode = ""
	K_INSERT       Keycode = ""
	K_HOME         Keycode = "Home"
	K_PAGEUP       Keycode = "PageUp"
	K_DELETE       Keycode = "Delete"
	K_END          Keycode = "End"
	K_PAGEDOWN     Keycode = "PageDown"
	K_RIGHT        Keycode = "ArrowRight"
	K_LEFT         Keycode = "ArrowLeft"
	K_DOWN         Keycode = "ArrowDown"
	K_UP           Keycode = "ArrowUp"
	K_MUTE         Keycode = ""
	K_VOLUMEUP     Keycode = ""
	K_VOLUMEDOWN   Keycode = ""
	K_LCTRL        Keycode = "ControlLeft"
	K_LSHIFT       Keycode = "ShiftLeft"
	K_LALT         Keycode = "AltLeft"
	K_LGUI         Keycode = ""
	K_RCTRL        Keycode = "ControlRight"
	K_RSHIFT       Keycode = "ShiftRight"
	K_RALT         Keycode = "AltRight"
	K_RGUI         Keycode = ""
	K_MODE         Keycode = ""
	K_AUDIONEXT    Keycode = ""
	K_AUDIOPREV    Keycode = ""
	K_AUDIOSTOP    Keycode = ""
	K_AUDIOPLAY    Keycode = ""
	K_AUDIOMUTE    Keycode = ""

	//keyboard scancodes
	SCANCODE_UNKNOWN Scancode = 0

	SCANCODE_0 Scancode = iota + 48
	SCANCODE_1
	SCANCODE_2
	SCANCODE_3
	SCANCODE_4
	SCANCODE_5
	SCANCODE_6
	SCANCODE_7
	SCANCODE_8
	SCANCODE_9
	SCANCODE_A Scancode = iota + 65
	SCANCODE_B
	SCANCODE_C
	SCANCODE_D
	SCANCODE_E
	SCANCODE_F
	SCANCODE_G
	SCANCODE_H
	SCANCODE_I
	SCANCODE_J
	SCANCODE_K
	SCANCODE_L
	SCANCODE_M
	SCANCODE_N
	SCANCODE_O
	SCANCODE_P
	SCANCODE_Q
	SCANCODE_R
	SCANCODE_S
	SCANCODE_T
	SCANCODE_U
	SCANCODE_V
	SCANCODE_W
	SCANCODE_X
	SCANCODE_Y
	SCANCODE_Z
	SCANCODE_RETURN       Scancode = 13
	SCANCODE_ESCAPE       Scancode = 27
	SCANCODE_BACKSPACE    Scancode = 8
	SCANCODE_TAB          Scancode = 9
	SCANCODE_SPACE        Scancode = 32
	SCANCODE_MINUS        Scancode = 189
	SCANCODE_EQUALS       Scancode = 187
	SCANCODE_LEFTBRACKET  Scancode = 219
	SCANCODE_BACKSLASH    Scancode = 220
	SCANCODE_RIGHTBRACKET Scancode = 221
	SCANCODE_APOSTROPHE   Scancode = 222
	SCANCODE_SEMICOLON    Scancode = 186
	SCANCODE_GRAVE        Scancode = 0
	SCANCODE_COMMA        Scancode = 188
	SCANCODE_PERIOD       Scancode = 190
	SCANCODE_SLASH        Scancode = 191
	SCANCODE_CAPSLOCK     Scancode = 20
	SCANCODE_F1           Scancode = iota + 112
	SCANCODE_F2
	SCANCODE_F3
	SCANCODE_F4
	SCANCODE_F5
	SCANCODE_F6
	SCANCODE_F7
	SCANCODE_F8
	SCANCODE_F9
	SCANCODE_F10
	SCANCODE_F11
	SCANCODE_F12
	SCANCODE_PRINTSCREEN        Scancode = 0
	SCANCODE_SCROLLLOCK         Scancode = 0
	SCANCODE_PAUSE              Scancode = 0
	SCANCODE_INSERT             Scancode = 0
	SCANCODE_HOME               Scancode = 36
	SCANCODE_PAGEUP             Scancode = 33
	SCANCODE_DELETE             Scancode = 46
	SCANCODE_END                Scancode = 35
	SCANCODE_PAGEDOWN           Scancode = 34
	SCANCODE_RIGHT              Scancode = 39
	SCANCODE_LEFT               Scancode = 37
	SCANCODE_DOWN               Scancode = 40
	SCANCODE_UP                 Scancode = 38
	SCANCODE_NUMLOCKCLEAR       Scancode = 0
	SCANCODE_NONUSBACKSLASH     Scancode = 0
	SCANCODE_APPLICATION        Scancode = 0
	SCANCODE_POWER              Scancode = 0
	SCANCODE_F13                Scancode = 0
	SCANCODE_F14                Scancode = 0
	SCANCODE_F15                Scancode = 0
	SCANCODE_F16                Scancode = 0
	SCANCODE_F17                Scancode = 0
	SCANCODE_F18                Scancode = 0
	SCANCODE_F19                Scancode = 0
	SCANCODE_F20                Scancode = 0
	SCANCODE_F21                Scancode = 0
	SCANCODE_F22                Scancode = 0
	SCANCODE_F23                Scancode = 0
	SCANCODE_F24                Scancode = 0
	SCANCODE_EXECUTE            Scancode = 0
	SCANCODE_HELP               Scancode = 0
	SCANCODE_MENU               Scancode = 0
	SCANCODE_SELECT             Scancode = 0
	SCANCODE_STOP               Scancode = 0
	SCANCODE_AGAIN              Scancode = 0
	SCANCODE_UNDO               Scancode = 0
	SCANCODE_CUT                Scancode = 0
	SCANCODE_COPY               Scancode = 0
	SCANCODE_PASTE              Scancode = 0
	SCANCODE_FIND               Scancode = 0
	SCANCODE_MUTE               Scancode = 0
	SCANCODE_VOLUMEUP           Scancode = 0
	SCANCODE_VOLUMEDOWN         Scancode = 0
	SCANCODE_INTERNATIONAL1     Scancode = 0
	SCANCODE_INTERNATIONAL2     Scancode = 0
	SCANCODE_INTERNATIONAL3     Scancode = 0
	SCANCODE_INTERNATIONAL4     Scancode = 0
	SCANCODE_INTERNATIONAL5     Scancode = 0
	SCANCODE_INTERNATIONAL6     Scancode = 0
	SCANCODE_INTERNATIONAL7     Scancode = 0
	SCANCODE_INTERNATIONAL8     Scancode = 0
	SCANCODE_INTERNATIONAL9     Scancode = 0
	SCANCODE_LANG1              Scancode = 0
	SCANCODE_LANG2              Scancode = 0
	SCANCODE_LANG3              Scancode = 0
	SCANCODE_LANG4              Scancode = 0
	SCANCODE_LANG5              Scancode = 0
	SCANCODE_LANG6              Scancode = 0
	SCANCODE_LANG7              Scancode = 0
	SCANCODE_LANG8              Scancode = 0
	SCANCODE_LANG9              Scancode = 0
	SCANCODE_ALTERASE           Scancode = 0
	SCANCODE_SYSREQ             Scancode = 0
	SCANCODE_CANCEL             Scancode = 0
	SCANCODE_CLEAR              Scancode = 0
	SCANCODE_PRIOR              Scancode = 0
	SCANCODE_RETURN2            Scancode = 13
	SCANCODE_SEPARATOR          Scancode = 0
	SCANCODE_OUT                Scancode = 0
	SCANCODE_OPER               Scancode = 0
	SCANCODE_CLEARAGAIN         Scancode = 0
	SCANCODE_CRSEL              Scancode = 0
	SCANCODE_EXSEL              Scancode = 0
	SCANCODE_THOUSANDSSEPARATOR Scancode = 0
	SCANCODE_DECIMALSEPARATOR   Scancode = 0
	SCANCODE_CURRENCYUNIT       Scancode = 0
	SCANCODE_CURRENCYSUBUNIT    Scancode = 0
	SCANCODE_LCTRL              Scancode = 17
	SCANCODE_LSHIFT             Scancode = 16
	SCANCODE_LALT               Scancode = 18
	SCANCODE_LGUI               Scancode = 0
	SCANCODE_RCTRL              Scancode = 17
	SCANCODE_RSHIFT             Scancode = 16
	SCANCODE_RALT               Scancode = 18
	SCANCODE_RGUI               Scancode = 0
	SCANCODE_MODE               Scancode = 0
	SCANCODE_AUDIONEXT          Scancode = 0
	SCANCODE_AUDIOPREV          Scancode = 0
	SCANCODE_AUDIOSTOP          Scancode = 0
	SCANCODE_AUDIOPLAY          Scancode = 0
	SCANCODE_AUDIOMUTE          Scancode = 0
	SCANCODE_MEDIASELECT        Scancode = 0

	CONTROLLER_AXIS_INVALID      GameControllerAxis = 0
	CONTROLLER_AXIS_LEFTX        GameControllerAxis = 0
	CONTROLLER_AXIS_LEFTY        GameControllerAxis = 0
	CONTROLLER_AXIS_RIGHTX       GameControllerAxis = 0
	CONTROLLER_AXIS_RIGHTY       GameControllerAxis = 0
	CONTROLLER_AXIS_TRIGGERLEFT  GameControllerAxis = 0
	CONTROLLER_AXIS_TRIGGERRIGHT GameControllerAxis = 0
	CONTROLLER_AXIS_MAX          GameControllerAxis = 0

	CONTROLLER_BUTTON_INVALID       GameControllerButton = -1
	CONTROLLER_BUTTON_A             GameControllerButton = 0
	CONTROLLER_BUTTON_B             GameControllerButton = 0
	CONTROLLER_BUTTON_X             GameControllerButton = 0
	CONTROLLER_BUTTON_Y             GameControllerButton = 0
	CONTROLLER_BUTTON_BACK          GameControllerButton = 0
	CONTROLLER_BUTTON_GUIDE         GameControllerButton = 0
	CONTROLLER_BUTTON_START         GameControllerButton = 0
	CONTROLLER_BUTTON_LEFTSTICK     GameControllerButton = 0
	CONTROLLER_BUTTON_RIGHTSTICK    GameControllerButton = 0
	CONTROLLER_BUTTON_LEFTSHOULDER  GameControllerButton = 0
	CONTROLLER_BUTTON_RIGHTSHOULDER GameControllerButton = 0
	CONTROLLER_BUTTON_DPAD_UP       GameControllerButton = 0
	CONTROLLER_BUTTON_DPAD_DOWN     GameControllerButton = 0
	CONTROLLER_BUTTON_DPAD_LEFT     GameControllerButton = 0
	CONTROLLER_BUTTON_DPAD_RIGHT    GameControllerButton = 0
	CONTROLLER_BUTTON_MAX           GameControllerButton = 0

	CONTROLLERBUTTONDOWN = iota
	CONTROLLERBUTTONUP
	JOYDEVICEADDED
	JOYDEVICEREMOVED
	CONTROLLERDEVICEADDED
	CONTROLLERDEVICEREMOVED
	CONTROLLERDEVICEREMAPPED

	FINGERMOTION
	FINGERDOWN
	FINGERUP

	MOUSEBUTTONDOWN
	MOUSEBUTTONUP

	WINDOWEVENT_NONE
	WINDOWEVENT_ENTER
	WINDOWEVENT_LEAVE
	WINDOWEVENT_SHOWN
	WINDOWEVENT_FOCUS_GAINED
	WINDOWEVENT_HIDDEN
	WINDOWEVENT_FOCUS_LOST
	WINDOWEVENT_RESIZED
	WINDOWEVENT_SIZE_CHANGED
	WINDOWEVENT_CLOSE
)
