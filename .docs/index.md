# Golden Language Overview

<!-- TOC -->

- [Golden Language Overview](#golden-language-overview)
  - [Comments](#comments)
  - [Expressions](#expressions)
  - [Functions](#functions)
  - [Modules](#modules)

<!-- /TOC -->

## Comments

Comments are prefixed by `--` and goes to the end of the line.

## Expressions

Everything in the language is an expression, which means that every command and computation has a type and a value. The most basic types are `Int`, `Float`, `String`, `Bool` and `Void`, but Void cannot be assigned to anything.

```go
const int = 1 + 1_000_000 - 100 
const float = 1.0 + 1e10 + -1e100
const bool = true or false
const string1 = "This is a string and
                 strings are multi line by default.
                 The ones starting with \" have automatic
                 offset, thus it will be evaluated without
                 the spaces in the left."
const string2 = `Raw strings are also multi line, but they
                 will be evaluate as is, with all spaces.`
```

Blocks are also expressions, and you can use blocks anywhere you would with other expressions. Blocks describe list of expressions and are evaluated to its last expression. If no expression is provided, block evaluates to `()`, which is a `Void` value.

```rust
const a = 1 * { 4 + 2 }
const b = 3 * { const x = 5; x*x }
```

Operations are:

```haskell
+a -- Int and Float
-a -- Int and Float
!a -- Bool
a + b -- Int, Float and String
a - b -- Int and Float
a * b -- Int and Float
a / b -- Int and Float
a % b -- Int
a < b   -- Int and Float
a > b   -- Int and Float
a <= b  -- Int and Float
a >= b  -- Int and Float
a <=> b -- Int and Float
a == b
a != b
a and b -- Bool
a or b  -- Bool
a xor b -- Bool
```

## Functions

Functions represents the behavior of the program. They can be declared in the module scope using a name. Functions must have the full declaration of types (see lambdas for non annotated declarations). Parenthesis are optional if the functions does not receive any parameter. 

```haskell
-- These are equivalents
fn func = <expr>
fn func() = <expr>

-- These are equivalents
fn func(a, b Int) Int = <expr>
fn func(a Int, b Int) Int = { <expr> }
```

Function call is pretty traditional: `func()`, `func(1, 2)`, etc. Function are typed as `Fn`, `Fn()`, `Fn(Int, Int) Int`, etc.

Inside expressions, functions cannot be named:

```rust
fn multiAdd(x Int) Fn(Int, Int) Int =
  fn(a, b Int) Int =
    x * (a + b)
```

## Modules

Modules are first class.