package hits

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	hit := NewHits("/test")
	fmt.Println(hit.String())

	assert.Equal(t, hit.Path, "/test")
	assert.Equal(t, hit.Count, uint(1))
}

func TestInc(t *testing.T) {
	hit := NewHits("/test")
	hit.Increment()
	fmt.Println(hit.String())

	assert.Equal(t, hit.Path, "/test")
	assert.Equal(t, hit.Count, uint(2))
}

func TestGetKey(t *testing.T) {
	hit := NewHits("/test")
	keys := hit.GetKey()
	expected := map[string]any(map[string]any{"path": "/test"})
	fmt.Printf("key: %#v\n", keys)

	assert.Equal(t, keys, expected)
}
