package epub

import (
	"archive/zip"
	"fmt"
	"komiko/utils"
	"komiko/utils/ziputil"
	"regexp"
	"strings"
)

type EpubIndex struct {
	ziputil.ZipIndex

	ImageNames   []string
	ChapterNames []string
	Catalog      string
	Opt          *Package
}

func BuildEpubIndex(filePath string) (*EpubIndex, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}

	imageNames := make([]string, 0, 5)
	zipIndex := &ziputil.ZipIndex{FileMap: make(map[string]*zip.File), Names: make([]string, 0, 20), ZipFile: r}
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			if utils.IsImageFile(f.Name) {
				imageNames = append(imageNames, f.Name)
			}
			zipIndex.Names = append(zipIndex.Names, f.Name)
			zipIndex.FileMap[f.Name] = f
		}
	}

	opt, err := LoadOpt(filePath)
	if err != nil {
		return nil, err
	}

	chapterNames := make([]string, 0, 5)
	for _, itemRef := range opt.Spine.Itemrefs {
		for _, item := range opt.Manifest.Items {
			if item.ID == itemRef.IDRef {
				chapterNames = append(chapterNames, utils.RelToAbs(item.Href, opt.path))
			}
		}
	}

	return &EpubIndex{ZipIndex: *zipIndex, ImageNames: imageNames, ChapterNames: chapterNames, Opt: opt}, nil
}

func (e *EpubIndex) Read(page uint) ([]byte, error) {
	filePath := e.ChapterNames[page-1]
	html, err := e.ReadByPath(filePath)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`src=(?:'([^']+)'|"([^"]+)")`)

	result := re.ReplaceAllFunc(html, func(m []byte) []byte {
		submatches := re.FindSubmatch(m)
		var path string

		if submatches[1] != nil {
			path = string(submatches[1])
		} else {
			path = string(submatches[2])
		}

		newPath := fmt.Sprintf(`src="file/%s"`, utils.RelToAbs(string(path), filePath))
		return []byte(newPath)
	})

	re = regexp.MustCompile(`href=(?:'([^']+)'|"([^"]+)")`)

	result = re.ReplaceAllFunc(result, func(m []byte) []byte {
		submatches := re.FindSubmatch(m)
		var path string

		if submatches[1] != nil {
			path = string(submatches[1])
		} else {
			path = string(submatches[2])
		}

		if strings.HasPrefix(path, "#") {
			return m
		}

		path = utils.RelToAbs(string(path), filePath)

		for i, c := range e.ChapterNames {
			if strings.HasSuffix(c, path) {
				return []byte(fmt.Sprintf(`href="%d"`, i+1))
			}
		}

		newPath := fmt.Sprintf(`href="file/%s"`, path)
		return []byte(newPath)
	})

	return result, nil
}

func (e *EpubIndex) Pages() []string {
	return e.ChapterNames
}
