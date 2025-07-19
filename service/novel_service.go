package service

import (
	"komiko/filecache"
	"komiko/utils/epub"
)

type NovelService struct {
	BookService *BookService
	FileCache   *filecache.FileCache
}

func NewNovelService(bookService *BookService, fileCache *filecache.FileCache) *NovelService {
	return &NovelService{
		BookService: bookService,
		FileCache:   fileCache,
	}
}

func (s *NovelService) GetOpenedFile(bookID uint) (*filecache.OpenedFile, error) {
	openedFile, ok := s.FileCache.Get(bookID)
	if !ok {
		book, err := s.BookService.repo.GetByID(bookID)
		if err != nil {
			return nil, err
		}
		path := book.Path
		file, err := epub.BuildEpubIndex(path)
		if err != nil {
			return nil, err
		}

		openedFile = filecache.NewOpenedFile(bookID, path, file)
		s.FileCache.Set(bookID, openedFile)
	}

	return openedFile, nil
}

func (s *NovelService) GetByPage(bookID uint, page uint) ([]byte, error) {
	openedFile, err := s.GetOpenedFile(bookID)
	if err != nil {
		return nil, err
	}

	return openedFile.Read(page)
}

func (s *NovelService) GetByPath(bookID uint, path string) ([]byte, error) {
	openedFile, err := s.GetOpenedFile(bookID)
	if err != nil {
		return nil, err
	}

	return openedFile.ReadByPath(path)
}

func (s *NovelService) GetChapters(bookID uint) ([]string, error) {
	openedFile, err := s.GetOpenedFile(bookID)
	if err != nil {
		return nil, err
	}

	return openedFile.File.Pages(), nil
}
