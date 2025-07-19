package service

import (
	"context"
	"errors"
	"fmt"
	"komiko/config"
	"komiko/imagearchive"
	"komiko/metadata"
	"komiko/model"
	"komiko/repo"
	"komiko/utils"
	"komiko/utils/epub"
	"komiko/utils/ziputil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var semCover = make(chan struct{}, 4)
var semImage = make(chan struct{}, 16)
var semBook = make(chan struct{}, 1)

type BookService struct {
	baseService[model.Book, *repo.BookRepo]
	repos         *repo.Repo
	seriesService *SeriesService
	genreService  *GenreService
	personService *PersonService
	taskService   *TaskManager
}

func NewBookService(repos *repo.Repo, seriesService *SeriesService, genreService *GenreService, personService *PersonService, taskService *TaskManager) *BookService {
	return &BookService{baseService: baseService[model.Book, *repo.BookRepo]{repo: repos.BookRepo}, repos: repos, seriesService: seriesService, genreService: genreService, personService: personService, taskService: taskService}
}

func (s *BookService) GetBySeriesID(seriesID uint) ([]*model.Book, error) {
	return s.repo.GetBySeriesID(seriesID)
}

func (s *BookService) GetBooks(seriesIDStr string) ([]*model.Book, error) {
	if seriesIDStr == "" {
		return s.repo.GetAll()
	}
	seriesID, err := utils.ParamToUint(seriesIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.GetBySeriesID(seriesID)
}

func (s *BookService) GetCoversByLibraryID(id uint) ([]string, error) {
	return s.repo.GetCoversByLibraryID(id)
}

func (s *BookService) UpdateCoverAll(library *model.Library) error {
	books, err := s.repo.GetByLibraryID(library.ID)
	if err != nil {
		return fmt.Errorf("get books by library id error: %w", err)
	}

	for _, book := range books {
		go func(book *model.Book) {
			semCover <- struct{}{}
			defer func() { <-semCover }()
			err := s.UpdateCover(book, s.GetSeriesCoverBookID(books))
			if err != nil {
				log.Printf("update cover error: %v", err)
			}
		}(book)
	}

	return nil
}

func (s *BookService) GetSeriesCoverBookID(books []*model.Book) uint {
	if len(books) == 0 {
		return 0
	} else if len(books) == 1 {
		books[0].Number = 1
	}

	var seriesCoverBook *model.Book

	for _, book := range books {
		if book.Number == 1 {
			seriesCoverBook = book
			break
		}
	}

	if seriesCoverBook == nil {
		seriesCoverBook = books[0]
	}

	return seriesCoverBook.ID
}

func (s *BookService) UpdateCover(book *model.Book, seriesCoverBookID uint) error {
	image := fmt.Sprintf("b%d.webp", book.ID)
	imagePath, err := filepath.Abs(filepath.Join(config.CoverDir, image))
	if err != nil {
		log.Printf("get absolute path error: %v", err)
	}
	err = SaveCoverFromArchive(book.Path, imagePath, book.Images[0].FileName)
	if err != nil {
		log.Printf("save cover error: %v", err)
	}

	s.repo.Update(&model.Book{ID: book.ID, Cover: image})

	if book.ID == seriesCoverBookID {
		s.repos.SeriesRepo.Update(&model.Series{
			ID:    book.SeriesID,
			Cover: image,
		})
	}

	return nil
}

func (s *BookService) ScanUpdate(library *model.Library) error {
	path := library.Path

	s.taskService.AddTask(fmt.Sprintf("scan update %d", library.ID), func(ctx context.Context, task *model.Task) {
		// remove not exist books
		series, err := s.repos.SeriesRepo.GetByLibraryID(library.ID)
		if err != nil {
			task.Error = fmt.Sprintf("get series by library id error: %v", err)
			return
		}
		for _, series := range series {
			select {
			case <-ctx.Done():
				return
			default:
				if !utils.FileExists(series.Dir) {
					s.repos.SeriesRepo.Delete(series)
				}
			}
		}
		book, err := s.repos.BookRepo.GetByLibraryID(library.ID)
		if err != nil {
			task.Error = fmt.Sprintf("get book by library id error: %v", err)
			return
		}
		for _, book := range book {
			select {
			case <-ctx.Done():
				return
			default:
				if !utils.FileExists(book.Path) {
					s.repos.BookRepo.Delete(book)
				}
			}
		}

		err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				if err != nil {
					return err
				}
				if info.IsDir() {
					return s.AddBook(ctx, filePath, library, true)
				}
				return nil
			}
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			task.Error = fmt.Sprintf("scan update error: %v", err)
			return
		}
	})

	return nil
}

func (s *BookService) ScanCreate(library *model.Library) error {
	path := library.Path

	s.taskService.AddTask(fmt.Sprintf("scan create %d", library.ID), func(ctx context.Context, task *model.Task) {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				if err != nil {
					return err
				}
				if info.IsDir() {
					return s.AddBook(ctx, filePath, library, false)
				}
				return nil
			}
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			task.Error = fmt.Sprintf("scan create error: %v", err)
			return
		}
	})

	return nil
}

func (s *BookService) ParseBook(filePath string, library *model.Library, update bool) (*model.Book, error) {

	existingBook, err := s.repo.GetByPath(filePath)
	if err != nil {
		return nil, fmt.Errorf("get book by path %s error: %w", filePath, err)
	}
	if existingBook != nil && !update {
		return existingBook, nil
	}

	series, book, err := s.ParseBookSeries(filePath, library)
	if err != nil {
		return nil, fmt.Errorf("add book series error for %s: %w", filePath, err)
	}
	if book == nil {
		return nil, nil
	}

	book.Series = series
	if existingBook != nil {
		book.ID = existingBook.ID
	}

	return book, nil
}

func (s *BookService) ParseBooksByDir(ctx context.Context, dirPath string, library *model.Library, update bool) ([]*model.Book, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	books := make([]*model.Book, 0, len(files))
	var mu sync.Mutex
	for _, entry := range files {

		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			if entry.IsDir() {
				continue
			}

			filePath := filepath.Join(dirPath, entry.Name())

			semBook <- struct{}{}
			go func() {
				defer func() { <-semBook }()
				book, err := s.ParseBook(filePath, library, update)
				if err != nil {
					log.Printf("parse book error: %v", err)
				}
				if book != nil {
					mu.Lock()
					books = append(books, book)
					mu.Unlock()
				}
			}()
		}
	}
	return books, nil
}

func (s *BookService) AddBook(ctx context.Context, dirPath string, library *model.Library, update bool) error {
	books, err := s.ParseBooksByDir(ctx, dirPath, library, update)
	if err != nil {
		return fmt.Errorf("get books by dir error: %w", err)
	}

	if len(books) == 0 {
		return nil
	}

	var series *model.Series
	for _, book := range books {
		if book.Series != nil {
			series = book.Series
			break
		}
	}

	if series == nil {
		series = &model.Series{
			Title:     filepath.Base(dirPath),
			Dir:       dirPath,
			LibraryID: library.ID,
		}
	}

	if update {
		err = s.repos.SeriesRepo.UpdateOrCreateByDir(series)
	} else {
		err = s.repos.SeriesRepo.FirstOrCreateByDir(series)
	}
	if err != nil {
		return fmt.Errorf("add series error: %w", err)
	}

	for _, book := range books {
		// not 0 is exist
		if book.ID != 0 && !update {
			continue
		}
		book.Series = nil
		book.SeriesID = series.ID
		book.Title = strings.Replace(book.Title, series.Title, "", 1)
		book.Title = strings.Trim(book.Title, "-_ .")

		if update {
			err = s.repo.UpdateOrCreateByPath(book)
		} else {
			// not exist and update is false
			err = s.repo.FirstOrCreateByPath(book)
		}
		if err != nil {
			return fmt.Errorf("add book error: %w", err)
		}

		// save cover
		go func(book *model.Book) {
			semCover <- struct{}{}
			defer func() { <-semCover }()
			err := s.UpdateCover(book, s.GetSeriesCoverBookID(books))
			if err != nil {
				log.Printf("update cover error: %v", err)
			}
		}(book)
	}

	return nil
}

func (s *BookService) ParseBookSeries(filePath string, library *model.Library) (series *model.Series, book *model.Book, err error) {
	switch library.Type {
	case model.Comic:
		return s.ParseComic(filePath, library)
	case model.Novel:
		return s.ParseNovel(filePath, library)
	default:
		return nil, nil, fmt.Errorf("unknown library type: %s", library.Type)
	}
}

func (s *BookService) ParseNovel(filePath string, library *model.Library) (series *model.Series, book *model.Book, err error) {
	if filepath.Ext(filePath) == ".epub" {
		return s.ParseEpub(filePath, library)
	}
	return nil, nil, nil
}

func (s *BookService) ParseEpub(filePath string, library *model.Library) (series *model.Series, book *model.Book, err error) {
	epubIndex, err := epub.BuildEpubIndex(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("build epub index error: %w", err)
	}
	defer epubIndex.Close()

	pages := []string{}

	if library.Type == model.Novel {
		pages = epubIndex.ChapterNames
	}

	images, err := GetEpubImages(epubIndex)
	if err != nil {
		return nil, nil, err
	}

	opt := epubIndex.Opt

	bookTitle := opt.Metadata.Title
	var number *float32
	if opt.Metadata.SeriesIndex != nil {
		n := float32(*opt.Metadata.SeriesIndex)
		number = &n
	}
	seriesTitle := opt.Metadata.Series

	if bookTitle == "" {
		bookTitle = utils.FileName(filePath)
	}
	if seriesTitle == "" {
		seriesTitle = utils.FileName(filepath.Dir(filePath))
	}
	if number == nil {
		num, ok := utils.ExtractVolume(filepath.Base(filePath))
		if ok {
			number = &num
		} else {
			n := float32(-1)
			number = &n
		}
	}

	genres, err := s.genreService.ToGenresFromList(opt.Metadata.Subjects)
	if err != nil {
		return nil, nil, err
	}
	author, err := s.personService.ToPersons(opt.Metadata.Creator)
	if err != nil {
		return nil, nil, err
	}
	publisher, err := s.personService.ToPersons(opt.Metadata.Publisher)
	if err != nil {
		return nil, nil, err
	}

	series = &model.Series{
		Title:       seriesTitle,
		Description: opt.Metadata.Description,
		Dir:         filepath.Dir(filePath),
		Author:      author,
		Genres:      genres,
		Publisher:   publisher,
		LibraryID:   library.ID,
	}

	book = &model.Book{
		Title:       bookTitle,
		Description: opt.Metadata.Description,
		Path:        filePath,
		Number:      *number,
		PageCount:   uint(len(opt.Spine.Itemrefs)),
		Pages:       pages,
		Images:      images,
		Type:        library.Type,
		Author:      author,
		Genres:      genres,
		Publisher:   publisher,
	}

	return series, book, nil
}

func (s *BookService) ParseComic(filePath string, library *model.Library) (series *model.Series, book *model.Book, err error) {
	if filepath.Ext(filePath) == ".cbz" || filepath.Ext(filePath) == ".zip" {
		comicInfo, err := metadata.LoadComicInfo(filePath)
		if err != nil {
			return nil, nil, err
		}

		images, err := GetImages(filePath)
		if err != nil {
			return nil, nil, err
		}

		if comicInfo == nil {
			book = s.ParseBookWithFilePath(filePath, library, images)

			return nil, book, nil
		}

		bookTitle := comicInfo.Title
		var number *float32
		if comicInfo.Number != nil {
			n := float32(*comicInfo.Number)
			number = &n
		} else if comicInfo.Volume != nil {
			n := float32(*comicInfo.Volume)
			number = &n
		}
		seriesTitle := comicInfo.Series

		if bookTitle == "" {
			bookTitle = utils.FileName(filePath)
		}
		if seriesTitle == "" {
			seriesTitle = utils.FileName(filepath.Dir(filePath))
		}
		if number == nil {
			num, ok := utils.ExtractVolume(filepath.Base(filePath))
			if ok {
				number = &num
			} else {
				n := float32(-1)
				number = &n
			}
		}

		genres, err := s.genreService.ToGenres(comicInfo.Genre)
		if err != nil {
			return nil, nil, err
		}
		author, err := s.personService.ToPersons(comicInfo.Writer)
		if err != nil {
			return nil, nil, err
		}
		publisher, err := s.personService.ToPersons(comicInfo.Publisher)
		if err != nil {
			return nil, nil, err
		}

		series = &model.Series{
			Title:       seriesTitle,
			Description: comicInfo.Summary,
			Dir:         filepath.Dir(filePath),
			Author:      author,
			Genres:      genres,
			Publisher:   publisher,
			LibraryID:   library.ID,
		}

		book = &model.Book{
			Title:       bookTitle,
			Description: comicInfo.Summary,
			Path:        filePath,
			Number:      *number,
			PageCount:   uint(len(images)),
			Images:      images,
			Type:        library.Type,
			Author:      author,
			Genres:      genres,
			Publisher:   publisher,
		}

	} else if filepath.Ext(filePath) == ".epub" {
		series, book, err = s.ParseEpub(filePath, library)
		if err != nil {
			return nil, nil, err
		}
	}
	return series, book, nil
}

func (s *BookService) ParseBookWithFilePath(filePath string, library *model.Library, images []*model.Image) *model.Book {
	var book *model.Book
	fileName := utils.FileName(filePath)
	number, ok := utils.ExtractVolume(filepath.Base(filePath))
	if ok {
		book = &model.Book{
			Title:     fileName,
			Path:      filePath,
			Number:    number,
			PageCount: uint(len(images)),
			Images:    images,
			Type:      library.Type,
		}
	} else {
		book = &model.Book{
			Title:     fileName,
			Path:      filePath,
			Number:    -1,
			PageCount: uint(len(images)),
			Images:    images,
			Type:      library.Type,
		}
	}

	return book
}

func SaveCoverFromArchive(filePath string, output string, coverPath string) error {
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	imageReader, zipReader, err := ziputil.GetReaderFromZip(filePath, coverPath)
	if err != nil {
		return fmt.Errorf("get cover from zip error: %w", err)
	}
	defer imageReader.Close()
	defer zipReader.Close()

	image, err := utils.ResizeImageToMinDimensionWebp(imageReader, 300, 400)
	if err != nil {
		return fmt.Errorf("convert image to avif error: %w", err)
	}
	err = os.WriteFile(output, image, 0644)
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}

func GetImages(filePath string) ([]*model.Image, error) {
	zipIndex, err := imagearchive.BuildZipIndex(filePath)
	if err != nil {
		return nil, err
	}
	defer zipIndex.Close()

	imageCh := make(chan *model.Image, len(zipIndex.Names))
	var wg sync.WaitGroup

	for i, n := range zipIndex.Names {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semImage <- struct{}{}
			defer func() { <-semImage }()

			rc, err := zipIndex.FileMap[n].Open()
			if err != nil {
				log.Printf("open file error: %v", err)
			}
			defer rc.Close()

			width, height, err := utils.GetImageSize(rc)
			if err != nil {
				log.Printf("get image size error: %v", err)
			}

			imageCh <- &model.Image{
				FileName:   n,
				Height:     height,
				Width:      width,
				PageNumber: uint(i + 1),
			}
		}()
	}

	wg.Wait()
	close(imageCh)

	images := make([]*model.Image, 0, len(zipIndex.Names))
	for img := range imageCh {
		images = append(images, img)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].PageNumber < images[j].PageNumber
	})

	println(images[0].FileName)

	return images, nil
}

func GetEpubImages(epubIndex *epub.EpubIndex) ([]*model.Image, error) {
	imgUrls := make([]string, 0, 10)
	for _, chapter := range epubIndex.ChapterNames {
		html, err := epubIndex.ReadByPath(chapter)
		if err != nil {
			return nil, err
		}

		re := regexp.MustCompile(`(?i)<img[^>]*src=(?:'([^']+)'|"([^"]+)")`)
		matches := re.FindAllStringSubmatch(string(html), -1)

		for _, match := range matches {
			if match[1] != "" {
				img := utils.RelToAbs(match[1], chapter)
				if utils.Contains(imgUrls, img) {
					continue
				}
				imgUrls = append(imgUrls, img)
			} else if match[2] != "" {
				img := utils.RelToAbs(match[2], chapter)
				if utils.Contains(imgUrls, img) {
					continue
				}
				imgUrls = append(imgUrls, img)
			}
		}
	}

	imageCh := make(chan *model.Image, len(imgUrls))
	var wg sync.WaitGroup

	for i, img := range imgUrls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semImage <- struct{}{}
			defer func() { <-semImage }()

			rc, err := epubIndex.FileMap[img].Open()
			if err != nil {
				log.Printf("open file error: %v", err)
			}
			defer rc.Close()

			width, height, err := utils.GetImageSize(rc)
			if err != nil {
				log.Printf("get image size error: %v", err)
			}

			imageCh <- &model.Image{
				FileName:   img,
				Height:     height,
				Width:      width,
				PageNumber: uint(i + 1),
			}
		}()
	}

	wg.Wait()
	close(imageCh)

	images := make([]*model.Image, 0, len(imgUrls))
	for img := range imageCh {
		images = append(images, img)
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].PageNumber < images[j].PageNumber
	})

	return images, nil
}
