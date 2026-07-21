package notionapi

import (
	"testing"
	"time"
)

func TestPdfBlockImplementsDownloadableFileBlock(t *testing.T) {
	// Test setup
	now := time.Now()
	pdfBlock := &PdfBlock{
		Pdf: Pdf{
			File: &FileObject{
				URL:        "https://example.com/file.pdf",
				ExpiryTime: &now,
			},
		},
	}

	// Test GetURL
	if url := pdfBlock.GetURL(); url != "https://example.com/file.pdf" {
		t.Errorf("Expected URL to be 'https://example.com/file.pdf', got %s", url)
	}

	// Test GetExpiryTime
	if expiry := pdfBlock.GetExpiryTime(); expiry != &now {
		t.Errorf("Expected expiry time to be %v, got %v", now, expiry)
	}
}

func TestFileBlockImplementsDownloadableFileBlock(t *testing.T) {
	// Test setup
	now := time.Now()
	fileBlock := &FileBlock{
		File: BlockFile{
			File: &FileObject{
				URL:        "https://example.com/file.txt",
				ExpiryTime: &now,
			},
		},
	}

	// Test GetURL
	if url := fileBlock.GetURL(); url != "https://example.com/file.txt" {
		t.Errorf("Expected URL to be 'https://example.com/file.txt', got %s", url)
	}

	// Test GetExpiryTime
	if expiry := fileBlock.GetExpiryTime(); expiry != &now {
		t.Errorf("Expected expiry time to be %v, got %v", now, expiry)
	}
}

func TestImageBlockImplementsDownloadableFileBlock(t *testing.T) {
	// Test setup
	now := time.Now()
	imageBlock := &ImageBlock{
		Image: Image{
			File: &FileObject{
				URL:        "https://example.com/image.jpg",
				ExpiryTime: &now,
			},
		},
	}

	// Test GetURL
	if url := imageBlock.GetURL(); url != "https://example.com/image.jpg" {
		t.Errorf("Expected URL to be 'https://example.com/image.jpg', got %s", url)
	}

	// Test GetExpiryTime
	if expiry := imageBlock.GetExpiryTime(); expiry != &now {
		t.Errorf("Expected expiry time to be %v, got %v", now, expiry)
	}
}

func TestExternalURLCases(t *testing.T) {
	// Test External URLs for each block type
	testCases := []struct {
		name     string
		block    DownloadableFileBlock
		expected string
	}{
		{
			name: "PDF with external URL",
			block: &PdfBlock{
				Pdf: Pdf{
					External: &FileObject{
						URL: "https://external.com/file.pdf",
					},
				},
			},
			expected: "https://external.com/file.pdf",
		},
		{
			name: "File with external URL",
			block: &FileBlock{
				File: BlockFile{
					External: &FileObject{
						URL: "https://external.com/file.txt",
					},
				},
			},
			expected: "https://external.com/file.txt",
		},
		{
			name: "Image with external URL",
			block: &ImageBlock{
				Image: Image{
					External: &FileObject{
						URL: "https://external.com/image.jpg",
					},
				},
			},
			expected: "https://external.com/image.jpg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if url := tc.block.GetURL(); url != tc.expected {
				t.Errorf("Expected URL to be '%s', got '%s'", tc.expected, url)
			}
		})
	}
}
