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

```
int1 := 1_000_000
int2 := -10000000
int3 := 0b1010
int4 := 0o1000
int5 := 0x10ff

float1 := 1.0
float2 := .0
float3 := 1e10
float4 := 1f
float5 := -1e10

string1 := 'golden'
string2 := '⚠️'

bool1 := true
bool2 := false

byte1 := 120b
```

## Composite Types

- Lists
- Maps
- Tuples
- Structs
- Tagged Unions

```
lists   := []Int{ 1, 2, 3 }
maps    := [String]Int{ 'a': 1, 'b': 2 }
tuples  := (1, 2, 'renato')
structs := Vector2{ x:0, y:0 }
unions  := Ok(0)               -- Result<Int>
```

## Assignments

