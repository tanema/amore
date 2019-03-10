package audio

import (
	"time"

	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/runtime"
)

var audioFunctions = runtime.LuaFuncs{
	"new":       audioNewSource,
	"getvolume": audioGetVolume,
	"setvolume": audioSetVolume,
	"pause":     audioPauseAll,
	"play":      audioPlayAll,
	"rewind":    audioRewindAll,
	"stop":      audioStopAll,
}

var audioMetaTables = runtime.LuaMetaTable{
	"Source": {
		"finished":   audioSourceIsFinished,
		"looping":    audioSourceIsLooping,
		"paused":     audioSourceIsPaused,
		"playing":    audioSourceIsPlaying,
		"static":     audioSourceIsStatic,
		"stopped":    audioSourceIsStopped,
		"duration":   audioSourceGetDuration,
		"pitch":      audioSourceGetPitch,
		"volume":     audioSourceGetVolume,
		"state":      audioSourceGetState,
		"setlooping": audioSourceSetLooping,
		"setpitch":   audioSourceSetPitch,
		"setvolume":  audioSourceSetVolume,
		"play":       audioSourcePlay,
		"pause":      audioSourcePause,
		"resume":     audioSourceResume,
		"rewind":     audioSourceRewind,
		"seek":       audioSourceSeek,
		"stop":       audioSourceStop,
		"tell":       audioSourceTell,
	},
}

func init() {
	runtime.RegisterModule("audio", audioFunctions, audioMetaTables)
}

func audioNewSource(ls *lua.LState) int {
	source, err := NewSource(ls.ToString(1), ls.ToBool(2))
	if err != nil {
		return 0
	}
	return returnSource(ls, source)
}

func audioGetVolume(ls *lua.LState) int {
	ls.Push(lua.LNumber(GetVolume()))
	return 1
}

func audioSetVolume(ls *lua.LState) int {
	SetVolume(float32(ls.ToNumber(1)))
	return 0
}

func audioPlayAll(ls *lua.LState) int {
	PlayAll()
	return 0
}

func audioPauseAll(ls *lua.LState) int {
	PlayAll()
	return 0
}

func audioRewindAll(ls *lua.LState) int {
	RewindAll()
	return 0
}

func audioStopAll(ls *lua.LState) int {
	StopAll()
	return 0
}

func audioSourceIsFinished(ls *lua.LState) int {
	ls.Push(lua.LBool(toSource(ls, 1).IsFinished()))
	return 1
}

func audioSourceIsLooping(ls *lua.LState) int {
	ls.Push(lua.LBool(toSource(ls, 1).IsLooping()))
	return 1
}

func audioSourceIsPaused(ls *lua.LState) int {
	ls.Push(lua.LBool(toSource(ls, 1).IsPaused()))
	return 1
}

func audioSourceIsPlaying(ls *lua.LState) int {
	ls.Push(lua.LBool(toSource(ls, 1).IsPlaying()))
	return 1
}

func audioSourceIsStatic(ls *lua.LState) int {
	ls.Push(lua.LBool(toSource(ls, 1).IsStatic()))
	return 1
}

func audioSourceIsStopped(ls *lua.LState) int {
	ls.Push(lua.LBool(toSource(ls, 1).IsStopped()))
	return 1
}

func audioSourceGetDuration(ls *lua.LState) int {
	ls.Push(lua.LNumber(toSource(ls, 1).GetDuration().Seconds()))
	return 1
}

func audioSourceTell(ls *lua.LState) int {
	ls.Push(lua.LNumber(toSource(ls, 1).Tell().Seconds()))
	return 1
}

func audioSourceGetPitch(ls *lua.LState) int {
	ls.Push(lua.LNumber(toSource(ls, 1).GetPitch()))
	return 1
}

func audioSourceGetVolume(ls *lua.LState) int {
	ls.Push(lua.LNumber(toSource(ls, 1).GetVolume()))
	return 1
}

func audioSourceGetState(ls *lua.LState) int {
	ls.Push(lua.LString(toSource(ls, 1).GetState()))
	return 1
}

func audioSourceSetLooping(ls *lua.LState) int {
	toSource(ls, 1).SetLooping(ls.ToBool(2))
	return 0
}

func audioSourceSetPitch(ls *lua.LState) int {
	toSource(ls, 1).SetPitch(float32(ls.ToNumber(2)))
	return 0
}

func audioSourceSetVolume(ls *lua.LState) int {
	toSource(ls, 1).SetVolume(float32(ls.ToNumber(2)))
	return 0
}

func audioSourceSeek(ls *lua.LState) int {
	toSource(ls, 1).Seek(time.Duration(float32(ls.ToNumber(2)) * 1e09))
	return 0
}

func audioSourcePlay(ls *lua.LState) int {
	toSource(ls, 1).Play()
	return 0
}

func audioSourcePause(ls *lua.LState) int {
	toSource(ls, 1).Pause()
	return 0
}

func audioSourceResume(ls *lua.LState) int {
	toSource(ls, 1).Resume()
	return 0
}

func audioSourceRewind(ls *lua.LState) int {
	toSource(ls, 1).Rewind()
	return 0
}

func audioSourceStop(ls *lua.LState) int {
	toSource(ls, 1).Stop()
	return 0
}

func toSource(ls *lua.LState, offset int) Source {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(Source); ok {
		return v
	}
	ls.ArgError(offset, "audio source expected")
	return nil
}

func returnSource(ls *lua.LState, source Source) int {
	f := ls.NewUserData()
	f.Value = source
	ls.SetMetatable(f, ls.GetTypeMetatable("Source"))
	ls.Push(f)
	return 1
}
