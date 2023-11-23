package cache

import (
	"io/ioutil"
	"path"
	"sync"

	"ecksbee.com/kushim/pkg/taxonomies"
	gocache "github.com/patrickmn/go-cache"
)

var (
	once     sync.Once
	appCache *gocache.Cache
	lock     sync.RWMutex
)

func InitRepo(gts string) {
	once.Do(func() {
		appCache = gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	})
}

func MarshalCatalog() ([]byte, error) {
	lock.RLock()
	if x, found := appCache.Get("_kushim"); found {
		ret := x.([]byte)
		lock.RUnlock()
		return ret, nil
	}
	lock.RUnlock()
	target := path.Join(taxonomies.VolumePath, "_.json")
	data, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, err
	}
	go func() {
		lock.Lock()
		defer lock.Unlock()
		appCache.Set("_kushim", data, gocache.DefaultExpiration)
	}()
	return data, err
}

func MarshalRenderable(hash string) ([]byte, error) {
	lock.RLock()
	if x, found := appCache.Get(hash); found {
		ret := x.([]byte)
		lock.RUnlock()
		return ret, nil
	}
	lock.RUnlock()
	target := path.Join(taxonomies.VolumePath, hash+".json")
	data, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, err
	}
	go func() {
		lock.Lock()
		defer lock.Unlock()
		appCache.Set(hash, data, gocache.DefaultExpiration)
	}()
	return data, err
}
