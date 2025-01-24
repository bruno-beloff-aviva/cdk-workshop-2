package hits

type Hits struct {
	Path  string `dynamodbav:"path"`
	Count int    `dynamodbav:"count"`
}

func NewHits(path string) Hits {
	return Hits{Path: path, Count: 0}
}

func (hit *Hits) GetKeys() map[string]interface{} {
	return map[string]interface{}{"path": hit.Path}
}

func (hit *Hits) Increment() {
	hit.Count++
}
