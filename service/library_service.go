package service

import (
	"fmt"
	"komiko/config"
	"komiko/model"
	"komiko/repo"
	"os"
	"path/filepath"
	"strings"
)

type LibraryService struct {
	baseService[model.Library, *repo.LibraryRepo]
	repos       *repo.Repo
	bookService *BookService
}

func NewLibraryService(repos *repo.Repo, bookService *BookService) *LibraryService {
	return &LibraryService{baseService: baseService[model.Library, *repo.LibraryRepo]{repo: repos.LibraryRepo}, repos: repos, bookService: bookService}
}

func (s *LibraryService) Create(library *model.Library) error {
	err := s.baseService.Create(library)
	if err != nil {
		return err
	}
	return s.bookService.ScanCreate(library)
}

func (s *LibraryService) ScanUpdate(id uint) error {
	library, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if library == nil {
		return nil
	}
	return s.bookService.ScanUpdate(library)
}

func (s *LibraryService) ScanCreate(id uint) error {
	library, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.bookService.ScanCreate(library)
}

func (s *LibraryService) UpdateCover(id uint) error {
	library, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.bookService.UpdateCoverAll(library)
}

func (s *LibraryService) DeleteByID(id uint) error {
	tasks := GetTaskManager().ListTasks()
	for _, t := range tasks {
		if t.Status == model.TaskPending || t.Status == model.TaskRunning {
			if strings.Contains(t.Name, fmt.Sprintf(" %d", id)) {
				GetTaskManager().StopTask(t.ID)
			}
		}
	}

	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	covers, err := s.bookService.GetCoversByLibraryID(id)
	if err != nil {
		return err
	}
	for _, cover := range covers {
		os.Remove(filepath.Join(config.CoverDir, cover))
	}
	return nil
}
