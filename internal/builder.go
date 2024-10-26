package internal

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/renatopp/golden/lang/strutils"
)

// A directory
type Package struct {
	Name    string
	Path    string // absolute path
	Private bool
	Modules []*Module
	Imports []*Package
}

func (p *Package) Debug() string {
	s := fmt.Sprint(strutils.PadRight("Package:", 15), "\n")
	s += fmt.Sprint(strutils.PadRight("  name:", 15), p.Name, "\n")
	s += fmt.Sprint(strutils.PadRight("  path:", 15), p.Path, "\n")
	s += fmt.Sprint(strutils.PadRight("  private:", 15), p.Private, "\n")
	s += fmt.Sprint(strutils.PadRight("  modules:", 15), "\n")
	for _, mod := range p.Modules {
		s += fmt.Sprint(strutils.PadRight("  - name:", 15), mod.Name, "\n")
		s += fmt.Sprint(strutils.PadRight("    path:", 15), mod.Path, "\n")
		s += fmt.Sprint(strutils.PadRight("    private:", 15), mod.Private, "\n")
		s += fmt.Sprint(strutils.PadRight("    package:", 15), mod.Package.Name, "\n")
	}
	return s
}

// A file
type Module struct {
	Name    string
	Package *Package
	Path    string // absolute path
	Private bool
	Imports []*Module
}

func (m *Module) Debug() string {
	s := fmt.Sprint(strutils.PadRight("Module:", 15), "\n")
	s += fmt.Sprint(strutils.PadRight("  name:", 15), m.Name, "\n")
	s += fmt.Sprint(strutils.PadRight("  path:", 15), m.Path, "\n")
	s += fmt.Sprint(strutils.PadRight("  private:", 15), m.Private, "\n")
	s += fmt.Sprint(strutils.PadRight("  package:", 15), m.Package.Name, "\n")
	return s
}

func ReadPackage(path string) (*Package, error) {
	pkg := createPackage(path)
	if !isValidName(pkg.Name) {
		return nil, fmt.Errorf("invalid package name: '%s'. Packages must be named strictly as `snake_case`", pkg.Name)
	}

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.ToLower(filepath.Ext(p)) == ".gold" {
			mod := createModule(pkg, p)
			if !isValidName(mod.Name) {
				return fmt.Errorf("invalid module name: '%s'. Modules must be named strictly as `snake_case`", mod.Name)
			}
			pkg.Modules = append(pkg.Modules, mod)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func createPackage(path string) *Package {
	name := filepath.Base(path)
	private := strings.HasPrefix(name, "_")
	modules := []*Module{}
	imports := []*Package{}
	return &Package{
		Name:    name,
		Path:    path,
		Private: private,
		Modules: modules,
		Imports: imports,
	}
}

func createModule(pkg *Package, path string) *Module {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	private := strings.HasPrefix(filepath.Base(path), "_")
	imports := []*Module{}
	return &Module{
		Package: pkg,
		Path:    path,
		Name:    name,
		Private: private,
		Imports: imports,
	}
}

func isValidName(name string) bool {
	match, _ := regexp.MatchString(`^[a-z][a-z_0-9]*$|^_[a-z_0-9]+$`, name)
	return match
}
