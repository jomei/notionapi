package notionapi_test

import (
	"encoding/json"
	"testing"

	"github.com/jomei/notionapi"
)

// TestBlockAudioUnmarshal covers issue #183: the Notion API returns "audio"
// blocks, and while the AudioBlock/Audio types already existed, the block type
// was not registered in the decoder, so audio blocks did not decode into an
// *AudioBlock. This asserts the full decode path now yields an *AudioBlock with
// the expected type and URL.
func TestBlockAudioUnmarshal(t *testing.T) {
	data := []byte(`[{
		"object": "block",
		"id": "audio-block-1",
		"type": "audio",
		"audio": {
			"caption": [],
			"type": "external",
			"external": { "url": "https://example.com/sound.mp3" }
		}
	}]`)

	var blocks notionapi.Blocks
	if err := json.Unmarshal(data, &blocks); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if got := blocks[0].GetType(); got != notionapi.BlockTypeAudio {
		t.Fatalf("expected block type %q, got %q", notionapi.BlockTypeAudio, got)
	}

	audio, ok := blocks[0].(*notionapi.AudioBlock)
	if !ok {
		t.Fatalf("expected *notionapi.AudioBlock, got %T", blocks[0])
	}

	if audio.Audio.GetURL() != "https://example.com/sound.mp3" {
		t.Fatalf("unexpected audio URL: %q", audio.Audio.GetURL())
	}
}
