# Feature Map

# Compiler

[ ] Execute `main` function
[ ] Assert function `@assert(bool, message)`

# Expressions and Types

[ ] Process operations (+, -, /, *, **, //, <, >, <=, >=, <=>, and, or, xor, !, ==, !=)
[ ] Int, Float, String, Bool, Byte


# Variables

[ ] Variable declaration with initialization `let x = 1` 
[ ] Variable declaration with default `let x Int`
[ ] Variable declaration with casting `let x Float = 1`
[ ] Conversion `let x = 1; let y Float = x`

# Functions

[ ] Function declaration `fn main() {}`
[ ] Function declaration with arguments `fn add(a Int, b Int) Int { return a + b }`
[ ] Second order function `fn plus2(f Fn(Int, Int) Int) Int { return 2 + f(2, 5)}`
[ ] Closure `fn multier(n Int) Fn(Int, Int) Int { return fn(a Int, b Int) Int { n * (a + b) } }`
[ ] Shortcut declaration `fn add(a, b Int) Int { return a + b }`
[ ] Partial application `let adder = add(_, 2); add(5) == 7`
[ ] Default values `fn triple(a=0, b=1, c=2 Int) Int { a + b + c }`
