package ziputil

import (
	"archive/zip"
	"fmt"
	"io"
	"komiko/utils"
	"path/filepath"
	"strings"
)

type ZipIndex struct {
	FileMap map[string]*zip.File
	Names   []string
	ZipFile *zip.ReadCloser
}

func (zi *ZipIndex) Close() error {
	return zi.ZipFile.Close()
}

func BuildZipIndex(zipPath string) (*ZipIndex, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}

	index := &ZipIndex{FileMap: make(map[string]*zip.File), Names: make([]string, 0, 20), ZipFile: r}
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			decodedName := utils.DecodeFileName(f.Name)
			index.Names = append(index.Names, decodedName)
			index.FileMap[decodedName] = f
		}
	}
	return index, nil
}

func (i *ZipIndex) ReadByPath(path string) ([]byte, error) {
	file, ok := i.FileMap[path]
	if !ok {
		return nil, nil
	}
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (i *ZipIndex) ReadByRelPath(path string, relPath string) ([]byte, error) {
	path = strings.ReplaceAll(filepath.Clean(filepath.Join(filepath.Dir(relPath), path)), "\\", "/")
	return i.ReadByPath(path)
}

func ReadString(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	buf, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func GetFromZip(zipPath string, filePath string) ([]byte, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		// 解码文件名以处理编码问题
		decodedName := utils.DecodeFileName(f.Name)
		if decodedName == filePath || f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			return content, nil
		}
	}
	return nil, nil
}

func GetReaderFromZip(zipPath string, filePath string) (fileReader io.ReadCloser, zipReader *zip.ReadCloser, err error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, nil, err
	}

	for _, f := range r.File {
		// 解码文件名以处理编码问题
		decodedName := utils.DecodeFileName(f.Name)
		if decodedName == filePath || f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, nil, err
			}

			return rc, r, nil
		}
	}

	r.Close()
	return nil, nil, fmt.Errorf("file %s not found in zip %s", filePath, zipPath)
}
