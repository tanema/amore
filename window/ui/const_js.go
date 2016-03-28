// +build js

package ui

const (
	//mouse buttons
	BUTTON_LEFT   MouseButton = 1
	BUTTON_MIDDLE MouseButton = 2
	BUTTON_RIGHT  MouseButton = 3
	BUTTON_X1     MouseButton = 4
	BUTTON_X2     MouseButton = 5

	//keyboard keys
	K_UNKNOWN            Keycode = ""
	K_RETURN             Keycode = ""
	K_ESCAPE             Keycode = ""
	K_BACKSPACE          Keycode = ""
	K_TAB                Keycode = ""
	K_SPACE              Keycode = ""
	K_EXCLAIM            Keycode = ""
	K_QUOTEDBL           Keycode = ""
	K_HASH               Keycode = ""
	K_PERCENT            Keycode = ""
	K_DOLLAR             Keycode = ""
	K_AMPERSAND          Keycode = ""
	K_QUOTE              Keycode = ""
	K_LEFTPAREN          Keycode = ""
	K_RIGHTPAREN         Keycode = ""
	K_ASTERISK           Keycode = ""
	K_PLUS               Keycode = ""
	K_COMMA              Keycode = ""
	K_MINUS              Keycode = ""
	K_PERIOD             Keycode = ""
	K_SLASH              Keycode = ""
	K_0                  Keycode = ""
	K_1                  Keycode = ""
	K_2                  Keycode = ""
	K_3                  Keycode = ""
	K_4                  Keycode = ""
	K_5                  Keycode = ""
	K_6                  Keycode = ""
	K_7                  Keycode = ""
	K_8                  Keycode = ""
	K_9                  Keycode = ""
	K_COLON              Keycode = ""
	K_SEMICOLON          Keycode = ""
	K_LESS               Keycode = ""
	K_EQUALS             Keycode = ""
	K_GREATER            Keycode = ""
	K_QUESTION           Keycode = ""
	K_AT                 Keycode = ""
	K_LEFTBRACKET        Keycode = ""
	K_BACKSLASH          Keycode = ""
	K_RIGHTBRACKET       Keycode = ""
	K_CARET              Keycode = ""
	K_UNDERSCORE         Keycode = ""
	K_BACKQUOTE          Keycode = ""
	K_a                  Keycode = ""
	K_b                  Keycode = ""
	K_c                  Keycode = ""
	K_d                  Keycode = ""
	K_e                  Keycode = ""
	K_f                  Keycode = ""
	K_g                  Keycode = ""
	K_h                  Keycode = ""
	K_i                  Keycode = ""
	K_j                  Keycode = ""
	K_k                  Keycode = ""
	K_l                  Keycode = ""
	K_m                  Keycode = ""
	K_n                  Keycode = ""
	K_o                  Keycode = ""
	K_p                  Keycode = ""
	K_q                  Keycode = ""
	K_r                  Keycode = ""
	K_s                  Keycode = ""
	K_t                  Keycode = ""
	K_u                  Keycode = ""
	K_v                  Keycode = ""
	K_w                  Keycode = ""
	K_x                  Keycode = ""
	K_y                  Keycode = ""
	K_z                  Keycode = ""
	K_CAPSLOCK           Keycode = ""
	K_F1                 Keycode = ""
	K_F2                 Keycode = ""
	K_F3                 Keycode = ""
	K_F4                 Keycode = ""
	K_F5                 Keycode = ""
	K_F6                 Keycode = ""
	K_F7                 Keycode = ""
	K_F8                 Keycode = ""
	K_F9                 Keycode = ""
	K_F10                Keycode = ""
	K_F11                Keycode = ""
	K_F12                Keycode = ""
	K_PRINTSCREEN        Keycode = ""
	K_SCROLLLOCK         Keycode = ""
	K_PAUSE              Keycode = ""
	K_INSERT             Keycode = ""
	K_HOME               Keycode = ""
	K_PAGEUP             Keycode = ""
	K_DELETE             Keycode = ""
	K_END                Keycode = ""
	K_PAGEDOWN           Keycode = ""
	K_RIGHT              Keycode = ""
	K_LEFT               Keycode = ""
	K_DOWN               Keycode = ""
	K_UP                 Keycode = ""
	K_NUMLOCKCLEAR       Keycode = ""
	K_APPLICATION        Keycode = ""
	K_POWER              Keycode = ""
	K_F13                Keycode = ""
	K_F14                Keycode = ""
	K_F15                Keycode = ""
	K_F16                Keycode = ""
	K_F17                Keycode = ""
	K_F18                Keycode = ""
	K_F19                Keycode = ""
	K_F20                Keycode = ""
	K_F21                Keycode = ""
	K_F22                Keycode = ""
	K_F23                Keycode = ""
	K_F24                Keycode = ""
	K_EXECUTE            Keycode = ""
	K_HELP               Keycode = ""
	K_MENU               Keycode = ""
	K_SELECT             Keycode = ""
	K_STOP               Keycode = ""
	K_AGAIN              Keycode = ""
	K_UNDO               Keycode = ""
	K_CUT                Keycode = ""
	K_COPY               Keycode = ""
	K_PASTE              Keycode = ""
	K_FIND               Keycode = ""
	K_MUTE               Keycode = ""
	K_VOLUMEUP           Keycode = ""
	K_VOLUMEDOWN         Keycode = ""
	K_CANCEL             Keycode = ""
	K_CLEAR              Keycode = ""
	K_PRIOR              Keycode = ""
	K_RETURN2            Keycode = ""
	K_SEPARATOR          Keycode = ""
	K_OUT                Keycode = ""
	K_OPER               Keycode = ""
	K_CLEARAGAIN         Keycode = ""
	K_CRSEL              Keycode = ""
	K_EXSEL              Keycode = ""
	K_THOUSANDSSEPARATOR Keycode = ""
	K_DECIMALSEPARATOR   Keycode = ""
	K_CURRENCYUNIT       Keycode = ""
	K_CURRENCYSUBUNIT    Keycode = ""
	K_LCTRL              Keycode = ""
	K_LSHIFT             Keycode = ""
	K_LALT               Keycode = ""
	K_LGUI               Keycode = ""
	K_RCTRL              Keycode = ""
	K_RSHIFT             Keycode = ""
	K_RALT               Keycode = ""
	K_RGUI               Keycode = ""
	K_MODE               Keycode = ""
	K_AUDIONEXT          Keycode = ""
	K_AUDIOPREV          Keycode = ""
	K_AUDIOSTOP          Keycode = ""
	K_AUDIOPLAY          Keycode = ""
	K_AUDIOMUTE          Keycode = ""
	K_MEDIASELECT        Keycode = ""

	//keyboard modifier keys
	KMOD_NONE     Keymod = ""
	KMOD_LSHIFT   Keymod = ""
	KMOD_RSHIFT   Keymod = ""
	KMOD_LCTRL    Keymod = ""
	KMOD_RCTRL    Keymod = ""
	KMOD_LALT     Keymod = ""
	KMOD_RALT     Keymod = ""
	KMOD_LGUI     Keymod = ""
	KMOD_RGUI     Keymod = ""
	KMOD_NUM      Keymod = ""
	KMOD_CAPS     Keymod = ""
	KMOD_MODE     Keymod = ""
	KMOD_CTRL     Keymod = ""
	KMOD_SHIFT    Keymod = ""
	KMOD_ALT      Keymod = ""
	KMOD_GUI      Keymod = ""
	KMOD_RESERVED Keymod = ""

	//keyboard scancodes
	SCANCODE_UNKNOWN            Scancode = 0
	SCANCODE_A                  Scancode = 0
	SCANCODE_B                  Scancode = 0
	SCANCODE_C                  Scancode = 0
	SCANCODE_D                  Scancode = 0
	SCANCODE_E                  Scancode = 0
	SCANCODE_F                  Scancode = 0
	SCANCODE_G                  Scancode = 0
	SCANCODE_H                  Scancode = 0
	SCANCODE_I                  Scancode = 0
	SCANCODE_J                  Scancode = 0
	SCANCODE_K                  Scancode = 0
	SCANCODE_L                  Scancode = 0
	SCANCODE_M                  Scancode = 0
	SCANCODE_N                  Scancode = 0
	SCANCODE_O                  Scancode = 0
	SCANCODE_P                  Scancode = 0
	SCANCODE_Q                  Scancode = 0
	SCANCODE_R                  Scancode = 0
	SCANCODE_S                  Scancode = 0
	SCANCODE_T                  Scancode = 0
	SCANCODE_U                  Scancode = 0
	SCANCODE_V                  Scancode = 0
	SCANCODE_W                  Scancode = 0
	SCANCODE_X                  Scancode = 0
	SCANCODE_Y                  Scancode = 0
	SCANCODE_Z                  Scancode = 0
	SCANCODE_1                  Scancode = 0
	SCANCODE_2                  Scancode = 0
	SCANCODE_3                  Scancode = 0
	SCANCODE_4                  Scancode = 0
	SCANCODE_5                  Scancode = 0
	SCANCODE_6                  Scancode = 0
	SCANCODE_7                  Scancode = 0
	SCANCODE_8                  Scancode = 0
	SCANCODE_9                  Scancode = 0
	SCANCODE_0                  Scancode = 0
	SCANCODE_RETURN             Scancode = 0
	SCANCODE_ESCAPE             Scancode = 0
	SCANCODE_BACKSPACE          Scancode = 0
	SCANCODE_TAB                Scancode = 0
	SCANCODE_SPACE              Scancode = 0
	SCANCODE_MINUS              Scancode = 0
	SCANCODE_EQUALS             Scancode = 0
	SCANCODE_LEFTBRACKET        Scancode = 0
	SCANCODE_RIGHTBRACKET       Scancode = 0
	SCANCODE_BACKSLASH          Scancode = 0
	SCANCODE_NONUSHASH          Scancode = 0
	SCANCODE_SEMICOLON          Scancode = 0
	SCANCODE_APOSTROPHE         Scancode = 0
	SCANCODE_GRAVE              Scancode = 0
	SCANCODE_COMMA              Scancode = 0
	SCANCODE_PERIOD             Scancode = 0
	SCANCODE_SLASH              Scancode = 0
	SCANCODE_CAPSLOCK           Scancode = 0
	SCANCODE_F1                 Scancode = 0
	SCANCODE_F2                 Scancode = 0
	SCANCODE_F3                 Scancode = 0
	SCANCODE_F4                 Scancode = 0
	SCANCODE_F5                 Scancode = 0
	SCANCODE_F6                 Scancode = 0
	SCANCODE_F7                 Scancode = 0
	SCANCODE_F8                 Scancode = 0
	SCANCODE_F9                 Scancode = 0
	SCANCODE_F10                Scancode = 0
	SCANCODE_F11                Scancode = 0
	SCANCODE_F12                Scancode = 0
	SCANCODE_PRINTSCREEN        Scancode = 0
	SCANCODE_SCROLLLOCK         Scancode = 0
	SCANCODE_PAUSE              Scancode = 0
	SCANCODE_INSERT             Scancode = 0
	SCANCODE_HOME               Scancode = 0
	SCANCODE_PAGEUP             Scancode = 0
	SCANCODE_DELETE             Scancode = 0
	SCANCODE_END                Scancode = 0
	SCANCODE_PAGEDOWN           Scancode = 0
	SCANCODE_RIGHT              Scancode = 0
	SCANCODE_LEFT               Scancode = 0
	SCANCODE_DOWN               Scancode = 0
	SCANCODE_UP                 Scancode = 0
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
	SCANCODE_RETURN2            Scancode = 0
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
	SCANCODE_LCTRL              Scancode = 0
	SCANCODE_LSHIFT             Scancode = 0
	SCANCODE_LALT               Scancode = 0
	SCANCODE_LGUI               Scancode = 0
	SCANCODE_RCTRL              Scancode = 0
	SCANCODE_RSHIFT             Scancode = 0
	SCANCODE_RALT               Scancode = 0
	SCANCODE_RGUI               Scancode = 0
	SCANCODE_MODE               Scancode = 0
	SCANCODE_AUDIONEXT          Scancode = 0
	SCANCODE_AUDIOPREV          Scancode = 0
	SCANCODE_AUDIOSTOP          Scancode = 0
	SCANCODE_AUDIOPLAY          Scancode = 0
	SCANCODE_AUDIOMUTE          Scancode = 0
	SCANCODE_MEDIASELECT        Scancode = 0
	NUM_SCANCODES               Scancode = 0

	CONTROLLER_AXIS_INVALID      GameControllerAxis = 0
	CONTROLLER_AXIS_LEFTX        GameControllerAxis = 0
	CONTROLLER_AXIS_LEFTY        GameControllerAxis = 0
	CONTROLLER_AXIS_RIGHTX       GameControllerAxis = 0
	CONTROLLER_AXIS_RIGHTY       GameControllerAxis = 0
	CONTROLLER_AXIS_TRIGGERLEFT  GameControllerAxis = 0
	CONTROLLER_AXIS_TRIGGERRIGHT GameControllerAxis = 0
	CONTROLLER_AXIS_MAX          GameControllerAxis = 0

	CONTROLLER_BUTTON_INVALID       GameControllerButton = 0
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
