//go:build !js

package game

func saveEditorPack(string) bool {
	return false
}

func loadEditorPack() string {
	return ""
}

func exportEditorImage(string, string) bool {
	return false
}

func requestEditorImageImport(int) bool {
	return false
}

func takeEditorImageImport() string {
	return ""
}

func requestEditorColorPicker(string) bool {
	return false
}

func takeEditorColorPicker() string {
	return ""
}

func requestEditorPackImport() bool {
	return false
}

func takeEditorPackImport() string {
	return ""
}

func communityFetchStatus() string {
	return "Community is available in the web build"
}
