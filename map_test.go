package concurrency_map_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	concurrencyMap "github.com/seymourtang/concurrency-map"
)

func TestConcurrentMap(t *testing.T) {
	var (
		key1   = "key1"
		value1 = "value1"
		key2   = "key2"
		value2 = "value2"
	)

	m := concurrencyMap.New()
	m.Set(key1, value1)
	assert.Equal(t, 1, m.Count())

	v, ok := m.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, value1, v)

	m.Set(key2, value2)
	assert.Equal(t, 2, m.Count())

	m.Delete(key1)

	_, ok = m.Get(key1)
	assert.False(t, ok)
	assert.Equal(t, 1, m.Count())
	assert.Equal(t, []string{key2}, m.Keys())
}
