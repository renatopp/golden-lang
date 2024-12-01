package javascript

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

//go:embed templates/main.mjs.tmpl
var raw_template_main string
var template_main, _ = template.New("main").Parse(raw_template_main)

type Javascript struct {
	entryRef        *Ref
	backendMainPath string
}

func NewBackend() *Javascript {
	return &Javascript{}
}

func (b *Javascript) Initialize(targetPath string) {
	targetDirectory = path.Join(targetPath, "javascript")
	b.backendMainPath = path.Join(targetDirectory, "main.mjs")
}

func (b *Javascript) BeforeCodeGeneration() {
	fs.GuaranteeDirectoryExists(targetDirectory)
}

func (b *Javascript) GenerateCode(goldenFilePath string, root *ast.Module, entry bool) {
	if entry {
		b.entryRef = R(goldenFilePath, "main")
	}
	backendFilePath := BackendPath(goldenFilePath)
	os.WriteFile(backendFilePath, []byte("export function main() { console.log('Hello, World from JS!') }"), 0644)
}

func (b *Javascript) AfterCodeGeneration() {
	os.WriteFile(b.backendMainPath, tmpl.GenerateBytes(template_main, map[string]any{
		"EntryImport": b.entryRef.BackendImportPath,
	}), 0644)
}

func (b *Javascript) Run() {
	cmd := exec.Command("node", b.backendMainPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run node: %v", err)
	}
}

func (b *Javascript) Build() {}

func (b *Javascript) Finalize() {}
