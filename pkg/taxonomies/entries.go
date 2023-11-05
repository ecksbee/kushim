package taxonomies

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ecksbee.com/kushim/internal/actions"
	"ecksbee.com/telefacts/pkg/attr"
	"ecksbee.com/telefacts/pkg/serializables"
	"github.com/joshuanario/r8lmt"
)

var (
	wLock     sync.Mutex
	out       chan interface{} = make(chan interface{})
	in        chan interface{} = make(chan interface{})
	dur       time.Duration    = 200 * time.Millisecond
	throttled bool             = false
)

func startSECThrottle() {
	if !throttled {
		r8lmt.Throttler(out, in, dur, false)
		throttled = true
	}
}

func Discover(entries []string) error {
	if VolumePath == "" {
		return fmt.Errorf("empty VolumePath")
	}
	startSECThrottle()
	serializables.GlobalTaxonomySetPath = VolumePath
	for _, entry := range entries {
		url, err := serializables.UrlToFilename(entry)
		if err != nil {
			return err
		}
		schemaFile, err := serializables.ReadSchemaFile(url)
		if err != nil {
			return err
		}
		if schemaFile == nil {
			continue
		}
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			ImportSchema(schemaFile)
		}()
		go func() {
			defer wg.Done()
			IncludeSchema(schemaFile)
		}()
		go func() {
			defer wg.Done()
			LinkbaseRefSchema(schemaFile)
		}()
		wg.Wait()
	}
	return nil
}

func LinkbaseRefSchema(file *serializables.SchemaFile) {
	if file == nil {
		return
	}
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
						go DiscoverRemoteURL(hrefAttr.Value)
						return
					}
				}(iitem)
			}
		}
	}
	wg.Wait()
}

func IncludeSchema(file *serializables.SchemaFile) {
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
				go DiscoverRemoteURL(schemaLocationAttr.Value)
				return
			}
		}(iitem)
	}
	wg.Wait()
}

func ImportSchema(file *serializables.SchemaFile) {
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
			namespaceAttr := attr.FindAttr(item.XMLAttrs, "namespace")
			if namespaceAttr == nil || namespaceAttr.Value == "" {
				return
			}
			schemaLocationAttr := attr.FindAttr(item.XMLAttrs, "schemaLocation")
			if schemaLocationAttr == nil || schemaLocationAttr.Value == "" {
				return
			}
			if attr.IsValidUrl(schemaLocationAttr.Value) {
				DiscoverRemoteURL(schemaLocationAttr.Value)
				return
			}
		}(iitem)
	}
	wg.Wait()
}

func throttle(urlString string) {
	urlStruct, err := url.Parse(urlString)
	if urlStruct.Hostname() != "sec.gov" {
		return
	}
	if err != nil {
		return
	}
	in <- struct{}{}
	<-out
}

func DiscoverRemoteURL(url string) {
	dest, err := serializables.UrlToFilename(url)
	if err != nil {
		return
	}
	_, err = os.Stat(dest)
	if os.IsNotExist(err) {
		targetDir := filepath.Dir(dest)
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return
		}
		body, err := actions.Scrape(url, throttle)
		if err != nil {
			return
		}
		wLock.Lock()
		err = actions.WriteFile(dest, body)
		if err != nil {
			return
		}
		wLock.Unlock()
	}
	discoveredSchema, err := serializables.ReadSchemaFile(dest)
	if err != nil {
		return
	}
	var wwg sync.WaitGroup
	wwg.Add(2)
	go func() {
		defer wwg.Done()
		ImportSchema(discoveredSchema)
	}()
	go func() {
		defer wwg.Done()
		IncludeSchema(discoveredSchema)
	}()
	wwg.Wait()
}
