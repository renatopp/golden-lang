// Given the module `@/foo/bar/hello.gold` at the absolute path `/d/project/foo/bar/hello.gold`:
//
// ImportName: 		 @/foo/bar/hello
// ModuleName: 		 hello
// ModuleFileName: hello.gold
// ModulePath: 		 /d/project/foo/bar/hello.gold
// PackageName: 	 @/foo/bar
// PackagePath:    /d/project/foo/bar

package fs

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Separator = "/"
var WorkingDir = ""

func init() {
	Separator = string(filepath.Separator)
	WorkingDir, _ = os.Getwd()
}

// Checks ---------------------------------------------------------------------

func CheckFileExists(path string) error {
	_, err := os.Stat(path)
	return err
}

func IsFileExtension(path, extension string, sensitive bool) bool {
	ext := path[len(path)-len(extension):]

	if sensitive {
		return ext == extension
	}

	return ext == extension
}

func CheckFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return os.ErrPermission
	}

	return nil
}

var validModuleName = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

func IsModuleNameValid(moduleName string) bool {
	return validModuleName.MatchString(moduleName)
}

// General Utilities ----------------------------------------------------------

// From a generic path, returns the absolute path
func GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// From a generic file path, returns the file extension (with dot)
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// From a generic path, returns the OS-compliant path
func ToOSSlash(path string) string {
	return strings.ReplaceAll(filepath.ToSlash(path), "/", Separator)
}

// From a generic path, returns the linux path
func ToLinuxSlash(path string) string {
	return filepath.ToSlash(path)
}

// Conversions ---------------------------------------------------------------
func ImportName_To_ModulePath(importName string) string {
	path := PackageName_To_PackagePath(importName)
	return path + ".gold"
}

func ModulePath_To_ModuleName(modulePath string) string {
	extension := GetFileExtension(modulePath)
	return filepath.Base(modulePath)[0 : len(filepath.Base(modulePath))-len(extension)]
}

func ModulePath_To_ModuleFileName(modulePath string) string {
	return filepath.Base(modulePath)
}

func ModulePath_To_PackageName(modulePath string) string {
	packagePath := ModulePath_To_PackagePath(modulePath)
	return PackagePath_To_PackageName(packagePath)
}

func ModulePath_To_PackagePath(modulePath string) string {
	return filepath.Dir(modulePath)
}

func PackagePath_To_PackageName(packagePath string) string {
	if strings.HasPrefix(packagePath, WorkingDir) {
		packagePath = filepath.Join("@", strings.TrimPrefix(packagePath, WorkingDir))
	}

	return ToLinuxSlash(packagePath)
}

func PackageName_To_PackagePath(packageName string) string {
	packageName = strings.ReplaceAll(packageName, "@", WorkingDir)

	path := ToOSSlash(packageName)
	return path
}

func Path_To_PackagePath(path string) string {
	if IsFileExtension(path, ".gold", false) {
		return ModulePath_To_PackagePath(path)
	}
	return path
}

func Path_To_PackageName(path string) string {
	if IsFileExtension(path, ".gold", false) {
		return ModulePath_To_PackageName(path)
	}
	return PackagePath_To_PackageName(path)
}

// ----------------------------------------------------------------------------

func ListFiles(path string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}

	return files, nil
}

// From the or package path, returns the package path
func DiscoverModules(path string) []string {
	var modules []string

	packagePath := Path_To_PackagePath(path)

	files, err := ListFiles(packagePath)
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		if !IsFileExtension(file, ".gold", false) {
			continue
		}

		modules = append(modules, file)
	}

	return modules
}