# TODO

clear && golden run --debug --working-dir=.samples/imports-cyclic/ .samples/imports-cyclic/a/a.gold
clear && golden run --debug --working-dir=.samples/imports-non-cyclic/ .samples/imports-non-cyclic/a/a.gold

[ ] try codegen to go
  [ ] builder must configure temp output directory
  [ ] codegen must write the main file
  [ ] builder must call go build/go run
  [ ] builder must copy the binary to the final output 

  [ ] codegen must create the package folder in the temp
  [ ] codegen must create the package file in the temp
  [ ] codegen must write the imports in the file
  [ ] codegen must write the function declaration
  ...