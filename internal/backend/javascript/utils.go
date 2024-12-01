package javascript

import (
	"path"
	"strings"

	"github.com/renatopp/golden/internal/helpers/fs"
)

type Ref struct {
	GoldenFilePath    string
	GoldenIdentifier  string
	BackendFilePath   string
	BackendImportPath string
	BackendIdentifier string
}

var targetDirectory = ""

func R(filepath, identifier string) *Ref {
	return &Ref{
		GoldenFilePath:    filepath,
		GoldenIdentifier:  identifier,
		BackendFilePath:   BackendPath(filepath),
		BackendImportPath: BackendImportPath(filepath),
		BackendIdentifier: BackendIdentifier(identifier),
	}
}

// Returns the absolute path of the backend file
func BackendPath(goldenFilepath string) string {
	file := _relativeBackendPath(goldenFilepath)
	return path.Join(targetDirectory, file)
}

func BackendImportPath(goldenFilepath string) string {
	file := _relativeBackendPath(goldenFilepath)
	return "./" + file
}

func BackendIdentifier(goldenIdentifier string) string {
	return goldenIdentifier
}

func _relativeBackendPath(goldenFilepath string) string {
	// if filepath is from project, root/**
	// if filepath is from core, core/**
	// if filepath is from packages, package<name>/**

	file := ""
	if fs.IsProjectPath(goldenFilepath) {
		relative := fs.ToLinuxSlash(fs.GetProjectRelativePath(goldenFilepath))
		file = "root" + strings.ReplaceAll(relative, "/", "_")
	} else {
		panic("BackendPath not implemented: " + goldenFilepath)
	}

	return strings.ReplaceAll(file, ".gold", ".mjs")
}
