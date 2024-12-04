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

```rust
let int = 1 + 1_000_000 - 100 
let float = 1.0 + 1e10 + -1e100
let bool = true or false
let string1 = "This is a string and
               strings are multi line by default.
               The ones starting with \" have automatic
               offset, thus it will be evaluated without
               the spaces in the left."
let string2 = `Raw strings are also multi line, but they
               will be evaluate as is, with all spaces.`
```

Blocks are also expressions, and you can use blocks anywhere you would with other expressions. Blocks describe list of expressions and are evaluated to its last expression. If no expression is provided, block evaluates to `()`, which is a `Void` value.

```rust
let a = 1 * { 4 + 2 }
let b = 3 * { let x = 5; x*x }
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
fn func() Void = <expr>

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

Modules can be defined in two ways: by file and by explicit declaration.

In Golden, every file represents a module and the file name will be used as name to the module. Thus, the file name must follow the same naming rules ([a-z_][a-zA-Z0-9_]*). Explicit declarations lets you create submodules inside the file module.

```
module name(param Type, param Type) = {
  ...
}
```

where params follow the same rules of function parameters.

Inside a module, you can declare types, functions and variables, and you can use them before declaration, because declaration in module-level does not follow a strict order.

Modules export everything as default unless it starts with `_`, which denotes private field.

You can import module files with:

```
import "@/sample/x/y/foo"
import "@/sample/x/y/bar" as baz
```

