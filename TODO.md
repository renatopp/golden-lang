# TODO

[ ] error with nested return `return { return 0 }`
  [ ] add state machine for feature toggle (with Inherit, Allow, Deny per feature)
[ ] if then else
  [ ] [optimization] return if -> if return else return (rodar depois do auto return)
  [ ] [optimization] transformar if em variavel temporaria pois `let x = if ...`
[ ] [optimization] add block to function `let a = 1 + {2; 4}` => `let a = 1 + (:{ 2; 4 })() 
[ ] [optimization] add xor to function
[ ] [optimization] add <=> to function `let x = a <=> b` => `let x = (:{ let a_=a;let b_=b; if a < b then -1 else if a > b then 1 else 0 })`
[ ] Review string usage: change string from " to ' (?) and add automatic offset for multiline. Also add `` for raw string
