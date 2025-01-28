package hits

type Hits struct {
	Path  string `dynamodbav:"path"`
	Count uint   `dynamodbav:"count"`
}

func NewHits(path string) Hits {
	return Hits{Path: path, Count: 0}
}

func (hit *Hits) GetKey() map[string]any {
	return map[string]any{"path": hit.Path}
}

func (hit *Hits) Increment() {
	hit.Count++
}
