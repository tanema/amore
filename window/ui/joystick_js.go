// +build js

package ui

func NewJoystick(idx int) *Joystick {
	return &Joystick{id: idx}
}

func (joystick *Joystick) Open() bool                                     { return false }
func (joystick *Joystick) IsGamepad() bool                                { return false }
func (joystick *Joystick) GetID() int                                     { return joystick.id }
func (joystick *Joystick) GetName() string                                { return "" }
func (joystick *Joystick) IsConnected() bool                              { return false }
func (joystick *Joystick) GetGUID() string                                { return "" }
func (joystick *Joystick) Close()                                         {}
func (joystick *Joystick) IsDown(button int) bool                         { return false }
func (joystick *Joystick) GetAxisCount() int                              { return 0 }
func (joystick *Joystick) GetButtonCount() int                            { return 0 }
func (joystick *Joystick) GetHatCount() int                               { return 0 }
func (joystick *Joystick) GetAxis(axisindex int) float32                  { return 0 }
func (joystick *Joystick) GetAxes() []float32                             { return []float32{} }
func (joystick *Joystick) GetHat(hatindex int) byte                       { return 0 }
func (joystick *Joystick) GetGamepadAxis(axis GameControllerAxis) float32 { return 0 }
func (joystick *Joystick) IsGamepadDown(button GameControllerButton) bool { return false }
func (joystick *Joystick) IsVibrationSupported() bool                     { return false }
func (joystick *Joystick) SetVibration(args ...float32) bool              { return false }
func (joystick *Joystick) GetVibration() (float32, float32)               { return 0, 0 }
