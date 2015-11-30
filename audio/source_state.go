package audio

import (
	"github.com/tanema/amore/audio/al"
)

// State indicates the current playing state of the player.
type State int

const (
	Unknown State = iota
	Initial
	Playing
	Paused
	Stopped
)

func (s State) String() string { return stateStrings[s] }

var stateStrings = [...]string{
	Unknown: "unknown",
	Initial: "initial",
	Playing: "playing",
	Paused:  "paused",
	Stopped: "stopped",
}

var codeToState = map[int32]State{
	0:          Unknown,
	al.Initial: Initial,
	al.Playing: Playing,
	al.Paused:  Paused,
	al.Stopped: Stopped,
}
