package internal

type Scope struct {
	Parent   *Scope
	Bindings map[string]*Node
	Types    map[string]*Node
}

func Analyze(scope *Scope, modules []*Node) {

}
