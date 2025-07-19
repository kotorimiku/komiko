package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"komiko/utils"
)

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
	Guide    Guide    `xml:"guide"`
	Rest     string   `xml:",innerxml"`
	path     string
}

type Meta struct {
	Name     string `xml:"name,attr"`
	Content  string `xml:"content,attr"`
	Property string `xml:"property,attr"`
	Value    string `xml:",chardata"`
}

type Metadata struct {
	Title       string   `xml:"http://purl.org/dc/elements/1.1/ title"`
	Creator     string   `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Description string   `xml:"http://purl.org/dc/elements/1.1/ description"`
	Publisher   string   `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Language    string   `xml:"http://purl.org/dc/elements/1.1/ language"`
	Identifier  string   `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Subjects    []string `xml:"http://purl.org/dc/elements/1.1/ subject"`
	Cover       string   `xml:"meta[name='cover']"`
	Modified    string   `xml:"meta[property='dcterms:modified']"`
	Series      string   `xml:"meta[name='calibre:series']"`
	SeriesIndex *uint    `xml:"meta[name='calibre:series_index']"`
	Metas       []Meta   `xml:"meta"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr"`
}

type Spine struct {
	Itemrefs []Itemref `xml:"itemref"`
}

type Itemref struct {
	IDRef string `xml:"idref,attr"`
}

type Guide struct {
	References []Reference `xml:"reference"`
}

type Reference struct {
	Href  string `xml:"href,attr"`
	Title string `xml:"title,attr"`
	Type  string `xml:"type,attr"`
}

func ReadOpf(epubPath string) (opt string, optPath string, err error) {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return "", "", err
	}
	defer r.Close()

	files := map[string]*zip.File{}
	optPath = ""

	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			if f.Name == "META-INF/container.xml" {
				rc, err := f.Open()
				if err != nil {
					return "", "", err
				}
				defer rc.Close()

				// Read container.xml
				decoder := xml.NewDecoder(rc)
				for {
					tok, err := decoder.Token()
					if err != nil {
						break
					}
					switch se := tok.(type) {
					case xml.StartElement:
						if se.Name.Local == "rootfile" {
							for _, attr := range se.Attr {
								if attr.Name.Local == "full-path" {
									optPath = attr.Value
									file, ok := files[optPath]
									if ok {
										rc, err := file.Open()
										if err != nil {
											return "", "", err
										}
										defer rc.Close()

										content, err := io.ReadAll(rc)
										if err != nil {
											return "", "", err
										}
										return string(content), optPath, nil

									}
								}
							}
						}
					}
				}
			} else if optPath != "" && f.Name == optPath {
				rc, err := f.Open()
				if err != nil {
					return "", "", err
				}
				defer rc.Close()

				content, err := io.ReadAll(rc)
				if err != nil {
					return "", "", err
				}
				return string(content), optPath, nil
			} else {
				files[f.Name] = f
			}
		}

	}
	return "", "", nil
}

func ParseOpt(xmlData string) (*Package, error) {
	var opt Package
	err := xml.Unmarshal([]byte(xmlData), &opt)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML: %w", err)
	}

	for _, item := range opt.Metadata.Metas {
		if item.Name == "calibre:series_index" {
			seriesIndex, err := utils.ParamToUint(item.Content)
			if err != nil {
				return nil, fmt.Errorf("error parsing series_index: %w", err)
			}
			opt.Metadata.SeriesIndex = &seriesIndex
		} else if item.Name == "calibre:series" {
			opt.Metadata.Series = item.Content
		} else if item.Content == "dcterms:modified" {
			opt.Metadata.Modified = item.Value
		} else if item.Name == "cover" {
			opt.Metadata.Cover = item.Content
		}
	}

	return &opt, nil
}

func LoadOpt(filePath string) (*Package, error) {
	xmlData, optPath, err := ReadOpf(filePath)
	if err != nil {
		return nil, fmt.Errorf("read xml failed: %w", err)
	}
	if xmlData == "" {
		return nil, fmt.Errorf("xml file does not exist or is empty")
	}

	opt, err := ParseOpt(xmlData)
	if err != nil {
		return nil, fmt.Errorf("parse xml failed: %w", err)
	}
	opt.path = optPath

	return opt, nil
}
