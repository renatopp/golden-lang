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
	backendProjectDirectory string
	backendMainPath         string
	backendGoModPath        string
}

func NewBackend() *Golang {
	return &Golang{}
}

func (b *Golang) Initialize(targetPath string) {
	targetDirectory = path.Join(targetPath, "golang")
	b.backendProjectDirectory = path.Join(targetDirectory, "root")
	b.backendMainPath = path.Join(targetDirectory, "main.go")
	b.backendGoModPath = path.Join(targetDirectory, "go.mod")
}

func (b *Golang) BeforeCodeGeneration() {
	fs.GuaranteeDirectoryExists(targetDirectory)
	fs.GuaranteeDirectoryExists(b.backendProjectDirectory)
}

func (b *Golang) GenerateCode(goldenFilePath string, root *ast.Module, entry bool) {
	backendFilePath := BackendPath(goldenFilePath)
	os.WriteFile(backendFilePath, []byte("package root\nfunc Main() { println(\"Hello, World from Go!\")}"), 0644)
}

func (b *Golang) AfterCodeGeneration() {
	os.WriteFile(b.backendMainPath, tmpl.GenerateBytes(template_main, nil), 0644)
	os.WriteFile(b.backendGoModPath, tmpl.GenerateBytes(template_mod, nil), 0644)
}

func (b *Golang) Run() {
	cmd := exec.Command("go", "run", b.backendMainPath)
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
