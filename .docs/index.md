# The Golden Language

<!-- TOC -->

- [The Golden Language](#the-golden-language)
  - [Design Pillars](#design-pillars)
  - [Naming Rules](#naming-rules)
  - [Packages, Modules and Imports](#packages-modules-and-imports)
  - [Visibility](#visibility)
  - [Types](#types)
    - [Primitives](#primitives)
    - [Functions](#functions)
  - [Value Declarations](#value-declarations)
    - [Functions](#functions)
  - [Expressions](#expressions)
    - [Arithmetic Expressions](#arithmetic-expressions)
    - [Comparison Expressions](#comparison-expressions)
    - [Logical Expressions](#logical-expressions)

<!-- /TOC -->

## Design Pillars

The heart of the language lies on the following designing pillars:

- **Consistency with ergonomy**: rules of the languages should be general and applied to all constructs, but shortcuts may be taken in order to improve quality of life.
- **Cohesion with simplicity**: future features must be developed over the old ones improving how they talk together, but it must keep the language simple and intuitive.
- **Safety with flexibility**: the constructs must find a balance between safety of the language and flexibility of development, providing a reliable tool to detect problems in compile time but still allowing fast prototyping and interation of the application development.

These pillars are not simply references to the syntax of the language, but they must be applied from the conception of ideas, going through the development of the compiler and finding its way to the applications itself, in both arquitecture and code.

## Naming Rules

- Files **must** be named in lower snake case (eg: `my_module.gold`), because their name reflects the module name in scope.
- Types **must** start with a capital letter (eg: `Type`), and **should** be in PascalCase, because types are not first class and they are treated fundamentally different from values.
- Values **must** start with a lower letter (eg: `value`), and **should** be in camelCase.
- Private names **must** start with an underscore (`_private`).

These naming decisions were taken to reduce the cognitive load of the programmer, so you can extract the maximum number of information just be the identifier alone, without resorting to jumping around the files looking for the declaration or relying on the intellisense for hints.

For example, just by looking at these generic names:

```
node             -- you now that's a value
Node             -- you now that's a type
_node            -- you now that's a private value
_Node            -- you now that's a private type
container.node() -- you now that's a function inside the value `container`
container.Node() -- you now that's a constructor inside the value `container`
```

## Packages, Modules and Imports

Golden uses the folder structure as basis for its package structure.

- A folder is a package.
- A file is a module.

Packages are merely intermediary namespaces to organize your application logically and cannot be used inside the code in any way. Modules on the other hand can be imported into the scope of another module, using its file name as the qualifier for its access. 

Modules can be imported by other other modules but they follow strict rules:

- Modules can import other modules in its scope, but always using the a qualifier to access the module content, i.e., you cannot inject a single type or value from another module directly into the scope.
- Modules inside the same package can access each other without explicity imports, but they still following the qualifier access.
- Modules can import modules from another package, but packages cannot have cyclic references, i.e., if modules from package `a` import modules from package `b`, then, modules from `b` cannot import modules from `a`.
- There is no relative imports, modules should always to imported using absolute paths.

These rules help us optimize the compilation time and analysis of your code, but they also force users to follow a similar structure from project to project. A golden developer should look for starting an application using a shallow structure and improving organization organically as the project grows. These rules also help this movement by keeping enabling simpler structural refactors.

The following example to illustrate this mechanics.

Consider the file structure:

```
project/
  main.gold
  internal/
    parser.gold
    token.gold
    node.gold
```

In `main.gold`, the entry point of our application, you find:

```rust
import '@/internal/parser'

fn main() {
  parser.parse('2 + 3')
}
```

where:

- `import '@/internal/parser'`:
  - `@` means the root of the project.
  - `internal` is a folders inside the root
  - `parser` is a file named `parser.gold`
- `fn main() {}`: is the entry point function
- `parser.parse('2 + 3')`:
  - `parser` is the module qualifier
  - `parse` is a function inside the module
  - `('2 + 3')` is the application if the string argument to the function

Optionally, you can change the module qualifier using the `as` keyword:

```python
import '@/internal/parser' as p
```

## Visibility

Types and values declared inside a module are exported by default unless they are prefixed with underscore (`_`). This rule is applied to modules inside a package.

For example, consider that we add a new file `_core.gold` inside the internal package:

```
project/
  main.gold
  internal/
    _core.gold
    parser.gold
    token.gold
    node.gold
```

Now, inside the `parser.gold` file we can use, let's say, `_core.randomId()`, but in `main.gold`, which lies in other package, you cannot import the private module.

## Types

### Primitives

Golden have very few primitive types:

- `Bool`: `true` and `false`
- `Int`: `1`, `-1`, `1_000_000`
- `Float`: `1.0`, `-1.0`, `1e10`
- `String`: `'Hello'`, `` `Hello` ``
- `Void`: special usable

### Functions

Functions have the type signature: `Fn(<Arg1Type>, ...) <ReturnType>`, where the return type is optional.

## Value Declarations

Values can be declared using the keyword `let` and they are immutable by default, i.e., their value cannot change and so their reference.

```rust
let int1 = 0
let int2 Int = 0
let int3 Int

let float1 = 0.0
let flaot2 Float = 0.0
let flaot3 Float

let bool1 = false
let bool2 Bool = false
let bool3 Bool

let str1 = ''
let str2 String = ''
let str3 String
```

All the declarations above for each type are equivalent.

An immutable variable cannot be reassigned, but it can be redeclared (shadowed):

```rust
let a = 0
let a = 'string'
```

### Functions

Functions can be created using the `fn` keyword and should follow strictly the signature: `fn(<param> <Type>, ...) <Type> { ... }`.

```rust
let add = fn(a Int, b Int) Int {
  a + b
}

-- or

fn add(a Int, b Int) Int {
  a + b
}
```

All the declarations above for each type are equivalent.

## Expressions

Pretty much everything in Golden is an expression except the statement `import`. Expressions return necessarily a value, which is attached to a type. Statements are special code instructions that don't return anything and cannot be used in expressions.

```rust
2 + 3 --> is evaluated to 5
```

Blocks, defined by the brace pair `{}`, are also considered an expression that returns a value:

```rust
{}                --> is evaluated to a Void value
{ 2 + 3 }         --> is evaluated to 5
{ 2 + 3; 5 + 10 } --> is evaluated to 15 (the previous expression is evaluated to a Void value)

let x = 2         --> is evaluated to 2
let f = fn() {}   --> is evaluated to a Fn() value
```

Notice that, by this rule, we don't need `()` to represent mathematical expressions, we can use `{}`. In fact, considering the block behavior, Golden reservers the behavior of parenthesis to other uses. For example:

```rust
2*{a + b}

x and {x or y}
```

### Arithmetic Expressions

```rust
a + b
a - b
a * b
a / b
a % b
```

### Comparison Expressions

```rust
a == b
a != b
a < b
a <= b
a > b
a >= b
a <=> b   --> spaceship operator, returns -1 when a < b, 0 when a == b and 1 when a > b
```

### Logical Expressions

```rust
a and b
a or b
a xor b
!a
```
