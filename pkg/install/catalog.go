package install

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"path"
	"strings"

	"ecksbee.com/kushim/internal/actions"
	"ecksbee.com/telefacts/pkg/attr"
	"golang.org/x/net/html/charset"
)

var (
	CATALOG string = "catalog.xml"
)

func getCatalog(unZipFiles []*zip.File) (*zip.File, error) {
	for _, unZipFile := range unZipFiles {
		dir := path.Base(path.Dir(unZipFile.Name))
		basename := path.Base(unZipFile.Name)
		if basename == CATALOG && dir == METAINF {
			return unZipFile, nil
		}
	}
	return nil, nil
}

func remap(unZipFiles []*zip.File) (map[string]string, error) {
	ret := make(map[string]string)
	file, err := getCatalog(unZipFiles)
	if err != nil {
		return make(map[string]string), err
	}
	unzipped, err := actions.UnzipFile(file)
	if err != nil {
		return make(map[string]string), err
	}
	catalog, err := DecodeCatalogFile(unzipped)
	if err != nil {
		return make(map[string]string), err
	}
	baseURI := path.Clean(path.Dir(file.Name) + "/")
	for _, rewrite := range catalog.RewriteURI {
		uriStartString := attr.FindAttr(rewrite.XMLAttrs, "uriStartString")
		if uriStartString == nil {
			continue
		}
		rewritePrefix := attr.FindAttr(rewrite.XMLAttrs, "rewritePrefix")
		if rewritePrefix == nil {
			continue
		}
		currURI := ""
		if path.IsAbs(rewritePrefix.Value) {
			currURI = rewritePrefix.Value
		} else {
			currURI = path.Clean(path.Join(baseURI, rewritePrefix.Value))
		}
		if oldrewrite, found := ret[currURI]; found {
			oldrewriteLevels := len(strings.Split(oldrewrite, "/"))
			newrewriteLevels := len(strings.Split(uriStartString.Value, "/"))
			if newrewriteLevels < oldrewriteLevels {
				continue
			}
		}
		ret[currURI] = uriStartString.Value
	}
	return ret, nil
}

type CatalogFile struct {
	//https://www.oasis-open.org/committees/download.php/14809/xml-catalogs.html
	XMLName    xml.Name   `xml:"catalog"`
	XMLAttrs   []xml.Attr `xml:",any,attr"`
	RewriteURI []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		CharData string     `xml:",chardata"`
	} `xml:"rewriteURI"`
}

func DecodeCatalogFile(xmlData []byte) (*CatalogFile, error) {
	reader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	decoded := CatalogFile{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}
