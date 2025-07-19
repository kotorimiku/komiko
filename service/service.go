package service

import (
	"komiko/filecache"
	"komiko/repo"
)

type Service struct {
	BookService     *BookService
	LibraryService  *LibraryService
	SeriesService   *SeriesService
	ComicService    *ComicService
	NovelService    *NovelService
	UserService     *UserService
	ProgressService *ProgressService
	TaskManager     *TaskManager
}

func NewService(repo *repo.Repo) *Service {
	taskManager := GetTaskManager()
	seriesService := NewSeriesService(repo)
	genreService := NewGenreService(repo)
	personService := NewPersonService(repo)
	bookService := NewBookService(repo, seriesService, genreService, personService, taskManager)
	filecache := &filecache.FileCache{Files: make(map[uint]*filecache.OpenedFile)}
	filecache.StartCleaner()
	comicService := NewComicService(bookService, filecache)
	novelService := NewNovelService(bookService, filecache)
	userService := NewUserService(repo)
	progressService := NewProgressService(repo)
	return &Service{
		BookService:     bookService,
		LibraryService:  NewLibraryService(repo, bookService),
		SeriesService:   seriesService,
		ComicService:    comicService,
		NovelService:    novelService,
		UserService:     userService,
		ProgressService: progressService,
		TaskManager:     taskManager,
	}
}
