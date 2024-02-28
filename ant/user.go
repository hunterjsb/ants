package ant

type User struct {
	Name     string    `json:"name"`
	Colonies []*Colony `json:"colonies"`
}
