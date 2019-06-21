package loader

// Loader loads graph
type Loader interface {
	Load(id string) (uint32, []string, error)
}
