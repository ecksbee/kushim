package librarian

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"strconv"

	"ecksbee.com/telefacts/pkg/hydratables"
	"ecksbee.com/telefacts/pkg/serializables"
	"golang.org/x/net/html/charset"
)

type EmbeddedLinkbaseFile struct {
	XMLName  xml.Name   `xml:"schema"`
	XMLAttrs []xml.Attr `xml:",any,attr"`
	RoleRef  []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
	} `xml:"roleRef"`
	Annotation []struct {
		XMLName    xml.Name
		XMLAttrs   []xml.Attr `xml:",any,attr"`
		Definition []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
			CharData string     `xml:",chardata"`
		} `xml:"definition"`
		Appinfo []struct {
			XMLName    xml.Name
			XMLAttrs   []xml.Attr `xml:",any,attr"`
			Definition []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr `xml:",any,attr"`
				CharData string     `xml:",chardata"`
			} `xml:"definition"`
			Linkbase []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr `xml:",any,attr"`
				RoleRef  []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr `xml:",any,attr"`
				} `xml:"roleRef"`

				PresentationLink []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr `xml:",any,attr"`
					Loc      []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"loc"`
					PresentationArc []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"presentationArc"`
				} `xml:"presentationLink"`

				DefinitionLink []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr `xml:",any,attr"`
					Loc      []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"loc"`
					DefinitionArc []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"definitionArc"`
				} `xml:"definitionLink"`

				CalculationLink []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr `xml:",any,attr"`
					Loc      []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"loc"`
					CalculationArc []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"calculationArc"`
				} `xml:"calculationLink"`

				LabelLink []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr `xml:",any,attr"`
					Loc      []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"loc"`
					Label []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
						CharData string     `xml:",chardata"`
					} `xml:"label"`
					LabelArc []struct {
						XMLName  xml.Name
						XMLAttrs []xml.Attr `xml:",any,attr"`
					} `xml:"labelArc"`
				} `xml:"labelLink"`
			} `xml:"linkbase"`
		} `xml:"appinfo"`
	} `xml:"annotation"`
}

func processEmbeddedLinkbaseFile(filepath string,
	schemaFile *serializables.SchemaFile, entry string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	decoded, err := decodeEmbeddedLinkbaseFile(data)
	if err != nil {
		return
	}
	for i, annotation := range decoded.Annotation {
		presentationLinkbase := serializables.PresentationLinkbaseFile{
			RoleRef: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
			}{},
			PresentationLink: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
				Loc      []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"loc\""
				PresentationArc []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"presentationArc\""
			}{},
		}
		definitionLinkbase := serializables.DefinitionLinkbaseFile{
			RoleRef: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
			}{},
			DefinitionLink: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
				Loc      []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"loc\""
				DefinitionArc []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"definitionArc\""
			}{},
		}
		calculationLinkbase := serializables.CalculationLinkbaseFile{
			RoleRef: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
			}{},
			CalculationLink: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
				Loc      []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"loc\""
				CalculationArc []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"calculationArc\""
			}{},
		}
		labelLinkbase := serializables.LabelLinkbaseFile{
			RoleRef: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
			}{},
			LabelLink: []struct {
				XMLName  xml.Name
				XMLAttrs []xml.Attr "xml:\",any,attr\""
				Loc      []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"loc\""
				Label []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
					CharData string     "xml:\",chardata\""
				} "xml:\"label\""
				LabelArc []struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				} "xml:\"labelArc\""
			}{},
		}
		for j, appInfo := range annotation.Appinfo {
			for k, linkbase := range appInfo.Linkbase {
				presentationLinkbase.PresentationLink = append(presentationLinkbase.PresentationLink, linkbase.PresentationLink...)
				definitionLinkbase.DefinitionLink = append(definitionLinkbase.DefinitionLink, linkbase.DefinitionLink...)
				calculationLinkbase.CalculationLink = append(calculationLinkbase.CalculationLink, linkbase.CalculationLink...)
				labelLinkbase.LabelLink = append(labelLinkbase.LabelLink, linkbase.LabelLink...)
				subentry := strconv.Itoa(i) + "_" + strconv.Itoa(j) + "_" + strconv.Itoa(k) + "_" + entry
				//todo roleref
				// presentationLinkbase.RoleRef = append(presentationLinkbase.RoleRef, linkbase.RoleRef...)
				presentationLinkbaseHydrated, err := hydratables.HydratePresentationLinkbase(&presentationLinkbase, entry)
				if err == nil {
					lock.Lock()
					superH.PresentationLinkbases[subentry] = *presentationLinkbaseHydrated
					lock.Unlock()
				}
				// definitionLinkbase.RoleRef = append(definitionLinkbase.RoleRef, linkbase.RoleRef...)
				defintionLinkbaseHydrated, err := hydratables.HydrateDefinitionLinkbase(&definitionLinkbase, entry)
				if err == nil {
					lock.Lock()
					superH.DefinitionLinkbases[subentry] = *defintionLinkbaseHydrated
					lock.Unlock()
				}
				// calculationLinkbase.RoleRef = append(calculationLinkbase.RoleRef, linkbase.RoleRef...)
				calculationLinkbaseHydrated, err := hydratables.HydrateCalculationLinkbase(&calculationLinkbase, entry)
				if err == nil {
					lock.Lock()
					superH.CalculationLinkbases[subentry] = *calculationLinkbaseHydrated
					lock.Unlock()
				}
				// labelLinkbase.RoleRef = append(labelLinkbase.RoleRef, linkbase.RoleRef...)
				labelLinkbaseHydrated, err := hydratables.HydrateLabelLinkbase(&labelLinkbase, entry)
				if err == nil {
					lock.Lock()
					superH.LabelLinkbases[subentry] = *labelLinkbaseHydrated
					lock.Unlock()
				}
			}
		}
	}
}

func decodeEmbeddedLinkbaseFile(xmlData []byte) (*EmbeddedLinkbaseFile, error) {
	reader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	decoded := EmbeddedLinkbaseFile{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}
