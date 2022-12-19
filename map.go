package concurrency_map

import (
	"hash/fnv"
	"sync"
)

const SharedCount = 1024

type ConcurrentMapShared struct {
	items map[string]any
	sync.RWMutex
}

type ConcurrentMap []*ConcurrentMapShared

func New() ConcurrentMap {
	m := make(ConcurrentMap, SharedCount)
	for i := 0; i < SharedCount; i++ {
		m[i] = &ConcurrentMapShared{
			items: map[string]any{},
		}
	}

	return m
}

func (m ConcurrentMap) getShared(key string) *ConcurrentMapShared {
	h := fnv.New64a()
	h.Write([]byte(key))

	return m[h.Sum64()&(SharedCount-1)]
}

func (m ConcurrentMap) Set(key string, v any) {
	shared := m.getShared(key)
	shared.Lock()
	shared.items[key] = v
	shared.Unlock()
}

func (m ConcurrentMap) Get(key string) (any, bool) {
	shared := m.getShared(key)
	shared.RLock()
	v, ok := shared.items[key]
	shared.RUnlock()
	return v, ok
}

func (m ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < SharedCount; i++ {
		shared := m[i]
		shared.RLock()
		count += len(shared.items)
		shared.RUnlock()
	}
	return count
}

func (m ConcurrentMap) Delete(key string) {
	shared := m.getShared(key)
	shared.Lock()
	delete(shared.items, key)
	shared.Unlock()
}

func (m ConcurrentMap) Keys() []string {
	count := m.Count()
	ch := make(chan string, count)
	go func() {
		var wg sync.WaitGroup
		wg.Add(SharedCount)
		for _, shared := range m {
			go func(shared *ConcurrentMapShared) {
				shared.RLock()
				for key := range shared.items {
					ch <- key
				}
				shared.RUnlock()
				wg.Done()
			}(shared)
		}
		close(ch)
	}()

	keys := make([]string, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}

	return keys
}
