# TODO

clear && golden run --debug --working-dir=.samples/imports-cyclic/ .samples/imports-cyclic/a/a.gold
clear && golden run --debug --working-dir=.samples/imports-non-cyclic/ .samples/imports-non-cyclic/a/a.gold

[ ] add package path to module type and ast
[ ] ir converter should keep track of package names
[ ] 