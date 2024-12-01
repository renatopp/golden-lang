# TODO

[ ] Implement golden build command to generate the binary
[ ] Convert <=> operator into a lambda function that is called inplace

    - `let x = a <=> b` => `let x = (:{ let a_=a;let b_=b; if a < b then -1 else if a > b then 1 else 0 })`

[ ] If the block contains more than 1 expression, and the block is within another expression (how to catch it?), the ast should convert the block into a lambda function with a direct call to it:

    - `let a = 1 + {2; 4}` => `let a = 1 + (:{ 2; 4 })() 

[ ] Review string usage: change string from " to ' (?) and add automatic offset for multiline. Also add `` for raw string
[ ] Review usage of consts in modules

[ ] Implement module imports
[ ] Implement functions