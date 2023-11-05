package cache

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"ecksbee.com/kushim/pkg/taxonomies"
	"ecksbee.com/telefacts/pkg/attr"
	"ecksbee.com/telefacts/pkg/serializables"
	gocache "github.com/patrickmn/go-cache"
)

func MarshalNamespace(namespace string) ([]byte, error) {
	data := nsrepo.Query(namespace)
	return json.Marshal(data)
}

type NamespaceRepo struct {
	lock  sync.RWMutex
	cache *gocache.Cache
}

func (repo *NamespaceRepo) Query(namespace string) *NamespaceView {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	if x, found := repo.cache.Get("ns:" + namespace); found {
		ret := x.(NamespaceView)
		return &ret
	}
	return nil
}

func NewNamespaceRepo(cache *gocache.Cache, gts string) (*NamespaceRepo, error) {
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
			var wg sync.WaitGroup
			wg.Add(3)
			go func() {
				defer wg.Done()
				processLocalLinkRoles(cache, schemaFile, metaentry)
			}()
			go func() {
				defer wg.Done()
				includeSchemaLinkRoles(cache, schemaFile, metaentry)
			}()
			go func() {
				defer wg.Done()
				importSchemaLinkRoles(cache, schemaFile, metaentry)
			}()
			wg.Wait()
		}
	}
	return &NamespaceRepo{
		lock:  sync.RWMutex{},
		cache: cache,
	}, nil
}

func processLocalLinkRoles(cache *gocache.Cache, file *serializables.SchemaFile, url string) {
	if file == nil {
		return
	}
	targetAttr := attr.FindAttr(file.XMLAttrs, "targetNamespace")
	if targetAttr == nil || targetAttr.Value == "" {
		return
	}
	ret := make([]string, 0)
	for _, annotation := range file.Annotation {
		for _, appinfo := range annotation.Appinfo {
			for _, roleType := range appinfo.RoleType {
				uriAttr := attr.FindAttr(roleType.XMLAttrs, "roleURI")
				if uriAttr == nil || uriAttr.Value == "" {
					continue
				}
				ret = append(ret, uriAttr.Value)
			}
		}
	}
	cache.Set("ns:"+targetAttr.Value, NamespaceView{
		RelationshipSets: ret,
	}, gocache.DefaultExpiration)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		includeSchemaLinkRoles(cache, file, url)
	}()
	go func() {
		defer wg.Done()
		importSchemaLinkRoles(cache, file, url)
	}()
	wg.Wait()
}

func includeSchemaLinkRoles(cache *gocache.Cache, file *serializables.SchemaFile, url string) {
	if file == nil {
		return
	}
	includes := file.Include
	var wg sync.WaitGroup
	wg.Add(len(includes))
	for _, iitem := range includes {
		go func(item struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr "xml:\",any,attr\""
		}) {
			defer wg.Done()
			if item.XMLName.Space != attr.XSD {
				return
			}
			schemaLocationAttr := attr.FindAttr(item.XMLAttrs, "schemaLocation")
			if schemaLocationAttr == nil || schemaLocationAttr.Value == "" {
				return
			}
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				myurl, err := serializables.UrlToFilename(schemaLocationAttr.Value)
				if err != nil {
					return
				}
				schemaFile, err := serializables.ReadSchemaFile(myurl)
				if err != nil {
					return
				}
				if schemaFile == nil {
					return
				}
				processLocalLinkRoles(cache, schemaFile, schemaLocationAttr.Value)
				return
			}
			urlDir := path.Dir(url)
			filepath := path.Join(urlDir, schemaLocationAttr.Value)
			discoveredSchema, err := serializables.ReadSchemaFile(filepath)
			if err != nil {
				return
			}
			if discoveredSchema == nil {
				return
			}
			processLocalLinkRoles(cache, discoveredSchema, schemaLocationAttr.Value)
		}(iitem)
	}
	wg.Wait()
}

func importSchemaLinkRoles(cache *gocache.Cache, file *serializables.SchemaFile, url string) {
	if file == nil {
		return
	}
	imports := file.Import
	var wg sync.WaitGroup
	wg.Add(len(imports))
	for _, iitem := range imports {
		go func(item struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr "xml:\",any,attr\""
		}) {
			defer wg.Done()
			if item.XMLName.Space != attr.XSD {
				return
			}
			schemaLocationAttr := attr.FindAttr(item.XMLAttrs, "schemaLocation")
			if schemaLocationAttr == nil || schemaLocationAttr.Value == "" {
				return
			}
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				myurl, err := serializables.UrlToFilename(schemaLocationAttr.Value)
				if err != nil {
					return
				}
				schemaFile, err := serializables.ReadSchemaFile(myurl)
				if err != nil {
					return
				}
				if schemaFile == nil {
					return
				}
				processLocalLinkRoles(cache, schemaFile, schemaLocationAttr.Value)
				return
			}
			urlDir := path.Dir(url)
			filepath := path.Join(urlDir, schemaLocationAttr.Value)
			discoveredSchema, err := serializables.ReadSchemaFile(filepath)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
			if discoveredSchema == nil {
				return
			}
			processLocalLinkRoles(cache, discoveredSchema, schemaLocationAttr.Value)
		}(iitem)
	}
	wg.Wait()
	return
}
