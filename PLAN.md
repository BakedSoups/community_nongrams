# Community Creation Plan

## Core Decisions

- A single artwork is the fundamental playable level.
- Packs are ordered collections of versioned levels.
- Guests can draw, playtest, export JPEGs, and save locally.
- Login is required for cloud saving, publishing, and official-game submission.
- Structured level data remains canonical; JPEG is only a personal export or thumbnail.
- Published versions are immutable.

## Phase 1: Community Structure

- Move the Editor out of the main menu.
- Add Community navigation for Featured, Levels, Packs, Create, and My Art.
- Open the editor through Community -> Create.
- Allow guests to use the editor without creating an account.
- Preserve local drafts in browser storage.

Success means a guest can create, save, refresh, reopen, playtest, and export artwork.

## Phase 2: Level Format

Each saved level contains:

- Before layer using black or transparency only.
- Full-color After layer.
- Solution generated from visible After pixels.
- Width and height.
- Title and description.
- Palette and tags.
- Creator and version information.
- Visibility and publication status.

Continue supporting `8x8`, `10x10`, `15x15`, and `20x20` puzzles.

## Phase 3: Import Wizard

Add three import modes:

- Single Before/After pair.
- Regular multi-level grid.
- Aseprite PNG plus JSON metadata.

Include an in-app Aseprite guide specifying:

- Horizontal sprite-sheet layout for simple pairs.
- No trim, padding, borders, or smoothing.
- Transparent background.
- Matching Before and After dimensions.
- Black or transparent Before artwork.

For large sheets, let users configure frame size, rows, columns, margins, spacing, and pairing order.

Show detected frames as thumbnails before importing. Allow users to pair, rename, reorder, or remove frames.

Support automatic filename pairing:

```text
cup_before
cup_after
flower_before
flower_after
```

A large sheet creates multiple drafts and can optionally start a new pack.

## Phase 4: Accounts And My Art

Use Supabase Auth with email magic links and Google login.

My Art should contain:

- Local drafts.
- Cloud drafts.
- Published levels.
- Pack-only levels.
- Unlisted levels.
- Packs.
- Official submissions.

After login, offer to move the current local draft into the user's account.

Use row-level security so only the owner can edit drafts.

## Phase 5: Single-Level Publishing

Before publication, require:

- Nonempty After artwork.
- Valid dimensions.
- Black-only Before artwork.
- Completed playtest.
- Title.
- Ownership confirmation.

Publication options:

- Public.
- Pack only.
- Unlisted.

Community cards show the thumbnail, title, creator, dimensions, difficulty, plays, and likes.

## Phase 6: Official Game Submission

Add an optional publication checkbox:

> **Consider this level for the main game**

Display this explanation:

> Submit this level for inspection by the game team. Submission does not guarantee inclusion. We will review its artwork, originality, puzzle quality, and suitability for the official game.

Require confirmation that the user owns the artwork and permits its use if accepted.

Submission statuses:

```text
Submitted
In review
Changes requested
Approved
Declined
```

Create an immutable review submission linked to the published level version.

A server-side function emails the administrator a review-page link. The email should not contain attachments or act as the review database.

The admin page provides:

- Before and After previews.
- Playable puzzle.
- Creator details.
- Approve action.
- Request-changes action.
- Decline action.
- Internal review notes.

Approval adds that exact version to the official catalog and notifies the creator.

## Phase 7: Art Packs

Allow creators to select their levels from My Art and create a pack.

Pack editing includes:

- Title, description, tags, and cover.
- Drag-to-reorder levels.
- Public, unlisted, or draft visibility.
- Pack-only levels.
- Difficulty progression.
- Pack playtesting.

Pack entries reference immutable level versions. A level can appear in multiple packs.

Track per-user pack progress and completion.

## Phase 8: Moderation And Operations

Add:

- Level and pack reporting.
- Submission rate limits.
- Duplicate-submission prevention.
- Creator blocking.
- Admin moderation queue.
- Publication audit history.
- Content ownership policy.
- Thumbnail regeneration.
- Server-side level validation.

## Suggested Database

```text
profiles
drafts
levels
level_versions
packs
pack_versions
pack_items
official_submissions
likes
play_events
reports
notifications
```

## Recommended Delivery Order

1. Community navigation and local creation.
2. Canonical level format and validation.
3. Aseprite and large-sheet importer.
4. Authentication and My Art.
5. Single-level publishing and Community browsing.
6. Official submission and admin review.
7. Packs and saved progression.
8. Moderation, reporting, and operational tooling.
