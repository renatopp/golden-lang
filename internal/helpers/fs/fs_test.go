package fs_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/stretchr/testify/assert"
)

// Forces OS-variant data to a fixed test scenario
func resets() {
	fs.WorkingDir = variant("/d/project")
}

func invariant(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func variant(path string) string {
	return strings.ReplaceAll(invariant(path), "/", string(filepath.Separator))
}

func Test_ModulePath2ModuleName(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	moduleName := fs.ModulePath2ModuleName(modulePath)
	assert.Equal(t, "hello", moduleName)
}

func Test_ModulePath2ModuleFileName(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	moduleFileName := fs.ModulePath2ModuleFileName(modulePath)
	assert.Equal(t, "hello.gold", moduleFileName)
}

func Test_ModulePath2PackageName(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	packageName := fs.ModulePath2PackageName(modulePath)
	assert.Equal(t, "@/foo/bar", packageName)
}

func Test_ModulePath2PackagePath(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	packagePath := fs.ModulePath2PackagePath(modulePath)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)
}

func Test_PackagePath2PackageName(t *testing.T) {
	resets()

	packagePath := variant("/d/project/foo/bar")
	packageName := fs.PackagePath2PackageName(packagePath)
	assert.Equal(t, "@/foo/bar", packageName)
}

func Test_PackageName2PackagePath(t *testing.T) {
	resets()

	packageName := "@/foo/bar"
	packagePath := fs.PackageName2PackagePath(packageName)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)
}

func Test_Path2PackagePath(t *testing.T) {
	resets()

	path := variant("/d/project/foo/bar/hello.gold")
	packagePath := fs.Path2PackagePath(path)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)

	path = variant("/d/project/foo/bar")
	packagePath = fs.Path2PackagePath(path)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)
}

func Test_Path2PackageName(t *testing.T) {
	resets()

	path := variant("/d/project/foo/bar/hello.gold")
	packagePath := fs.Path2PackageName(path)
	assert.Equal(t, "@/foo/bar", packagePath)

	path = variant("/d/project/foo/bar")
	packagePath = fs.Path2PackageName(path)
	assert.Equal(t, "@/foo/bar", packagePath)
}
