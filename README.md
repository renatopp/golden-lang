# Golden Programming Language

This project aims to create a programming language that is simple but expressive, that is easy to learn and write but also easy to read and understand, that supports complex programs but encorages simple arquitectures. 

The Golden Programming Language is procedural static and strong typed strongly language based on Go and Gleam, borrowing some inspirations from other languages such as Rust, Austral and other functional languages.

> This repository contains ongoing work and most of the content here is dynamic or temporary accordingly to the development phase.

## TODO

## Language Overview

### Modules and Organization

In golden every file represents a module and every folder represents a package. Packages cannot have circular dependency, and modules inside the same package don't have to import each other.

```
project/          -- root is already the package `@`
  main.gold       -- package @ | module main
  utilities/
    random.gold   -- package @/utilities | module random
    lists.gold    -- package @/utilities | module lists
```

Modules should be named in `snake_case` and the language enforces it.

Imports always insert the whole module in the scope, never a package and never the content inside the module directly. After importing, you can access the module based on its name.

```
import @/utilities/random

random.stuff()
```

This organization is intended to encourage flat hierarchies without too much abstraction (or indirection) in the project layers.

### Types and Variables

Types must start with an uppercase letter and recommended to be in `PascalCase`. Variables must start with a lowercase letter and recommended to be in `camelCase`.

Simples types include:

- Integers: I8, I16, I32, I64
- Unsigned Integers: U8, U16, U32, U64
- Floats: F32, F64
- Char
- String
- Byte
- Bool
- Function

Composite types include:

- Lists
- Maps
- Structs
- Tuples
- Aliases
- Union: Number, Int, UInt, Float

Golden uses an algebric type system and have an universal definition syntax:

```
type NamedTuple = (Int, Int, Bool)
type NamedStruct = (scoreA, scoreB Int, result Bool)
type ListAlias = List<Int>
type IntUnion = I8 | I16 | I32 | I64
type NamedUnion = Bomb | Cell | Value(Int)

type NamedTupleShort(Int, Int, Bool)
type NamedStructShort(scoreA, scoreB Int, result Bool)
```

Initialization patterns follow some strict rules. Heterogeneous fixed sized structures (known at compile time) uses `()` and homogeneous dynamic sized structures (dynamic in runtime) uses `[]`. Unlabelled structures (tuples and lists) uses a list of elements `a, b` while labelled (structs and maps) uses pair key-value: `a=1, b=2`:

```
let namedTuple = NamedTuple(1, 2, true)
let anonTuple = (1, 2, true)
let namedStruct = NamedStruct(scoreA=1, scoreB=2, result=true)
let anonStruct = (scoreA=1, scoreB=2, result=true)

let list = List<Int>([1, 2, 3])
let map = Map<String, Int>([a=1, b=2])
```

#### Ownership

TODO

#### Visibility

Public as default, private using `_` as prefix. Works for both types and variables. Inside the same module, you can access public and private elements.

```
type User(
  _private Int
  public   Int
)

let _private Int
let public Int
```

### Functions!

Functions are first class and follow the same variable rules. The general form of a function is:

```
let function = fn (scoreA, scoreB Int, result Bool) Bool { ... }
```

In module scope, we provide the shortcut for named functions:

```
fn function(...) {}
```

#### Partial Application

Partial application can be done using:

```
let add = fn(a, b Number) Number { return a + b }
let add2 = add(_, 2)
```

#### Lambdas

Lambdas can be defined by:

```
:0        // javascript equivalent: `() => 0`
a:0       // javascript equivalent: `a => 0`
(a, b): 0 // javascript equivalent: `(a, b) => 0`
```

#### Pipelines

```
value
| functionCall()
| otherFunction()
```

#### Labelled Arguments

```
fn add(a, b Int) Int { ... }

add(b=2, a=2)
```

#### Generics

```
fn add<T any>() { ... }
```

### Flow Controls

```
if expression {
} else if expression {
} else {
}


for {}
for expression {}
for a in iterator {}

match expression {
  pattern: expression
  pattern: expression
}
```