package core

type GraphNode struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Category string `json:"category"`
	IsDirty  bool   `json:"is_dirty"`
	Label    string `json:"label"`
}

type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type RenderTree struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}
