package router

import "sync"

// CORSCache holds the live set of allowed origins.
// Static origins (from env) are loaded once; DB origins can be added/removed at runtime.
type CORSCache struct {
	mu      sync.RWMutex
	static  map[string]struct{} // from env var, read-only after init
	dynamic map[string]struct{} // from DB, mutable
}

func NewCORSCache(staticOrigins []string) *CORSCache {
	c := &CORSCache{
		static:  make(map[string]struct{}, len(staticOrigins)),
		dynamic: make(map[string]struct{}),
	}
	for _, o := range staticOrigins {
		if o != "" {
			c.static[o] = struct{}{}
		}
	}
	return c
}

func (c *CORSCache) LoadDynamic(origins []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dynamic = make(map[string]struct{}, len(origins))
	for _, o := range origins {
		c.dynamic[o] = struct{}{}
	}
}

func (c *CORSCache) Add(origin string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dynamic[origin] = struct{}{}
}

func (c *CORSCache) Remove(origin string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.dynamic, origin)
}

func (c *CORSCache) Allow(origin string) bool {
	if _, ok := c.static[origin]; ok {
		return true
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.dynamic[origin]
	return ok
}
