package librarian

import (
	"encoding/xml"
	"fmt"
	"path"
	"sync"

	"ecksbee.com/telefacts/pkg/attr"
	"ecksbee.com/telefacts/pkg/hydratables"
	"ecksbee.com/telefacts/pkg/serializables"
)

var (
	IndexingMode bool
	entries      []string
	lock         sync.RWMutex
	superH       hydratables.Hydratable
	// pgridMap     map[string]renderables.PGrid
	// cgridMap     map[string]renderables.CGrid
	// dgridMap     map[string]renderables.DGrid
)

func init() {
	entries = make([]string, 0)
	superH = hydratables.Hydratable{
		Schemas:               make(map[string]hydratables.Schema),
		LabelLinkbases:        make(map[string]hydratables.LabelLinkbase),
		PresentationLinkbases: make(map[string]hydratables.PresentationLinkbase),
		DefinitionLinkbases:   make(map[string]hydratables.DefinitionLinkbase),
		CalculationLinkbases:  make(map[string]hydratables.CalculationLinkbase),
	}
	// pgridMap = make(map[string]renderables.PGrid)
	// cgridMap = make(map[string]renderables.CGrid)
	// dgridMap = make(map[string]renderables.DGrid)
}

func BuildIndex(entry string) {
	lock.Lock()
	entries = append(entries, entry)
	lock.Unlock()
}

func ProcessIndex() {
	for _, entry := range entries {
		processEntry(entry)
	}
}

func processEntry(entry string) error {
	urlPath, err := serializables.UrlToFilename(entry)
	if err != nil {
		return err
	}
	schemaFile, err := serializables.ReadSchemaFile(urlPath)
	if err != nil {
		return err
	}
	if schemaFile == nil {
		return fmt.Errorf("no schema")
	}
	processSchema(schemaFile, entry)
	processLinkbases(schemaFile, entry)
	imports := schemaFile.Import
	var wg sync.WaitGroup
	wg.Add(len(imports))
	for _, item := range imports {
		go func(iitem struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
		}) {
			defer wg.Done()
			schemaLocationAttr := attr.FindAttr(iitem.XMLAttrs, "schemaLocation")
			if schemaLocationAttr == nil || schemaLocationAttr.Value == "" {
				return
			}
			newentry := ""
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				newentry = schemaLocationAttr.Value
			} else {
				return
			}
			err = processEntry(newentry)
		}(item)
	}
	wg.Wait()
	includes := schemaFile.Include
	wg.Add(len(includes))
	for _, item := range includes {
		go func(iitem struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
		}) {
			defer wg.Done()
			schemaLocationAttr := attr.FindAttr(iitem.XMLAttrs, "schemaLocation")
			if schemaLocationAttr == nil || schemaLocationAttr.Value == "" {
				return
			}
			newentry := ""
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				newentry = schemaLocationAttr.Value
			} else {
				return
			}
			err = processEntry(newentry)
		}(item)
	}
	wg.Wait()
	return err
}

func processLinkbases(schemaFile *serializables.SchemaFile, entry string) {
	var wg sync.WaitGroup
	for _, annotation := range schemaFile.Annotation {
		if annotation.XMLName.Space != attr.XSD {
			continue
		}
		for _, appinfo := range annotation.Appinfo {
			if appinfo.XMLName.Space != attr.XSD {
				continue
			}
			for _, iitem := range appinfo.LinkbaseRef {
				wg.Add(1)
				go func(item struct {
					XMLName  xml.Name
					XMLAttrs []xml.Attr "xml:\",any,attr\""
				}) {
					defer wg.Done()
					if item.XMLName.Space != attr.LINK {
						return
					}
					arcroleAttr := attr.FindAttr(item.XMLAttrs, "arcrole")
					if arcroleAttr == nil || arcroleAttr.Name.Space != attr.XLINK || arcroleAttr.Value != attr.LINKARCROLE {
						return
					}
					typeAttr := attr.FindAttr(item.XMLAttrs, "type")
					if typeAttr == nil || typeAttr.Name.Space != attr.XLINK || typeAttr.Value != "simple" {
						return
					}
					roleAttr := attr.FindAttr(item.XMLAttrs, "role")
					if roleAttr == nil || roleAttr.Name.Space != attr.XLINK || roleAttr.Value == "" {
						return
					}
					hrefAttr := attr.FindAttr(item.XMLAttrs, "href")
					if hrefAttr == nil || hrefAttr.Name.Space != attr.XLINK || hrefAttr.Value == "" {
						return
					}
					if attr.IsValidUrl(hrefAttr.Value) {
						go serializables.DiscoverGlobalFile(hrefAttr.Value)
						return
					}
					filepath := path.Join(serializables.GlobalTaxonomySetPath, hrefAttr.Value)
					switch roleAttr.Value {
					case attr.PresentationLinkbaseRef:
						discoveredPre, err := serializables.ReadPresentationLinkbaseFile(filepath)
						if err != nil {
							return
						}
						presentationLinkbase, err := hydratables.HydratePresentationLinkbase(discoveredPre, filepath)
						if err != nil {
							return
						}
						lock.Lock()
						superH.PresentationLinkbases[filepath] = *presentationLinkbase
						lock.Unlock()
						break
					case attr.DefinitionLinkbaseRef:
						discoveredDef, err := serializables.ReadDefinitionLinkbaseFile(filepath)
						if err != nil {
							return
						}
						definitionLinkbase, err := hydratables.HydrateDefinitionLinkbase(discoveredDef, filepath)
						if err != nil {
							return
						}
						lock.Lock()
						superH.DefinitionLinkbases[filepath] = *definitionLinkbase
						lock.Unlock()
						break
					case attr.CalculationLinkbaseRef:
						discoveredCal, err := serializables.ReadCalculationLinkbaseFile(filepath)
						if err != nil {
							return
						}
						calculationLinkbase, err := hydratables.HydrateCalculationLinkbase(discoveredCal, filepath)
						if err != nil {
							return
						}
						lock.Lock()
						superH.CalculationLinkbases[filepath] = *calculationLinkbase
						lock.Unlock()
						break
					case attr.LabelLinkbaseRef:
						discoveredLab, err := serializables.ReadLabelLinkbaseFile(filepath)
						if err != nil {
							return
						}
						labelLinkbase, err := hydratables.HydrateLabelLinkbase(discoveredLab, filepath)
						if err != nil {
							return
						}
						lock.Lock()
						superH.LabelLinkbases[filepath] = *labelLinkbase
						lock.Unlock()
						break
					default:
						break
					}
				}(iitem)
			}
		}
	}
	wg.Wait()
}

func processSchema(schemaFile *serializables.SchemaFile, entry string) {
	hydratedSchema, err := hydratables.HydrateSchema(schemaFile, entry)
	if err != nil {
		return
	}
	lock.Lock()
	superH.Schemas[entry] = *hydratedSchema
	lock.Unlock()
}
