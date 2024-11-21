# TODO

clear && golden run --debug --working-dir=.samples/imports-cyclic/ .samples/imports-cyclic/a/a.gold
clear && golden run --debug --working-dir=.samples/imports-non-cyclic/ .samples/imports-non-cyclic/a/a.gold

## Step 1
[x] parse types
[x] parse function types
[x] parse function declarations
[x] parse function applications

## Step 2
[x] validate if all imports are valid
[x] create dependency graph
[x] check for circular imports

## Step 3
[ ] add type annotation in the node
[ ] create the scopes
[ ] create the module type
[ ] pre resolve all modules
[ ] resolve all modules for each package

## Step 4
- create IR