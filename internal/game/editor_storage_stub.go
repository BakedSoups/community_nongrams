//go:build !js

package game

func saveEditorPack(string) bool {
	return false
}

func loadEditorPack() string {
	return ""
}

func exportEditorPack(string, string) bool {
	return false
}

func requestEditorImageImport(int) bool {
	return false
}

func takeEditorImageImport() string {
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
