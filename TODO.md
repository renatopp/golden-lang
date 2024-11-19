# TODO

[ ] Build Pipeline
  [ ] Step 1:
    [ ] Discover Package from file
    [ ] Lex and Parse all modules in the package
    [ ] Repeat for all imports
  [ ] Step 2:
    [ ] Construct dependency graph between imports
  [ ] Step 3:
    [ ] Semantic analysis all files given the dependency order
  [ ] Step 4:
    [ ] Generate IR for all packages