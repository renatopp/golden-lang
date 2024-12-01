package backend

import "github.com/renatopp/golden/internal/compiler/ast"

type Backend interface {
	Initialize(targetPath string)
	BeforeCodeGeneration()
	GenerateCode(filePath string, root *ast.Module, entry bool)
	AfterCodeGeneration()
	Run()
	Build()
	Finalize()
}
