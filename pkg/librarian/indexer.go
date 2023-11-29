package librarian

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"sync"
	"time"

	"ecksbee.com/kushim/internal/actions"
	myrenderables "ecksbee.com/kushim/pkg/renderables"
	"ecksbee.com/telefacts/pkg/attr"
	"ecksbee.com/telefacts/pkg/cache"
	"ecksbee.com/telefacts/pkg/hydratables"
	"ecksbee.com/telefacts/pkg/renderables"
	"ecksbee.com/telefacts/pkg/serializables"
	gocache "github.com/patrickmn/go-cache"
)

var (
	IndexingMode bool
	entries      []string
	lock         sync.RWMutex
	superH       hydratables.Hydratable
	cards        map[string]myrenderables.ConceptCard
)

func init() {
	cards = make(map[string]myrenderables.ConceptCard)
	entries = make([]string, 0)
	nowisforever := time.Now().UTC().Format("2006-01-02")
	superH = hydratables.Hydratable{
		Folder: &serializables.Folder{
			EntryFileName:         "kushim.xsd",
			Namespaces:            make(map[string]string),
			Instances:             make(map[string]serializables.InstanceFile),
			Schemas:               make(map[string]serializables.SchemaFile),
			PresentationLinkbases: make(map[string]serializables.PresentationLinkbaseFile),
			DefinitionLinkbases:   make(map[string]serializables.DefinitionLinkbaseFile),
			CalculationLinkbases:  make(map[string]serializables.CalculationLinkbaseFile),
			LabelLinkbases:        make(map[string]serializables.LabelLinkbaseFile),
		},
		Instances: map[string]hydratables.Instance{
			"ecksbee_" + nowisforever + ".xml": hydratables.Instance{
				Contexts: []hydratables.Context{
					hydratables.Context{
						ID: "ctx",
						Period: struct {
							Instant  hydratables.Instant
							Duration hydratables.Duration
						}{
							Instant: hydratables.Instant{
								CharData: nowisforever,
							},
						},
						Entity: hydratables.Entity{
							Identifier: struct {
								Scheme   string
								CharData string
							}{
								Scheme:   "https://ecksbee.com",
								CharData: "kushim",
							},
						},
					},
				},
			},
		},
		Schemas:               make(map[string]hydratables.Schema),
		LabelLinkbases:        make(map[string]hydratables.LabelLinkbase),
		PresentationLinkbases: make(map[string]hydratables.PresentationLinkbase),
		DefinitionLinkbases:   make(map[string]hydratables.DefinitionLinkbase),
		CalculationLinkbases:  make(map[string]hydratables.CalculationLinkbase),
	}
}

func BuildIndex(entry string) {
	lock.Lock()
	entries = append(entries, entry)
	lock.Unlock()
}

func ProcessIndex(gts string) {
	appCache := cache.NewCache(false)
	appCache.Set("names.json", map[string]map[string]string{
		"https://ecksbee.com": map[string]string{
			"kushim": "kushim - XBRL Taxonomy Package Manager",
		},
	}, gocache.DefaultExpiration)
	hydratables.InjectCache(appCache)
	for _, entry := range entries {
		processEntry(entry, entry)
	}
	bytes, err := renderables.MarshalCatalog(&superH)
	if err != nil {
		return
	}
	var catalog renderables.Catalog
	json.Unmarshal(bytes, &catalog)
	schemedEntity := stringify(&catalog.Subjects[0].Entity)
	rsetMap := catalog.Networks[schemedEntity]
	for _, slug := range rsetMap {
		bytes2, err := renderables.MarshalRenderable(slug, &superH)
		if err != nil {
			return
		}
		dest := path.Join(gts, slug+".json")
		actions.WriteFile(dest, bytes2)
	}
	for href, card := range cards {
		bytes3, err := json.Marshal(card)
		if err != nil {
			return
		}
		dest := path.Join(gts, url.QueryEscape(href)+".json")
		actions.WriteFile(dest, bytes3)
	}
	dest := path.Join(gts, "_.json")
	data, _ := json.Marshal(catalog)
	actions.WriteFile(dest, data)
}

func processEntry(trueEntry string, virtualEntry string) error {
	urlPath, err := serializables.UrlToFilename(trueEntry)
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
	processEmbeddedLinkbaseFile(urlPath, schemaFile, virtualEntry)
	processSchema(schemaFile, virtualEntry)
	processLinkbaseFileRefs(schemaFile, virtualEntry)
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
			myDir := filepath.Dir(trueEntry)
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				newentry = schemaLocationAttr.Value
			} else {
				newentry = path.Join(myDir, schemaLocationAttr.Value)
			}
			namespaceAttr := attr.FindAttr(iitem.XMLAttrs, "namespace")
			if namespaceAttr == nil || namespaceAttr.Value == "" {
				return
			}
			lock.Lock()
			superH.Folder.Namespaces[namespaceAttr.Value] = schemaLocationAttr.Value
			lock.Unlock()
			err = processEntry(newentry, schemaLocationAttr.Value)
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
			myDir := filepath.Dir(trueEntry)
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				newentry = schemaLocationAttr.Value
			} else {
				newentry = path.Join(myDir, schemaLocationAttr.Value)
			}
			targetNS := attr.FindAttr(iitem.XMLAttrs, "targetNamespace")
			if targetNS == nil || targetNS.Value == "" {
				return
			}
			lock.Lock()
			superH.Folder.Namespaces[targetNS.Value] = schemaLocationAttr.Value
			lock.Unlock()
			err = processEntry(newentry, schemaLocationAttr.Value)
		}(item)
	}
	wg.Wait()
	return err
}

func processLinkbaseFileRefs(schemaFile *serializables.SchemaFile, entry string) {
	for _, annotation := range schemaFile.Annotation {
		if annotation.XMLName.Space != attr.XSD {
			continue
		}
		for _, appinfo := range annotation.Appinfo {
			if appinfo.XMLName.Space != attr.XSD {
				continue
			}
			for _, item := range appinfo.LinkbaseRef {
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
					presentationLinkbase = reroutePresentationLocs(presentationLinkbase, entry, filepath)
					lock.Lock()
					superH.PresentationLinkbases[filepath] = *presentationLinkbase
					superH.Folder.PresentationLinkbases[filepath] = *discoveredPre
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
					definitionLinkbase = rerouteDefinitionLocs(definitionLinkbase, entry, filepath)
					lock.Lock()
					superH.DefinitionLinkbases[filepath] = *definitionLinkbase
					superH.Folder.DefinitionLinkbases[filepath] = *discoveredDef
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
					calculationLinkbase = rerouteCalculationLocs(calculationLinkbase, entry, filepath)
					lock.Lock()
					superH.CalculationLinkbases[filepath] = *calculationLinkbase
					superH.Folder.CalculationLinkbases[filepath] = *discoveredCal
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
					labelLinkbase = rerouteLabelLocs(labelLinkbase, entry, filepath)
					lock.Lock()
					superH.LabelLinkbases[filepath] = *labelLinkbase
					superH.Folder.LabelLinkbases[filepath] = *discoveredLab
					lock.Unlock()
					break
				default:
					break
				}

			}
		}
	}
}

func processSchema(schemaFile *serializables.SchemaFile, entry string) {
	hydratedSchema, err := hydratables.HydrateSchema(schemaFile, entry)
	if err != nil {
		return
	}
	lock.Lock()
	superH.Schemas[entry] = *hydratedSchema
	superH.Folder.Schemas[entry] = *schemaFile
	lock.Unlock()
	for _, hydratedElement := range hydratedSchema.Element {
		processElement(&hydratedElement, entry)
	}
}

func processElement(concept *hydratables.Concept, source string) {
	href := source + "#" + concept.ID
	lock.RLock()
	namespace := superH.Folder.Namespaces[source]
	lock.RUnlock()
	card := myrenderables.ConceptCard{
		Source:            source,
		ID:                concept.ID,
		Namespace:         namespace,
		Name:              concept.XMLName.Local,
		SubstitutionGroup: concept.SubstitutionGroup.Local,
		PeriodType:        concept.PeriodType,
		ItemType:          concept.Type.Local,
		BalanceType:       concept.Balance,
	}
	lock.Lock()
	cards[href] = card
	lock.Unlock()
}
