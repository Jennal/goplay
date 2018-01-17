package session

import (
	"sync"

	"github.com/jennal/goplay/log"
)

type pushCacheItem struct {
	Route string
	Data  interface{}
}

type pushCache struct {
	sync.Mutex

	cache     []*pushCacheItem
	cacheOnce map[string]interface{}
}

func (pc *pushCache) PushCache(route string, data interface{}) {
	pc.Lock()
	defer pc.Unlock()

	pc.cache = append(pc.cache, &pushCacheItem{
		Route: route,
		Data:  data,
	})
	log.Logf("PushCache: %v %v", len(pc.cache), pc.cache)
}

func (pc *pushCache) PushCacheOnce(route string, data interface{}) {
	pc.Lock()
	defer pc.Unlock()

	if pc.cacheOnce == nil {
		pc.cacheOnce = make(map[string]interface{})
	}

	pc.cacheOnce[route] = data
}

func (pc *pushCache) PopAllCaches() []*pushCacheItem {
	pc.Lock()
	defer pc.Unlock()

	result := pc.cache
	for route, data := range pc.cacheOnce {
		result = append(result, &pushCacheItem{
			Route: route,
			Data:  data,
		})
	}

	pc.cache = nil
	pc.cacheOnce = nil

	return result
}
