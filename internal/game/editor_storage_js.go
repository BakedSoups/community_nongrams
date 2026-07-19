//go:build js

package game

import "syscall/js"

const (
	editorPackKey       = "community_nongrams.editor.pack"
	communityLibraryKey = "community_nongrams.community.library"
	communityProfileKey = "community_nongrams.community.profile"
	communityBioKey     = "community_nongrams.community.bio"
	communitySocialKey  = "community_nongrams.community.social"
	communityPaletteKey = "community_nongrams.community.palette"
	communityColorKey   = "community_nongrams.community.favorite_color"
	communityNameKey    = "community_nongrams.community.name"
)

func legacyStorageKey(suffix string) string {
	return "pix" + "aross" + suffix
}

func loadStorageValue(key, legacySuffix string) string {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return ""
	}
	value := storage.Call("getItem", key)
	if !value.IsUndefined() && !value.IsNull() {
		return value.String()
	}
	if legacySuffix == "" {
		return ""
	}
	value = storage.Call("getItem", legacyStorageKey(legacySuffix))
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	text := value.String()
	storage.Call("setItem", key, text)
	return text
}

func saveCommunityProfile(raw string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communityProfileKey, raw)
	return true
}

func saveCommunityBio(bio string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communityBioKey, bio)
	return true
}

func loadCommunityBio() string {
	return loadStorageValue(communityBioKey, ".community.bio")
}

func saveCommunitySocial(social string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communitySocialKey, social)
	return true
}

func loadCommunitySocial() string {
	return loadStorageValue(communitySocialKey, ".community.social")
}

func saveCommunityPalette(palette string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communityPaletteKey, palette)
	return true
}

func loadCommunityPalette() string {
	return loadStorageValue(communityPaletteKey, ".community.palette")
}

func saveCommunityFavoriteColor(color string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communityColorKey, color)
	return true
}

func loadCommunityFavoriteColor() string {
	return loadStorageValue(communityColorKey, ".community.favorite_color")
}

func saveCommunityName(name string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communityNameKey, name)
	return true
}

func loadCommunityName() string {
	if value := loadStorageValue(communityNameKey, ".community.name"); value != "" {
		return value
	}
	fn := js.Global().Get("communityDisplayName")
	if !fn.IsUndefined() && !fn.IsNull() {
		return fn.Invoke().String()
	}
	return ""
}

func takeTextPaste() string {
	fn := js.Global().Get("takeTextPaste")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func loadCommunityProfile() string {
	return loadStorageValue(communityProfileKey, ".community.profile")
}

func communityAccountLabel() string {
	fn := js.Global().Get("communityAccountLabel")
	if fn.IsUndefined() || fn.IsNull() {
		return "Sign in"
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return "Sign in"
	}
	return value.String()
}

func communitySignedIn() bool {
	fn := js.Global().Get("communitySignedIn")
	return !fn.IsUndefined() && !fn.IsNull() && fn.Invoke().Bool()
}

func saveCommunityData(raw string) bool {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}
	storage.Call("setItem", communityLibraryKey, raw)
	return true
}

func loadCommunityData() string {
	return loadStorageValue(communityLibraryKey, ".community.library")
}

func requestCommunityImport(size int, vertical bool) bool {
	fn := js.Global().Get("requestCommunityImport")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(size, vertical)
	return true
}

func takeCommunityImport() string {
	fn := js.Global().Get("takeCommunityImport")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func requestCommunitySignIn(email string) bool {
	fn := js.Global().Get("requestCommunitySignIn")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(email)
	return true
}

func requestCommunitySignOut() bool {
	fn := js.Global().Get("requestCommunitySignOut")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func requestCommunityGoogleSignIn() bool {
	fn := js.Global().Get("requestCommunityGoogleSignIn")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func requestCommunityPublish(raw string, submitOfficial, rightsConfirmed bool, preview string) bool {
	fn := js.Global().Get("requestCommunityPublish")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(raw, submitOfficial, rightsConfirmed, preview)
	return true
}

func requestCommunityPackPublish(raw, preview string) bool {
	fn := js.Global().Get("requestCommunityPackPublish")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(raw, preview)
	return true
}

func takeCommunityResult() string {
	fn := js.Global().Get("takeCommunityResult")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func takeCommunityPublishedID() string {
	fn := js.Global().Get("takeCommunityPublishedID")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func takeCommunityPublishedPackID() string {
	fn := js.Global().Get("takeCommunityPublishedPackID")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func requestCommunityCatalog(kind string) bool {
	fn := js.Global().Get("requestCommunityCatalog")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind)
	return true
}

func takeCommunityCatalog() string {
	fn := js.Global().Get("takeCommunityCatalog")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func syncCommunityDraft(raw string) {
	fn := js.Global().Get("syncCommunityDraft")
	if !fn.IsUndefined() && !fn.IsNull() {
		fn.Invoke(raw)
	}
}

func requestCommunityCloudDrafts() bool {
	fn := js.Global().Get("requestCommunityCloudDrafts")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func takeCommunityCloudDrafts() string {
	fn := js.Global().Get("takeCommunityCloudDrafts")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func requestCommunityCreators() bool {
	fn := js.Global().Get("requestCommunityCreators")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func takeCommunityCreators() string {
	fn := js.Global().Get("takeCommunityCreators")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func syncCommunityProfile(raw, bio, name, social, palette, favoriteColor string) {
	fn := js.Global().Get("syncCommunityProfile")
	if !fn.IsUndefined() && !fn.IsNull() {
		fn.Invoke(raw, bio, name, social, palette, favoriteColor)
	}
}

func deleteCommunityCloudDraft(id string) {
	fn := js.Global().Get("deleteCommunityCloudDraft")
	if !fn.IsUndefined() && !fn.IsNull() {
		fn.Invoke(id)
	}
}

func requestCommunityGallery(kind, sort string) bool {
	fn := js.Global().Get("requestCommunityGallery")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, sort)
	return true
}

func takeCommunityGallery() string {
	fn := js.Global().Get("takeCommunityGallery")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func requestCommunityChat(kind, id string) bool {
	fn := js.Global().Get("requestCommunityChat")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, id)
	return true
}

func takeCommunityChat() string {
	fn := js.Global().Get("takeCommunityChat")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func postCommunityChat(kind, id, body string) bool {
	fn := js.Global().Get("postCommunityChat")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, id, body)
	return true
}

func recordCommunityPlay(levelID string, completed bool) {
	fn := js.Global().Get("recordCommunityPlay")
	if !fn.IsUndefined() && !fn.IsNull() {
		fn.Invoke(levelID, completed)
	}
}

func requestCommunityCompleted() bool {
	fn := js.Global().Get("requestCommunityCompleted")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func takeCommunityCompleted() string {
	fn := js.Global().Get("takeCommunityCompleted")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func requestCommunityPublished() bool {
	fn := js.Global().Get("requestCommunityPublished")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke()
	return true
}

func takeCommunityPublished() string {
	fn := js.Global().Get("takeCommunityPublished")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func unpublishCommunityItem(kind, id string) bool {
	fn := js.Global().Get("unpublishCommunityItem")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, id)
	return true
}

func unpublishCommunityLocalArt(id string) bool {
	fn := js.Global().Get("unpublishCommunityLocalArt")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(id)
	return true
}

func updateCommunityPublishedItem(kind, id, title, description, levelsRaw string) bool {
	fn := js.Global().Get("updateCommunityPublishedItem")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, id, title, description, levelsRaw)
	return true
}

func toggleCommunityLike(kind, id string) bool {
	fn := js.Global().Get("toggleCommunityLike")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, id)
	return true
}

func promoteCommunityItem(kind, id string) bool {
	fn := js.Global().Get("promoteCommunityItem")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(kind, id)
	return true
}

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

func exportEditorImage(filename, raw string) bool {
	fn := js.Global().Get("downloadEditorImage")
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

func requestEditorColorPicker(initial string) bool {
	fn := js.Global().Get("requestEditorColorPicker")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(initial)
	return true
}

func takeEditorColorPicker() string {
	fn := js.Global().Get("takeEditorColorPicker")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func requestEditorTitle(current string) bool {
	fn := js.Global().Get("requestEditorTitle")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(current)
	return true
}

func takeEditorTitle() string {
	fn := js.Global().Get("takeEditorTitle")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func clearEditorTitle() {
	fn := js.Global().Get("clearEditorTitle")
	if !fn.IsUndefined() && !fn.IsNull() {
		fn.Invoke()
	}
}

func requestCommunityCoverImport(size int) bool {
	fn := js.Global().Get("requestCommunityCoverImport")
	if fn.IsUndefined() || fn.IsNull() {
		return false
	}
	fn.Invoke(size)
	return true
}

func takeCommunityCoverImport() string {
	fn := js.Global().Get("takeCommunityCoverImport")
	if fn.IsUndefined() || fn.IsNull() {
		return ""
	}
	value := fn.Invoke()
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
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
