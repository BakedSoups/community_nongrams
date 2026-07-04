//go:build js

package game

import "syscall/js"

const editorPackKey = "pixaross.editor.pack"

func saveEditorPack(raw string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", editorPackKey, raw)
	return true
}

func loadEditorPack() string {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return ""
	}
	raw := storage.Call("getItem", editorPackKey)
	if raw.IsUndefined() || raw.IsNull() {
		return ""
	}
	return raw.String()
}

func exportEditorPack(filename, raw string) bool {
	fn := js.Global().Get("downloadEditorPack")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(filename, raw)
	return true
}

func requestEditorImageImport(size int) bool {
	fn := js.Global().Get("requestEditorImageImport")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(size)
	return true
}

func takeEditorImageImport() string {
	fn := js.Global().Get("takeEditorImageImport")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	raw := fn.Invoke()
	if raw.IsUndefined() || raw.IsNull() {
		return ""
	}
	return raw.String()
}

func requestEditorPackImport() bool {
	fn := js.Global().Get("requestEditorPackImport")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func takeEditorPackImport() string {
	fn := js.Global().Get("takeEditorPackImport")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	raw := fn.Invoke()
	if raw.IsUndefined() || raw.IsNull() {
		return ""
	}
	return raw.String()
}

func communityFetchStatus() string {
	fn := js.Global().Get("communityFetchStatus")
	if fn.IsUndefined() || fn.IsNull() {
		return "Supabase not configured"
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return "Supabase not configured"
	}
	return value.String()
}
