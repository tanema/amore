package joystick

import ()

//type Joystick struct {
//id JoystickId
//}

////Checks if a button on the Joystick is pressed.
//func (joystick *Joystick) IsDown() {}

////Gets the direction of each axis.
//func (joystick *Joystick) GetAxes() []float32 {
//return glfw.GetJoystickAxes(glfw.Joystick(joystick.id))
//}

////Gets the direction of an axis.
//func (joystick *Joystick) GetAxis(index int) float32 {
//if index < joystick.GetAxisCount() {
//return joystick.GetAxes()[index]
//}
//return 0.0
//}

////Gets the number of axes on the joystick.
//func (joystick *Joystick) GetAxisCount() int {
//return len(joystick.GetAxes())
//}

////Gets the number of buttons on the joystick.
//func (joystick *Joystick) GetButtonCount() {
//return len(joystick.getButtons())
//}

////Gets the number of buttons on the joystick.
//func (joystick *Joystick) getButtons() []byte {
//return glfw.GetJoystickButtons(glfw.Joystick(joystick.id))
//}

////Gets the joystick's unique identifier.
//func (joystick *Joystick) GetID() JoystickId {
//return joystick.id
//}

////Gets the name of the joystick.
//func (joystick *Joystick) GetName() string {
//return glfw.GetJoystickName(glfw.Joystick(joystick.id))
//}

////Gets whether the Joystick is connected.
//func (joystick *Joystick) IsConnected() bool {
//return glfw.JoystickPresent(glfw.Joystick(joystick.id))
//}

//func (joystick *Joystick) Refresh() {
////is connected
////axis change
////button change
//}
