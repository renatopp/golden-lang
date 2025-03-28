package interpreter

import "github.com/renatopp/golden/internal/compiler/ast"

type Interpreter struct {
}

func NewBackend() *Interpreter {
	return &Interpreter{}
}

func (b *Interpreter) Initialize(targetPath string) {}

func (b *Interpreter) BeforeCodeGeneration() {}

func (b *Interpreter) GenerateCode(filePath string, root *ast.Module, entry bool) {}

func (b *Interpreter) AfterCodeGeneration() {}

func (b *Interpreter) Run() {}

func (b *Interpreter) Build(outputPath string) {}

func (b *Interpreter) Finalize() {}
