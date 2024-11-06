# Golden Programming Language

This project aims to create a programming language that is simple but expressive, that is easy to learn and write but also easy to read and understand, that supports complex programs but encorages simple arquitectures. 

The Golden Programming Language is a procedural, static and strong typed language based on Go and Gleam, borrowing some inspirations from other languages such as Rust, Austral and other functional languages.

> This repository contains ongoing work and most of the content here is dynamic or temporary accordingly to the development phase.

## Current State

```mermaid
flowchart LR

classDef uncertain fill:#F05D5E,color:black,stroke-width:0
classDef fair fill:#F7D002,color:black,stroke-width:0
classDef certain fill:#018E42,color:white,stroke-width:0
build[Build]:::uncertain
lexer[Lexer]:::certain
parser[Parser]:::certain
analyser[Analyser]:::fair
backend[Backend]:::uncertain

build --> lexer
lexer --> parser
parser --> analyser
analyser --> backend
```

| Step     | Progress                                        | Description                                           |
|----------|-------------------------------------------------|-------------------------------------------------------|
| Build    | <span style="color:#F05D5E">uncertain</span>    | Package loading, module loading, caching, etc.        |
| Lexer    | <span style="color:#018E42">certain</span>      | Lexer working as intended.                            |
| Parser   | <span style="color:#018E42">certain</span>      | Parser working as intended.                           |
| Analyser | <span style="color:#F7D002">fair certain</span> | Testing type checks and type inference.               |
| Backend  | <span style="color:#F05D5E">uncertain</span>    | Unsure about which backend to use. Probabily using C. |

## The Language Foundation

> For now, this section is for development reference only.

Design Pillars:

- **Consistency With Ergonomy**
- **Intuitive Simplicity**
- **Safety With Flexibility**

The foundation is a core set of features that defines the minimum base of the language, which will support all following developments. They should be simple to expand and useful enough to be used without additions.

### Modules, Packages and Imports
  - Modules are files
  - Packages are folders
  - `@` denotes the main package (root package of the project)
  - Packages cannot have circular dependency
  - Modules inside the same package auto import other modules
  - `import <package>*/<module>` is the base syntax

### Visibility

Definitions (both types and variables) starting with `_` are private and can only be accessed inside the module it was declared. 

### Expressions

Expressions are terms that will be evaluated to a value and everything in the language is an expression, however, some just evaluates to `void`.

Anywhere a term asks for an expression, you can use a block instead, i.e. `{}`, which will return the last list executed. Empty blocks also evaluates to `void`.

Notice that common mathematical expressions should use `{}` instead of `()`:

```rust
{x*x} + {y*y}
```

### Bindings

Bindings assign an expression to a name, which will be registered in the scope.

```rust
let pi = 3.1415      // Inferred type
let time Float = 10  // Complete information
let name String      // Default initialization
```

Variables are immutable by default and cannot be reassigned. However, they may be redeclared:

```rust
let x = 1
let x = "renato"
```

### Functions

Functions can be declared as:

```rust
fn add(a Int, b Int) Int { a + b }

add(1, 2)
```

Functions declared with a name inserts the name into the scope and is equivalent to:

```rust
let add = fn(a Int, b Int) Int { a + b }
```

Function type notation is a bit special:

```rust
let mult = fn(x Int) Fn(Int, Int) Int {
  fn (a Int, b Int) { {a + b}*x }
}
```

### Types

Golden have the following primitive types:

- `Void`
- `Bool`
- `Int`
- `Float`
- `String`
- `Fn`

Custom types can be declared using Algebraic data types:

```rust
type Weekday = Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday
```

where `Weekday` is the type and `Monday`, `Tuesday`, etc... are the constructors for that type. For example, calling 

```rust
let x = Monday
```

x is of type `Weekday` and its value is `Monday`.

Constructor may use structs or tuples:

```rust
type Vector = 
| Vector1(Int)       // 1-tuple
| Vector2(Int, Int)  // 2-tuple

let vec = Vector2(0, 0)

type Point = 
| Point1(x Int)        // struct
| Point2(x Int, y Int) // struct

let pnt = Point2(x=0, y=0)
```

Sum types may have a simple element:

```rust
type Dollar(Float)
type Euro(Float)
```

If a constructor have the same name as the type and is the only element, it can be written in a simplified way:

```rust
type Person(name String, age Int)
```

Anonymous types can be declared and used as:

```rust
fn bounds() (Int, Int) {
  (0, 10)
}

fn rect() (x Int, y Int, w Int, h Int) {
  (x=0, y=0, w=0, h=0)
}
```

