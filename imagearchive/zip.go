package imagearchive

import (
	"archive/zip"
	"komiko/utils"
	"komiko/utils/ziputil"

	"github.com/facette/natsort"
)

type ZipIndex struct {
	ziputil.ZipIndex
}

func BuildZipIndex(zipPath string) (*ZipIndex, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}

	index := &ZipIndex{ZipIndex: ziputil.ZipIndex{FileMap: make(map[string]*zip.File), Names: make([]string, 0, 20), ZipFile: r}}
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && utils.IsImageFile(f.Name) {
			// 解码文件名以处理编码问题
			decodedName := utils.DecodeFileName(f.Name)
			index.Names = append(index.Names, decodedName)
			index.FileMap[decodedName] = f
		}
	}

	natsort.Sort(index.Names)

	return index, nil
}

func (i *ZipIndex) Read(index uint) ([]byte, error) {
	return i.ReadByPath(i.Names[index])
}

func (i *ZipIndex) Pages() []string {
	return i.Names
}

func ImageCount(zipPath string) (int, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	count := 0
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && utils.IsImageFile(f.Name) {
			count++
		}
	}
	return count, nil
}
