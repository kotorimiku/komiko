package service

import (
	"komiko/filecache"
	"komiko/imagearchive"
)

type ComicService struct {
	BookService *BookService
	FileCache   *filecache.FileCache
}

func NewComicService(bookService *BookService, openingFiles *filecache.FileCache) *ComicService {
	return &ComicService{
		BookService: bookService,
		FileCache:   openingFiles,
	}
}

func (s *ComicService) GetByPage(bookID uint, page uint) ([]byte, error) {
	openedFile, ok := s.FileCache.Get(bookID)
	if !ok {
		book, err := s.BookService.repo.GetByID(bookID)
		if err != nil {
			return nil, err
		}
		path := book.Path
		file, err := imagearchive.BuildZipIndex(path)
		if err != nil {
			return nil, err
		}

		openedFile = filecache.NewOpenedFile(bookID, path, file)
		s.FileCache.Set(bookID, openedFile)
	}

	return openedFile.Read(page - 1)
}
