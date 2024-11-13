# TODOr

## Refactor
[x] errors
  [x] refactor errors on lexer
  [x] refactor errors on parser
  [x] refactor errors on builder
  [x] add pretty error messages
  
[x] organize semantic step
  [x] move types to own package
  [x] improve panic usage in the analyser - remove it?
  [x] finish resolver refactor
  [x] pass resolve function to the module type

[x] organize build step
  [x] improve panic usage (at least catch before returning to command)
  [x] refactor steps for better readability
  [x] visibility check for modules `_`

[ ] improve semantic step
  [x] add binding to scope instead of value/type alone
  [x] visibility check (access `_`)
  [ ] allow raw expression like: 1 + 1.0

[ ] others
  [ ] parse comments inside the code

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
