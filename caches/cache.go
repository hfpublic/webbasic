package caches

type cacheFunc func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type request struct {
	key      string
	response chan<- result
}

type entry struct {
	res   result
	ready chan struct{}
}

// Cache 缓存
type Cache struct {
	requests chan request
	cf       cacheFunc
	reget    bool
}

func New(reget bool, f cacheFunc) *Cache {
	c := &Cache{
		requests: make(chan request),
		cf:       f,
		reget:    reget,
	}
	go c.server()
	return c
}

func (c *Cache) Get(key string) (interface{}, error) {
	response := make(chan result)
	c.requests <- request{key: key, response: response}
	res := <-response
	return res.value, res.err
}

func (c *Cache) Close() {
	close(c.requests)
}

func (c *Cache) server() {
	cache := make(map[string]*entry)
	for req := range c.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(c.cf, req.key, c.reget)
		}
		go e.deliver(req.response, c.cf, req.key, c.reget)
	}
}

func (e *entry) call(f cacheFunc, key string, reget bool) {
	e.res.value, e.res.err = f(key)
	if reget && e.res.err != nil {
		e.ready <- struct{}{}
	} else {
		close(e.ready)
	}
}

func (e *entry) deliver(response chan<- result, cf cacheFunc, key string, reget bool) {
	_, ok := <-e.ready
	response <- e.res
	if ok {
		go e.call(cf, key, reget)
	}
}
