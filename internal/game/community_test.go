package game

import "testing"

func TestNormalizeProfileSocialRejectsLinks(t *testing.T) {
	rejected := []string{
		"https://example.com/me",
		"http://example.com/me",
		"www.example.com/me",
		"instagram.com/name",
		"bsky.app\\profile\\name",
	}
	for _, value := range rejected {
		if got, ok := normalizeProfileSocial(value); ok {
			t.Fatalf("normalizeProfileSocial(%q) = %q, true; want rejection", value, got)
		}
	}
}

func TestNormalizeProfileSocialAllowsHandles(t *testing.T) {
	got, ok := normalizeProfileSocial("  instagram   @pixaross  ")
	if !ok {
		t.Fatal("handle was rejected")
	}
	if got != "instagram @pixaross" {
		t.Fatalf("handle = %q, want normalized text", got)
	}
}

func TestLoadCommunityChat(t *testing.T) {
	var game Game
	raw := `[{"id":"m1","authorId":"u1","authorName":"Alex","body":"nice puzzle","createdAt":"2026-07-19T00:00:00Z"}]`
	if err := game.loadCommunityChat(raw); err != nil {
		t.Fatal(err)
	}
	if len(game.communityChatMessages) != 1 {
		t.Fatalf("messages = %d, want 1", len(game.communityChatMessages))
	}
	if game.communityChatMessages[0].Body != "nice puzzle" {
		t.Fatalf("body = %q", game.communityChatMessages[0].Body)
	}
}

func TestCommunityChatBackReturnsToPreviousView(t *testing.T) {
	game := Game{communityView: communityChat, chatReturn: communityGalleryPack}
	game.communityBack()
	if game.communityView != communityGalleryPack {
		t.Fatalf("communityView = %v, want gallery pack", game.communityView)
	}
}

func TestEditorTitleDialogSaveUpdatesDraftTitle(t *testing.T) {
	game := Game{editor: newEditorState(8)}
	game.editor.Title = "Old"
	game.openEditorTitleDialog()
	game.editorTitleDraft = "  New Title  "
	game.closeEditorTitleDialog(true)
	if game.editorTitleEditing {
		t.Fatal("title dialog stayed open")
	}
	if game.editor.Title != "New Title" {
		t.Fatalf("title = %q, want New Title", game.editor.Title)
	}
}
