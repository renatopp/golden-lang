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

func Test_ModulePath_To_ModuleName(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	moduleName := fs.ModulePath_To_ModuleName(modulePath)
	assert.Equal(t, "hello", moduleName)
}

func Test_ModulePath_To_ModuleFileName(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	moduleFileName := fs.ModulePath_To_ModuleFileName(modulePath)
	assert.Equal(t, "hello.gold", moduleFileName)
}

func Test_ModulePath_To_PackageName(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	packageName := fs.ModulePath_To_PackageName(modulePath)
	assert.Equal(t, "@/foo/bar", packageName)
}

func Test_ModulePath_To_PackagePath(t *testing.T) {
	resets()

	modulePath := variant("/d/project/foo/bar/hello.gold")
	packagePath := fs.ModulePath_To_PackagePath(modulePath)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)
}

func Test_PackagePath_To_PackageName(t *testing.T) {
	resets()

	packagePath := variant("/d/project/foo/bar")
	packageName := fs.PackagePath_To_PackageName(packagePath)
	assert.Equal(t, "@/foo/bar", packageName)
}

func Test_PackageName_To_PackagePath(t *testing.T) {
	resets()

	packageName := "@/foo/bar"
	packagePath := fs.PackageName_To_PackagePath(packageName)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)
}

func Test_Path_To_PackagePath(t *testing.T) {
	resets()

	path := variant("/d/project/foo/bar/hello.gold")
	packagePath := fs.Path_To_PackagePath(path)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)

	path = variant("/d/project/foo/bar")
	packagePath = fs.Path_To_PackagePath(path)
	assert.Equal(t, variant("/d/project/foo/bar"), packagePath)
}

func Test_Path_To_PackageName(t *testing.T) {
	resets()

	path := variant("/d/project/foo/bar/hello.gold")
	packagePath := fs.Path_To_PackageName(path)
	assert.Equal(t, "@/foo/bar", packagePath)

	path = variant("/d/project/foo/bar")
	packagePath = fs.Path_To_PackageName(path)
	assert.Equal(t, "@/foo/bar", packagePath)
}
