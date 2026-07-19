//go:build js

package game

import "syscall/js"

func setWebMusicMuted(muted bool) {
	fn := js.Global().Get("setMusicMuted")
	if fn.Type() == js.TypeFunction {
		fn.Invoke(muted)
	}
}

func setWebMusicMode(mode string) {
	fn := js.Global().Get("setMusicMode")
	if fn.Type() == js.TypeFunction {
		fn.Invoke(mode)
	}
}

func playWebSFX(name string) {
	fn := js.Global().Get("playSFX")
	if fn.Type() == js.TypeFunction {
		fn.Invoke(name)
	}
}
