package filecache

import "time"

type File interface {
	Read(uint) ([]byte, error)
	ReadByPath(string) ([]byte, error)
	Pages() []string
	// ReadByRelPath(string) ([]byte, error)
	Close() error
}

type FileType string

const (
	FileTypeZip FileType = "zip"
	FileTypeRar FileType = "rar"
)

type OpenedFile struct {
	BookID     uint
	Path       string
	File       File
	LastAccess time.Time
}

func NewOpenedFile(bookID uint, path string, file File) *OpenedFile {
	return &OpenedFile{
		BookID:     bookID,
		Path:       path,
		File:       file,
		LastAccess: time.Now(),
	}
}

func (f *OpenedFile) Close() (err error) {
	return f.File.Close()
}

func (f *OpenedFile) Read(index uint) (content []byte, err error) {
	return f.File.Read(index)
}

func (f *OpenedFile) ReadByPath(path string) (content []byte, err error) {
	return f.File.ReadByPath(path)
}

// func (f *OpenedFile) ReadByRelPath(path string) ([]byte, error) {
// 	return f.File.ReadByRelPath(path)
// }
