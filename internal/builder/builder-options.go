package builder

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/events"
	"github.com/renatopp/golden/internal/helpers/fs"
)

type BuildOptions struct {
	EntryFilePath    string // Absolute path of the entry file containing main function
	OutputFilePath   string // Absolute path of the output file for storing the final executable
	WorkingDir       string // Absolute path of the project working directory
	LocalCachePath   string // Absolute path of the local cache directory for storing project's ASTs
	GlobalCachePath  string // Absolute path of the global cache directory for storing bundles and core ASTs
	LocalTargetPath  string // Absolute path of the local target directory for storing transpiled files
	GlobalTargetPath string // Absolute path of the global target directory for storing transpiled files

	// Events
	OnTokensReady          *events.Signal2[*File, []token.Token]
	OnAstReady             *events.Signal2[*File, ast.Module]
	OnDependencyGraphReady *events.Signal1[[]*File]
	// OnTypeCheckReady       *events.Signal3[*File, *ast.Module, *env.Scope]
}

func NewBuildOptions(fileName string) *BuildOptions {
	return &BuildOptions{
		EntryFilePath:    fileName,
		OutputFilePath:   fs.GetBinaryName(fileName),
		WorkingDir:       fs.GetWorkingDir(),
		LocalCachePath:   fs.JoinProjectPath(".cache"),
		GlobalCachePath:  fs.JoinLangPath("cache"),
		LocalTargetPath:  fs.JoinProjectPath(".target"),
		GlobalTargetPath: fs.JoinLangPath("target"),

		OnTokensReady:          events.NewSignal2[*File, []token.Token](),
		OnAstReady:             events.NewSignal2[*File, ast.Module](),
		OnDependencyGraphReady: events.NewSignal1[[]*File](),
		// OnTypeCheckReady:       events.NewSignal3[*File, *ast.Module, *env.Scope](),
	}
}
