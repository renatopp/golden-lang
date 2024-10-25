# Golden Programming Language

This project aims to create a programming language that is simple but expressive, that is easy to learn and write but also easy to read and understand, that supports complex programs but encorages simple arquitectures. 

The Golden Programming Language is a procedural, static and strong typed language based on Go and Gleam, borrowing some inspirations from other languages such as Rust, Austral and other functional languages.

> This repository contains ongoing work and most of the content here is dynamic or temporary accordingly to the development phase.

## The Language

At this point, this section works only for development reference.

### The Foundation

- Modules, Packages and Imports
  - Modules are files
  - Packages are folders
  - `@` denotes the main package (root package of the project)
  - Packages cannot have circular dependency
  - Modules inside the same package auto import other modules
  - `import <package>*/<module>` is the base syntax

- Variable Definitions
  - Variables can be named as the regex: `_?[a-z]([a-z][A-Z][0-9]_)*`
    - first letter MUST be lower case
    - optionally, it can have a single `_` as first character
  - Variables are declared as:
    - `let <name> = <expression>` with compiler deciding the type
    - `let <name> <typeref>` with default initialization
    - `let <name> <typeref> = <expression>` with complete information
  - Variables are immutable by default, thus:
    - It cannot have reassignment
    - It can be passed as argument
    - It can be redeclared and shadowed

- Functions
  - Functions can have the same name as variables
  - Function declaration follows:
    - `fn <name>(<params>) <return> { <expression> }`
    - `fn (<params>) <return> { <expression> }`
    - params are a list of:
      - `<name> <typeref>`
    - return is a `<typeref>`
  - Functions can be called
  - Functions can be used as variable
  - Functions must have closure

- Types
  - Basic types: `I64`, `F64`, `Bool` and `String`

- Temporary Builtin Functions
  - `debug(string)` and `debugln(string)`

### Planned Features

- Type definition: struct
- Type definition: tuple
- Type definition: untagged union
- Type definition: tagged union
- Mutable variables
- Generics
- Result object (with `?` and `!`)
- Short parameters (`<name>, <name> <type>`)