# TODOr

## Refactor
[ ] errors
  [ ] refactor errors on lexer
  [ ] refactor errors on parser
  [ ] refactor errors on builder
  [ ] add pretty error messages
  
[ ] organize semantic step
  [x] move types to own package
  [x] improve panic usage in the analyser - remove it?
  [x] finish resolver refactor
  [ ] pass resolve function to the module type
  [ ] shouldn't scope be improved?
  
  [ ] allow raw expression like: 1 + 1.0
  [ ] visibility check (access `_`)

[ ] organize build step
  [ ] remove Package references from places other than module
  [ ] split worker into worker and /manager/
  [ ] improve panic usage (at least catch before returning to command)
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
