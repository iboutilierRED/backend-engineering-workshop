package models

type WikipediaChange struct {
	User string `json:"user"`
	Uri  string `json:"uri"`
	Bot  bool   `json:"bot"`
	Meta Meta   `json:"meta"`
}

type Meta struct {
	Uri string `json:"uri"`
}
