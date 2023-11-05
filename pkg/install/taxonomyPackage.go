package install

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"path"

	"ecksbee.com/kushim/internal/actions"
	"ecksbee.com/telefacts/pkg/attr"
	"golang.org/x/net/html/charset"
)

var (
	TAXONOMYPACKAGE string = "taxonomyPackage.xml"
)

func getTaxonomyPackage(unZipFiles []*zip.File) (*zip.File, error) {
	for _, unZipFile := range unZipFiles {
		dir := path.Base(path.Dir(unZipFile.Name))
		basename := path.Base(unZipFile.Name)
		if basename == TAXONOMYPACKAGE && dir == METAINF {
			return unZipFile, nil
		}
	}
	return nil, nil
}

func name(unZipFiles []*zip.File) (string, error) {
	file, err := getTaxonomyPackage(unZipFiles)
	if err != nil {
		return "", err
	}
	unzipped, err := actions.UnzipFile(file)
	if err != nil {
		return "", err
	}
	tp, err := DecodeTPFile(unzipped)
	if err != nil {
		return "", err
	}
	return tp.Name[0].CharData, nil
}

func entries(unZipFiles []*zip.File) ([]string, error) {
	ret := make([]string, 0)
	file, err := getTaxonomyPackage(unZipFiles)
	if err != nil {
		return []string{}, err
	}
	unzipped, err := actions.UnzipFile(file)
	if err != nil {
		return []string{}, err
	}
	tp, err := DecodeTPFile(unzipped)
	if err != nil {
		return []string{}, err
	}
	entryPoints := tp.EntryPoints[0].EntryPoint
	for _, entry := range entryPoints {
		if len(entry.EntryPointDocument) > 0 {
			epDoc := entry.EntryPointDocument[0]
			if len(epDoc.XMLAttrs) > 0 {
				href := attr.FindAttr(epDoc.XMLAttrs, "href")
				if href != nil && href.Value != "" {
					ret = append(ret, href.Value)
				}
			}
		}
	}
	return ret, nil
}

type TPFile struct {
	XMLName    xml.Name   `xml:"taxonomyPackage"`
	XMLAttrs   []xml.Attr `xml:",any,attr"`
	Identifier []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"identifier"`
	Name []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"name"`
	Description []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"description"`
	Version []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"version"`
	License []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"license"`
	Publisher []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"publisher"`
	PublisherURL []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"publisherURL"`
	PublisherCountry []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"publisherCountry"`
	PublisherDate []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"publisherDate"`
	EntryPoints []struct {
		XMLName    xml.Name
		XMLAttrs   []xml.Attr `xml:",any,attr"`
		EntryPoint []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			Name     []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr `xml:",any,attr"`
				CharData string     `xml:",chardata"`
			} `xml:"name"`
			Description []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr `xml:",any,attr"`
				CharData string     `xml:",chardata"`
			} `xml:"description"`
			Version []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr `xml:",any,attr"`
				CharData string     `xml:",chardata"`
			} `xml:"version"`
			EntryPointDocument []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr `xml:",any,attr"`
				CharData string     `xml:",chardata"`
			} `xml:"entryPointDocument"`
		} `xml:"entryPoint"`
	} `xml:"entryPoints"`
}

func DecodeTPFile(xmlData []byte) (*TPFile, error) {
	reader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	decoded := TPFile{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}
