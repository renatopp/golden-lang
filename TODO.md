# TODO

[ ] re-structure
  [ ] organize lexer & parser into its own folder
  [ ] organize semantic analysis in its own folder

[ ] organize steps: type checking, inference, binding checks, cyclic reference check, argument validations, control flow checks, visibility checks, mutability, etc.
  [ ] Treat panics during build
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







## Pending

[ ] Print pretty error messages



## Notes

```
./.tools/tcc/win/tcc.exe -run ./.tools/tcc/win/examples/fib.c 30
```