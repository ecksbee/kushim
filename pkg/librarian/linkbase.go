package librarian

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"ecksbee.com/telefacts/pkg/attr"
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
	for _, annotation := range decoded.Annotation {
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
		for _, appInfo := range annotation.Appinfo {
			for _, linkbase := range appInfo.Linkbase {
				presentationLinkbase.PresentationLink = append(presentationLinkbase.PresentationLink, linkbase.PresentationLink...)
				definitionLinkbase.DefinitionLink = append(definitionLinkbase.DefinitionLink, linkbase.DefinitionLink...)
				calculationLinkbase.CalculationLink = append(calculationLinkbase.CalculationLink, linkbase.CalculationLink...)
				labelLinkbase.LabelLink = append(labelLinkbase.LabelLink, linkbase.LabelLink...)
				subentry := entry
				presentationLinkbase.RoleRef = append(presentationLinkbase.RoleRef, linkbase.RoleRef...)
				presentationLinkbaseHydrated, err := hydratables.HydratePresentationLinkbase(&presentationLinkbase, entry)
				if err == nil {
					presentationLinkbaseHydrated := reroutePresentationLocs(presentationLinkbaseHydrated, entry, filepath)
					lock.Lock()
					superH.PresentationLinkbases[subentry] = *presentationLinkbaseHydrated
					superH.Folder.PresentationLinkbases[subentry] = presentationLinkbase
					lock.Unlock()
				}
				definitionLinkbase.RoleRef = append(definitionLinkbase.RoleRef, linkbase.RoleRef...)
				defintionLinkbaseHydrated, err := hydratables.HydrateDefinitionLinkbase(&definitionLinkbase, entry)
				if err == nil {
					defintionLinkbaseHydrated := rerouteDefinitionLocs(defintionLinkbaseHydrated, entry, filepath)
					lock.Lock()
					superH.DefinitionLinkbases[subentry] = *defintionLinkbaseHydrated
					superH.Folder.DefinitionLinkbases[subentry] = definitionLinkbase
					lock.Unlock()
				}
				calculationLinkbase.RoleRef = append(calculationLinkbase.RoleRef, linkbase.RoleRef...)
				calculationLinkbaseHydrated, err := hydratables.HydrateCalculationLinkbase(&calculationLinkbase, entry)
				if err == nil {
					calculationLinkbaseHydrated = rerouteCalculationLocs(calculationLinkbaseHydrated, entry, filepath)
					lock.Lock()
					superH.CalculationLinkbases[subentry] = *calculationLinkbaseHydrated
					superH.Folder.CalculationLinkbases[subentry] = calculationLinkbase
					lock.Unlock()
				}
				labelLinkbase.RoleRef = append(labelLinkbase.RoleRef, linkbase.RoleRef...)
				labelLinkbaseHydrated, err := hydratables.HydrateLabelLinkbase(&labelLinkbase, entry)
				if err == nil {
					labelLinkbaseHydrated = rerouteLabelLocs(labelLinkbaseHydrated, entry, filepath)
					lock.Lock()
					superH.LabelLinkbases[subentry] = *labelLinkbaseHydrated
					superH.Folder.LabelLinkbases[subentry] = labelLinkbase
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

func rerouteLocs(oldLocs []hydratables.Loc, entry string, myFilepath string) []hydratables.Loc {
	ret := make([]hydratables.Loc, len(oldLocs))
	dir := filepath.Dir(entry)
	filename := filepath.Base(myFilepath)
	for i, oldLoc := range oldLocs {
		newHref := oldLoc.Href
		if !attr.IsValidUrl(oldLoc.Href) {
			i := strings.IndexRune(oldLoc.Href, '#')
			if i < 0 {
				continue
			}
			base := oldLoc.Href[:i]
			if len(base) <= 0 {
				base = filename
			}
			fragment := oldLoc.Href[i+1:]
			if len(fragment) <= 0 {
				continue
			}
			newBase := path.Join(dir, base)
			newBase = strings.Replace(newBase, "https:/", "https://", -1)
			newBase = strings.Replace(newBase, "http:/", "http://", -1)
			newHref = newBase + "#" + fragment
		}
		ret[i] = hydratables.Loc{
			Label: oldLoc.Label,
			Href:  newHref,
		}
	}
	return ret
}

func reroutePresentationLocs(old *hydratables.PresentationLinkbase, entry string, myFilepath string) *hydratables.PresentationLinkbase {
	newPresentationLinks := make([]hydratables.PresentationLink, len(old.PresentationLinks))
	for i, dlink := range old.PresentationLinks {
		newPresentationLinks[i] = hydratables.PresentationLink{
			Role:             dlink.Role,
			PresentationArcs: dlink.PresentationArcs,
			Locs:             rerouteLocs(dlink.Locs, entry, myFilepath),
		}
	}
	return &hydratables.PresentationLinkbase{
		FileName:          old.FileName,
		RoleRefs:          old.RoleRefs,
		PresentationLinks: newPresentationLinks,
	}
}

func rerouteDefinitionLocs(old *hydratables.DefinitionLinkbase, entry string, myFilepath string) *hydratables.DefinitionLinkbase {
	newDefinitionLinks := make([]hydratables.DefinitionLink, len(old.DefinitionLinks))
	for i, dlink := range old.DefinitionLinks {
		newDefinitionLinks[i] = hydratables.DefinitionLink{
			Role:           dlink.Role,
			DefinitionArcs: dlink.DefinitionArcs,
			Locs:           rerouteLocs(dlink.Locs, entry, myFilepath),
		}
	}
	return &hydratables.DefinitionLinkbase{
		FileName:        old.FileName,
		RoleRefs:        old.RoleRefs,
		DefinitionLinks: newDefinitionLinks,
	}
}

func rerouteCalculationLocs(old *hydratables.CalculationLinkbase, entry string, myFilepath string) *hydratables.CalculationLinkbase {
	newCalculationLinks := make([]hydratables.CalculationLink, len(old.CalculationLinks))
	for i, clink := range old.CalculationLinks {
		newCalculationLinks[i] = hydratables.CalculationLink{
			Role:            clink.Role,
			CalculationArcs: clink.CalculationArcs,
			Locs:            rerouteLocs(clink.Locs, entry, myFilepath),
		}
	}
	return &hydratables.CalculationLinkbase{
		FileName:         old.FileName,
		RoleRefs:         old.RoleRefs,
		CalculationLinks: newCalculationLinks,
	}
}

func rerouteLabelLocs(old *hydratables.LabelLinkbase, entry string, myFilepath string) *hydratables.LabelLinkbase {
	newLabelLink := make([]hydratables.LabelLink, len(old.LabelLink))
	for i, llink := range old.LabelLink {
		newLabelLink[i] = hydratables.LabelLink{
			Role:      llink.Role,
			LabelArcs: llink.LabelArcs,
			Locs:      rerouteLocs(llink.Locs, entry, myFilepath),
		}
	}
	return &hydratables.LabelLinkbase{
		FileName:  old.FileName,
		RoleRefs:  old.RoleRefs,
		LabelLink: newLabelLink,
	}
}
