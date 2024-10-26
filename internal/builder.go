package internal

// A directory
type Package struct {
	Name    string
	Path    string // absolute path
	Private bool
	Modules []*Module
	Imports []*Package
}

// A file
type Module struct {
	Package *Package
	Path    string // absolute path
	Private bool
	Imports []*Module
}
