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
	entryRef *Ref
}

func NewBackend() *Javascript {
	return &Javascript{}
}

func (b *Javascript) Initialize(targetPath string) {
	targetDirectory = path.Join(targetPath, "javascript")
}

func (b *Javascript) BeforeCodeGeneration() {
	fs.GuaranteeDirectoryExists(targetDirectory)
}

func (b *Javascript) GenerateCode(goldenFilePath string, root *ast.Module, entry bool) {
	backendFilePath := BackendPath(goldenFilePath)
	os.WriteFile(backendFilePath, []byte("export function main() { console.log('hello, world!') }"), 0644)
}

func (b *Javascript) AfterCodeGeneration() {
	fileName := path.Join(targetDirectory, "main.mjs")
	os.WriteFile(fileName, tmpl.GenerateBytes(template_main, map[string]any{
		"EntryImport": b.entryRef.BackendImportPath,
	}), 0644)
}

func (b *Javascript) Run() {
	cmd := exec.Command("node", path.Join(targetDirectory, "main.mjs"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run node: %v", err)
	}
}

func (b *Javascript) Build() {}

func (b *Javascript) Finalize() {}
