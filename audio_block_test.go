package notionapi_test

import (
	"github.com/jomei/notionapi"
	"testing"
	"time"
)

func TestAudioBlock_GetURL(t *testing.T) {
	tests := []struct {
		name  string
		block *notionapi.AudioBlock
		want  string
	}{
		{
			name: "returns internal file URL",
			block: &notionapi.AudioBlock{
				Audio: notionapi.Audio{
					File: &notionapi.FileObject{
						URL: "https://example.com/internal.mp3",
					},
				},
			},
			want: "https://example.com/internal.mp3",
		},
		{
			name: "returns external file URL",
			block: &notionapi.AudioBlock{
				Audio: notionapi.Audio{
					External: &notionapi.FileObject{
						URL: "https://example.com/external.mp3",
					},
				},
			},
			want: "https://example.com/external.mp3",
		},
		{
			name:  "returns empty string when no URL",
			block: &notionapi.AudioBlock{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.block.GetURL(); got != tt.want {
				t.Errorf("AudioBlock.GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudioBlock_GetExpiryTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		block *notionapi.AudioBlock
		want  *time.Time
	}{
		{
			name: "returns expiry time for internal file",
			block: &notionapi.AudioBlock{
				Audio: notionapi.Audio{
					File: &notionapi.FileObject{
						ExpiryTime: &now,
					},
				},
			},
			want: &now,
		},
		{
			name: "returns nil for external file",
			block: &notionapi.AudioBlock{
				Audio: notionapi.Audio{
					External: &notionapi.FileObject{},
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.block.GetExpiryTime()
			if got != tt.want {
				t.Errorf("AudioBlock.GetExpiryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Compile-time interface check
var _ notionapi.DownloadableFileBlock = (*notionapi.AudioBlock)(nil)
