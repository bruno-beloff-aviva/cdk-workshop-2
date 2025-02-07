package hits

import "fmt"

type Hits struct {
	Path  string `dynamodbav:"path"`
	Count uint   `dynamodbav:"count"`
}

func NewHits(path string) Hits {
	return Hits{Path: path, Count: 1}
}

func (h *Hits) String() string {
	return fmt.Sprintf("Hit:{Path:%s Count:%v}", h.Path, h.Count)
}

func (h *Hits) GetKey() map[string]any {
	return map[string]any{"path": h.Path}
}

func (h *Hits) Increment() {
	h.Count++
}
