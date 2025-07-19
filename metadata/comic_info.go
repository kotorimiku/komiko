package metadata

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ComicInfo struct {
	XMLName   xml.Name `xml:"ComicInfo"`
	Series    string   `xml:"Series"`
	Writer    string   `xml:"Writer"`
	Publisher string   `xml:"Publisher,omitempty"`
	Genre     string   `xml:"Genre"`
	Summary   string   `xml:"Summary"`
	Title     string   `xml:"Title"`
	Number    *uint    `xml:"Number,omitempty"`
	Volume    *uint    `xml:"Volume,omitempty"`
	PageCount uint     `xml:"PageCount"`
}

func (c *ComicInfo) Build(outputPath string) error {
	fileName := filepath.Join(outputPath, "ComicInfo.xml")

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("error serializing to XML: %v", err)
	}

	return nil
}

func ParseComicInfo(xmlData string) (*ComicInfo, error) {
	var info ComicInfo
	err := xml.Unmarshal([]byte(xmlData), &info)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML: %w", err)
	}

	return &info, nil
}

func ReadComicInfo(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "ComicInfo.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			content, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			return string(content), nil
		}
	}
	return "", nil
}

func LoadComicInfo(filePath string) (*ComicInfo, error) {
	xmlData, err := ReadComicInfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("read ComicInfo.xml failed: %w", err)
	}
	if xmlData == "" {
		return nil, nil
	}

	comicInfo, err := ParseComicInfo(xmlData)
	if err != nil {
		return nil, fmt.Errorf("error parsing ComicInfo: %w", err)
	}

	return comicInfo, nil
}
