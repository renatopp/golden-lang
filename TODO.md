# TODOr

## Refactor
[ ] try to create a traverser (tick => check a node recursively)
  [ ] hook every node on: open (before ticking children), close (after ticking children)
  [ ] allow stopping ticking
  [ ] allow manual tick continuation (eg: for just a child)
[ ] organize semantic step
  [x] move types to own package
  [ ] try to split type checking from other checks
  [ ] shouldn't scope be improved?
  
  [ ] allow raw expression like: 1 + 1.0
  [ ] visibility check (access `_`)
  [ ] improve panic usage in the analyser - remove it?
[ ] organize build step
  [ ] remove Package references from places other than module
  [ ] split worker into worker and /manager/
  [ ] improve panic usage (at least catch before returning to command)
  [ ] add pretty error messages
  [ ] visibility check for modules `_`

## IR
[ ] create ir representation
  all calls will reside in the same "scope" for each package
  should reduce the ast into small blocks like (`bind <name> <type> <value>`, `register function`, `call function`) 
[ ] convert ast to ir
  [ ] flatify function definitions
  [ ] flatify expressions
  [ ] flatify function calls
  [ ] flatify scopes
  [ ] unique naming
[ ] optimize ir
[ ] plan codegen
[ ] create blog post:
  [ ] golden language
  [ ] how to develop a language
