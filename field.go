package clogger

import (
	"sync"
)

type field struct {
	key   string
	value interface{}
}

type fieldCollection struct {
	m  map[string]interface{}
	mu *sync.Mutex
}

func newFieldCollection() *fieldCollection {
	return &fieldCollection{
		m:  make(map[string]interface{}),
		mu: &sync.Mutex{},
	}
}

func (fc *fieldCollection) add(key string, value interface{}) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.m[key] = value
}

func (fc *fieldCollection) len() int {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	return len(fc.m)
}

func (fc *fieldCollection) retrieve(key string) interface{} {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	return fc.m[key]
}

func (fc *fieldCollection) addField(f field) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.m[f.key] = f.value
}

func (fc *fieldCollection) merge(f *fieldCollection) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	for k, v := range f.m {
		fc.m[k] = v
	}
}

func (fc *fieldCollection) fields() []field {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	out := make([]field, 0, len(fc.m))
	for k, v := range fc.m {
		out = append(out, field{k, v})
	}

	return out
}
