# TODO

[ ] If the block contains more than 1 expression, and the block is within another expression (how to catch it?), the ast should convert the block into a lambda function with a direct call to it:

    - `let a = 1 + {2; 4}` => `let a = 1 + (:{ 2; 4 })() 

[ ] Review string usage: change string from " to ' (?) and add automatic offset for multiline. Also add `` for raw string