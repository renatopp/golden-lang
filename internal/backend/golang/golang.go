package golang

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/internal/helpers/tmpl"
)

//go:embed templates/main.go.tmpl
var raw_template_main string
var template_main, _ = template.New("main").Parse(raw_template_main)

//go:embed templates/go.mod.tmpl
var raw_template_mod string
var template_mod, _ = template.New("mod").Parse(raw_template_mod)

type Golang struct {
}

func NewBackend() *Golang {
	return &Golang{}
}

func (b *Golang) Initialize(targetPath string) {
	targetDirectory = path.Join(targetPath, "golang")
}

func (b *Golang) BeforeCodeGeneration() {
	fs.GuaranteeDirectoryExists(targetDirectory)
	fs.GuaranteeDirectoryExists(path.Join(targetDirectory, "root"))
}

func (b *Golang) GenerateCode(goldenFilePath string, root *ast.Module, entry bool) {
	backendFilePath := BackendPath(goldenFilePath)
	os.WriteFile(backendFilePath, []byte("package root\nfunc Main() { println(\"Hello, World from Go!\")}"), 0644)
}

func (b *Golang) AfterCodeGeneration() {
	fileName := path.Join(targetDirectory, "main.go")
	os.WriteFile(fileName, tmpl.GenerateBytes(template_main, nil), 0644)

	fileName = path.Join(targetDirectory, "go.mod")
	os.WriteFile(fileName, tmpl.GenerateBytes(template_mod, nil), 0644)
}

func (b *Golang) Run() {
	cmd := exec.Command("go", "run", path.Join(targetDirectory, "main.go"))
	cmd.Dir = targetDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run node: %v", err)
	}
}

func (b *Golang) Build() {}

func (b *Golang) Finalize() {
}
