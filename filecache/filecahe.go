package filecache

import (
	"sync"
	"time"
)

type FileCache struct {
	Files map[uint]*OpenedFile
	mu    sync.RWMutex
}

func (f *FileCache) Get(bookID uint) (*OpenedFile, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	file, ok := f.Files[bookID]
	if ok {
		file.LastAccess = time.Now()
	}
	return file, ok
}

func (f *FileCache) Set(bookID uint, file *OpenedFile) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Files[bookID] = file
}

func (f *FileCache) Delete(bookID uint) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.Files, bookID)
}

func (f *FileCache) StartCleaner() {
	timeout := 5 * time.Minute
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			now := time.Now()
			var toDelete []uint

			f.mu.RLock()
			for bookID, of := range f.Files {
				if now.Sub(of.LastAccess) > timeout {
					toDelete = append(toDelete, bookID)
				}
			}
			f.mu.RUnlock()

			for _, bookID := range toDelete {
				f.mu.Lock()
				if of, ok := f.Files[bookID]; ok {
					_ = of.Close()
					delete(f.Files, bookID)
				}
				f.mu.Unlock()
			}
		}
	}()
}
