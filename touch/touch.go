// The touch Pacakge handles touch events in the gl context
// To capture events as they happen you can use the callbacks OnTouchPress,
// OnTouchRelease, and OnTouchMove. Define them by calling touch.OnTouchPress =
// func(x, y, dx, dy, pressure float32){}.
package touch

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Touch struct {
	ID int64
	*sdl.TouchFingerEvent
}

var touches = make(map[int64]*sdl.TouchFingerEvent)

// GetTouches return a slice of all the current touches
func GetTouches() []Touch {
	fingers := []Touch{}
	for id, touch := range touches {
		fingers = append(fingers, Touch{
			ID:               id,
			TouchFingerEvent: touch,
		})
	}
	return fingers
}

// GetPosition will return the x, y coordinates of the touch with the provied id
//This may cause a panic if the id does not exist. Safer to use GetPosition on a given touch.
func GetPosition(id int32) (float32, float32) {
	return touches[int64(id)].X, touches[int64(id)].Y
}

// GetPressure will return the pressure for the touch with the given id. This may
// cause a panic if the id does not exist. Safer to use GetPressure on a given touch.
func GetPressure(id int32) float32 {
	return touches[int64(id)].Pressure
}

// GetPosition will return the x, y coordinates of the touch
func (touch *Touch) GetPosition() (float32, float32) {
	return touch.X, touch.Y
}

// GetPressure will return the pressure for the touch
func (touch *Touch) GetPressure() float32 {
	return touch.Pressure
}
