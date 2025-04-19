# Basics

Here's a hello world in Golden:

```golden
fn main() {
  println('Hello, World!')
}
```

## Rules and Convention

Rules:

- `lower` lower first letter for values (and modules).
- `Upper` upper first letter for types.
- `_lower` lower first alpha, prefixed by underscore for private values.
- `_Upper` upper first alpha, prefixed by underscore for private types.
- `_` for ignore and wildcard patterns.

Conventions:

- `snake_case` for value names (?).
- `CamelCase` for type names.

Comments:

- Comments are always prefixed by `--`.

## Modules and Imports

Modules in Golden follow the same structure and rules of Golang:

- Modules are defined by folders.
- Files are "merged" into the same module definition.
- Modules cannot have circular dependency.
- Imports works upon the module, there is no unqualified imports.

```
import 'mypackage/name'

name.stuff()
```

## Type Checking

Types in Golden are strong, static and explicit. Pretty much anywhere you need to declare types in Go, you will need in Golden too. Conversion should alway be explicit.

## Primitive Types

There are a few primitive types:

- Int
- Float
- String
- Bool
- Byte
- Void

```
let int1 = 1_000_000
let int2 = -10000000
let int3 = 0b1010
let int4 = 0o1000
let int5 = 0x10ff

let float1 = 1.0
let float2 = .0
let float3 = 1e10
let float4 = 1f
let float5 = -1e10

let string1 = 'golden'
let string2 = '⚠️'

let bool1 = true
let bool2 = false

let byte1 = 120b
let void = () -- called unit value
```

## Composite Types

- Lists
- Maps
- Tuples
- Structs
- Tagged Unions

```
let lists   = []Int{ 1, 2, 3 }
let maps    = [String]Int{ 'a': 1, 'b': 2 }
let tuples  = (1, 2, 'renato')
let structs = Vector2{ x:0, y:0 }
let unions  = Ok(0)               -- Result<Int>
```

## Assignments

Values can be declared using the `let` keyword. Variables are mutable and may be reassigned at any time. The `const` keyword can be used to declare values that cannot be reassigned, however, the reference may still be mutable.

```
const name1 = 'Golden'
const name2 String = 'Golden'
const name3 String -- empty string

let x1 = 0
let x2 Int = 0
let x3 Int -- 0

name1 = 'renato'
x2 = 10
```

## Expressions

Golden will not support bitwise operator initially, their syntax will be between studied later. Most other operators are implemented:

```
+a
-a
a + b
a - b
a / b
a * b
a ^ b
a % b

a < b
a > b
a <= b
a >= b
a <=> b

a == b
a != b

!a
a or b
a and b
a xor b
```

Excepting by imports and comments, everything is an expression in Golden. See more details for each feature.

```
let x = if condition { 0 } else { 1 }
```

## Blocks

At any time, you can start a block using `{}`, notice that blocks generates a value, which will be the last expression. If the block is empty, it will return an unit `()` by default.

```
let {
  let a = 1
  let b = 2
  a + b
}
```