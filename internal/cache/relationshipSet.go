package cache

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"ecksbee.com/kushim/pkg/taxonomies"
	"ecksbee.com/telefacts/pkg/attr"
	"ecksbee.com/telefacts/pkg/serializables"
	gocache "github.com/patrickmn/go-cache"
)

func MarshalRelationshipSet(roleuri string) ([]byte, error) {
	data := rsrepo.Query(roleuri)
	return json.Marshal(data)
}

type RelationshipSetRepo struct {
	lock  sync.RWMutex
	cache *gocache.Cache
}

func (repo *RelationshipSetRepo) Query(roleuri string) *RelationshipSetView {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	if x, found := repo.cache.Get("rs:" + roleuri); found {
		ret := x.(RelationshipSetView)
		return &ret
	}
	return nil
}

func NewRelationshipSetRepo(cache *gocache.Cache, gts string) (*RelationshipSetRepo, error) {
	taxonomiesDir := path.Join(gts, "taxonomies")
	entries, err := os.ReadDir(taxonomiesDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() == false {
			continue
		}
		id := entry.Name()
		underscore := path.Join(taxonomiesDir, id, "_")
		data, err := ioutil.ReadFile(underscore)
		if err != nil {
			return nil, err
		}
		meta := taxonomies.Meta{}
		err = json.Unmarshal(data, &meta)
		if err != nil {
			return nil, err
		}
		metaEntries := meta.Entries
		for _, metaentry := range metaEntries {
			url, err := serializables.UrlToFilename(metaentry)
			if err != nil {
				return nil, err
			}
			schemaFile, err := serializables.ReadSchemaFile(url)
			if err != nil {
				return nil, err
			}
			if schemaFile == nil {
				continue
			}
			roleuri := "changeme"
			arcs := processLinkbaseRef(schemaFile, roleuri)
			cache.Set("rs:"+roleuri, RelationshipSetView{
				Arcs: arcs,
			}, gocache.DefaultExpiration)
		}
	}
	return &RelationshipSetRepo{
		lock:  sync.RWMutex{},
		cache: cache,
	}, nil
}

func processLinkbaseRef(file *serializables.SchemaFile, roleuri string) []string {
	if file == nil {
		return []string{}
	}
	ret := make([]string, 0)
	var wg sync.WaitGroup
	for _, annotation := range file.Annotation {
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
						//todo
						switch roleAttr.Value {
						case attr.PresentationLinkbaseRef:
							discoveredPre, err := serializables.ReadPresentationLinkbaseFile(hrefAttr.Value)
							if err != nil {
								return
							}
							fmt.Println("get arcs from " + discoveredPre.XMLName.Local)
						case attr.DefinitionLinkbaseRef:
							discoveredDef, err := serializables.ReadDefinitionLinkbaseFile(hrefAttr.Value)
							if err != nil {
								return
							}
							fmt.Println("get arcs from " + discoveredDef.XMLName.Local)
						case attr.CalculationLinkbaseRef:
							discoveredCal, err := serializables.ReadCalculationLinkbaseFile(hrefAttr.Value)
							if err != nil {
								return
							}
							fmt.Println("get arcs from " + discoveredCal.XMLName.Local)
						case attr.LabelLinkbaseRef:
							discoveredLab, err := serializables.ReadLabelLinkbaseFile(hrefAttr.Value)
							if err != nil {
								return
							}
							fmt.Println("get arcs from " + discoveredLab.XMLName.Local)
						default:
							break
						}
					}
				}(iitem)
			}
		}
	}
	wg.Wait()
	return ret
}
