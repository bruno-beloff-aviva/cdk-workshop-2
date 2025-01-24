package hits

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	hit := NewHits("/test")
	fmt.Printf("Hit: %#v\n", hit)

	assert.Equal(t, hit.Path, "/test")
	assert.Equal(t, hit.Count, 0)
}

func TestInc(t *testing.T) {
	hit := NewHits("/test")
	hit.Increment()
	fmt.Printf("Hit: %#v\n", hit)

	assert.Equal(t, hit.Path, "/test")
	assert.Equal(t, hit.Count, 1)
}

func TestGetKeys(t *testing.T) {
	hit := NewHits("/test")
	keys := hit.GetKeys()
	expected := map[string]interface{}(map[string]interface{}{"path": "/test"})
	fmt.Printf("keys: %#v\n", keys)

	assert.Equal(t, keys, expected)
}
