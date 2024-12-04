# TODO

[ ] Implement functions
  [x] Implement function calls
  [x] Add return statement
  [x] Create pass structure
  [x] pass: add auto return
  [ ] pass: add block to function
  [ ] pass: add xor to function
  [ ] pass: add <=> to function
  [ ] Javascript:
    [x] Gen return
    [x] Gen function with blocks
    [ ] Gen function call
  [ ] Go:
    [ ] Gen function signature
    [ ] Gen type signature
    [ ] Gen return
    [ ] Gen functions with blocks
    [ ]Gen function call
  [ ] Test closure

[ ] Review string usage: change string from " to ' (?) and add automatic offset for multiline. Also add `` for raw string
[ ] Check parse errors `fn(a, b`
[ ] Implement module imports

## After functions

[ ] Implement xor operator into a lambda function
[ ] Convert <=> operator into a lambda function that is called inplace

    - `let x = a <=> b` => `let x = (:{ let a_=a;let b_=b; if a < b then -1 else if a > b then 1 else 0 })`

[ ] If the block contains more than 1 expression, and the block is within another expression (how to catch it?), the ast should convert the block into a lambda function with a direct call to it:

    - `let a = 1 + {2; 4}` => `let a = 1 + (:{ 2; 4 })() 

