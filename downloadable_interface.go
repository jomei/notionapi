package notionapi

import "time"

// DownloadableFileBlock is an interface for blocks that can be downloaded
// such as Pdf, FileBlock, and Image
type DownloadableFileBlock interface {
	Block
	GetURL() string
	GetExpiryTime() *time.Time
}

// GetURL implements DownloadableFileBlock interface for PdfBlock
func (b *PdfBlock) GetURL() string {
	if b.Pdf.File != nil {
		return b.Pdf.File.URL
	}
	if b.Pdf.External != nil {
		return b.Pdf.External.URL
	}
	return ""
}

// GetExpiryTime implements DownloadableFileBlock interface for PdfBlock
func (b *PdfBlock) GetExpiryTime() *time.Time {
	if b.Pdf.File != nil {
		return b.Pdf.File.ExpiryTime
	}
	return nil
}

// GetURL implements DownloadableFileBlock interface for FileBlock
func (b *FileBlock) GetURL() string {
	if b.File.File != nil {
		return b.File.File.URL
	}
	if b.File.External != nil {
		return b.File.External.URL
	}
	return ""
}

// GetExpiryTime implements DownloadableFileBlock interface for FileBlock
func (b *FileBlock) GetExpiryTime() *time.Time {
	if b.File.File != nil {
		return b.File.File.ExpiryTime
	}
	return nil
}

// GetURL implements DownloadableFileBlock interface for ImageBlock
func (b *ImageBlock) GetURL() string {
	return b.Image.GetURL()
}

// GetExpiryTime implements DownloadableFileBlock interface for ImageBlock
func (b *ImageBlock) GetExpiryTime() *time.Time {
	if b.Image.File != nil {
		return b.Image.File.ExpiryTime
	}
	return nil
}

// Verify that types implement DownloadableFileBlock interface
var (
	_ DownloadableFileBlock = (*PdfBlock)(nil)
	_ DownloadableFileBlock = (*FileBlock)(nil)
	_ DownloadableFileBlock = (*ImageBlock)(nil)
)
