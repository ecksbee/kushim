package install

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path"

	"ecksbee.com/kushim/internal/actions"
	"ecksbee.com/kushim/pkg/taxonomies"
	"ecksbee.com/kushim/pkg/throttle"
	"ecksbee.com/telefacts/pkg/serializables"
	"golang.org/x/net/html/charset"
)

func InstallSECTaxonomies(gts string) error {
	if gts == "" {
		return fmt.Errorf("empty Global Taxonomy Set path")
	}
	err := os.MkdirAll(gts, 0755)
	if err != nil {
		return err
	}
	throttle.StartSECThrottle()
	taxonomiesFile, err := actions.Scrape("https://www.sec.gov/info/edgar/edgartaxonomies.xml", throttle.Throttle)
	if err != nil {
		return err
	}
	dest := path.Join(gts, "sec-taxonomies.xml")
	err = actions.WriteFile(dest, taxonomiesFile)
	if err == nil {
		decoded, err := DecodeSECTaxonomiesFile(taxonomiesFile)
		if err != nil {
			return err
		}
		for _, loc := range decoded.Loc {
			schemaUrl := loc.Href[0].CharData
			taxonomies.VolumePath = gts
			throttle.Throttle(schemaUrl)
			taxonomies.DiscoverRemoteURL(schemaUrl)
		}
	}
	taxonomies.VolumePath = gts
	serializables.GlobalTaxonomySetPath = gts
	err = DownloadUTR(throttle.Throttle)
	if err != nil {
		return err
	}
	err = DownloadLRR(throttle.Throttle)
	if err != nil {
		return err
	}
	return DownloadDTRs(throttle.Throttle)
}

type SECTaxonomiesFile struct {
	XMLName  xml.Name   `xml:"Erxl"`
	XMLAttrs []xml.Attr `xml:",any,attr"`
	Loc      []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		Family   []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"Family"`
		Version []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"Version"`
		Href []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"Href"`
		AttType []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"AttType"`
		FileTypeName []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"FileTypeName"`
		Elements []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"Elements"`
		Namespace []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"Namespace"`
		Prefix []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"Prefix"`
	} `xml:"Loc"`
}

func DecodeSECTaxonomiesFile(xmlData []byte) (*SECTaxonomiesFile, error) {
	reader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	decoded := SECTaxonomiesFile{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}
