package game

import (
	"testing"

	"github.com/BakedSoups/community_nongrams/internal/community"
	"github.com/BakedSoups/community_nongrams/internal/nonogram"
)

func validPuzzle(id string, size int) *nonogram.Puzzle {
	rows := make([]string, size)
	before := make([][]string, size)
	reveal := make([][]string, size)
	for y := 0; y < size; y++ {
		row := make([]byte, size)
		before[y] = make([]string, size)
		reveal[y] = make([]string, size)
		for x := 0; x < size; x++ {
			if x == y {
				row[x] = '1'
				reveal[y][x] = "#111111FF"
			} else {
				row[x] = '0'
			}
		}
		rows[y] = string(row)
	}
	puzzle := &nonogram.Puzzle{
		ID:          id,
		Title:       id,
		Width:       size,
		Height:      size,
		SolutionRaw: rows,
		SkeletonRaw: before,
		RevealRaw:   reveal,
	}
	_ = puzzle.ParseSolution()
	return puzzle
}

func TestNormalizeProfileSocialRejectsUnknownLinks(t *testing.T) {
	rejected := []string{
		"https://example.com/me",
		"http://example.com/me",
		"www.example.com/me",
		"not-a-real-social.example/name",
	}
	for _, value := range rejected {
		if got, ok := normalizeProfileSocial(value); ok {
			t.Fatalf("normalizeProfileSocial(%q) = %q, true; want rejection", value, got)
		}
	}
}

func TestNormalizeProfileSocialAllowsHandles(t *testing.T) {
	got, ok := normalizeProfileSocial("  instagram   @community_nongrams  ")
	if !ok {
		t.Fatal("handle was rejected")
	}
	if got != "instagram @community_nongrams" {
		t.Fatalf("handle = %q, want normalized text", got)
	}
}

func TestNormalizeProfileSocialAcceptsKnownLinks(t *testing.T) {
	tests := map[string]string{
		"https://github.com/BakedSoups":                           "github: bakedsoups",
		"https://twitter.com/community_nongrams":                  "x: community_nongrams",
		"https://x.com/community_nongrams":                        "x: community_nongrams",
		"https://instagram.com/community_nongrams/":               "instagram: community_nongrams",
		"https://bsky.app/profile/community_nongrams.bsky.social": "bluesky: community_nongrams.bsky.social",
	}
	for value, want := range tests {
		got, ok := normalizeProfileSocial(value)
		if !ok {
			t.Fatalf("normalizeProfileSocial(%q) was rejected", value)
		}
		if got != want {
			t.Fatalf("normalizeProfileSocial(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestNormalizeProfileSocialListCombinesThreeEntries(t *testing.T) {
	got, ok := normalizeProfileSocialList([3]string{
		"https://github.com/BakedSoups",
		"https://x.com/community_nongrams",
		"instagram: community_nongrams",
	})
	if !ok {
		t.Fatal("social list was rejected")
	}
	want := "github: bakedsoups | x: community_nongrams | instagram: community_nongrams"
	if got != want {
		t.Fatalf("social list = %q, want %q", got, want)
	}
}

func TestSplitProfileSocials(t *testing.T) {
	got := splitProfileSocials("github: alex | x: pixel | instagram: art")
	if got != [3]string{"github: alex", "x: pixel", "instagram: art"} {
		t.Fatalf("split socials = %#v", got)
	}
}

func TestAppendAllowedTextSanitizesPaste(t *testing.T) {
	got, changed := appendAllowedText("hi", "\tthere\n世界!", 12, allowPrintableText)
	if !changed {
		t.Fatal("paste did not change the field")
	}
	if got != "hi there !" {
		t.Fatalf("field = %q, want sanitized printable text", got)
	}
}

func TestAppendAllowedTextCapsLength(t *testing.T) {
	got, changed := appendAllowedText("abc", "defgh", 6, allowPrintableText)
	if !changed {
		t.Fatal("paste did not change the field")
	}
	if got != "abcdef" {
		t.Fatalf("field = %q, want max length applied", got)
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

func TestOpenChatAuthorProfileSelectsCreator(t *testing.T) {
	game := Game{
		communityView: communityChat,
		communityCreators: []community.CreatorProfile{
			{ID: "u1", DisplayName: "Alex"},
			{ID: "u2", DisplayName: "Sam"},
		},
	}
	if !game.openChatAuthorProfile("u2") {
		t.Fatal("author profile was not opened")
	}
	if game.communityView != communityCreatorProfile || game.selectedCreator != 1 {
		t.Fatalf("view = %v creator = %d, want profile index 1", game.communityView, game.selectedCreator)
	}
}

func TestMarkCommunityItemUnpublishedClearsLocalArtStatus(t *testing.T) {
	game := Game{
		communityLibrary: community.Library{
			Drafts: []community.LevelDraft{{
				ID:         "art1",
				Status:     community.LevelPublishedStatus,
				Visibility: community.VisibilityPublic,
				Puzzle:     validPuzzle("art1", 8),
			}},
			Packs: []community.Pack{{Items: []community.PackItem{{LevelID: "art1"}}}},
		},
		communityPublished: []community.GalleryItem{{Kind: "art", ID: "cloud-art1", LocalID: "art1"}},
	}
	game.markCommunityItemUnpublished("art", "cloud-art1")
	if game.communityLibrary.Drafts[0].Status != community.LevelDraftStatus {
		t.Fatalf("status = %q, want draft", game.communityLibrary.Drafts[0].Status)
	}
	if game.communityLibrary.Drafts[0].Visibility != community.VisibilityDraft {
		t.Fatalf("visibility = %q, want draft", game.communityLibrary.Drafts[0].Visibility)
	}
	if game.communityLibrary.Drafts[0].ID == "art1" {
		t.Fatal("draft id was not rotated after unpublish")
	}
	if game.communityLibrary.Drafts[0].Puzzle.ID != game.communityLibrary.Drafts[0].ID {
		t.Fatalf("puzzle id = %q, want %q", game.communityLibrary.Drafts[0].Puzzle.ID, game.communityLibrary.Drafts[0].ID)
	}
	if got := game.communityLibrary.Packs[0].Items[0].LevelID; got != game.communityLibrary.Drafts[0].ID {
		t.Fatalf("pack item level id = %q, want %q", got, game.communityLibrary.Drafts[0].ID)
	}
	if len(game.communityPublished) != 0 {
		t.Fatalf("published items = %d, want 0", len(game.communityPublished))
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
