package core

type Module struct {
	Name     string   // Name of the module, ex: `hello`
	Path     string   // Absolute path of the module in the file system, ex: `/d/project/foo/bar/hello.gold`
	FileName string   // Name of the file, ex: `hello.gold`
	Package  *Package // Package that contains the module
	Root     Node     // Root node of the module, type is `*ast.Module`
}
