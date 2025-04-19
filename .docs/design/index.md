# Language Design Description

Golden language is being designed to """improve""" Golang in a general way. It tries to follow the same design principles adding more syntactical and semantical consistency, safety and ergonomics improvements. It removes completely the presence of `nil`s, which changes drastically the whole foundational structure of the language.

## Core Features

- Algebraic Data Types (ADT) as a form of eliminate `nil`s and invalid states.
- ADT-based error propagation, by using `Result` object and `!` operator to immediately return the error.
- Promise-like layer upon goroutines allowing more flexibility and control over async.
- Differentiating exception and panics (on is recoverable and other is not).
- Reference only, no more ambiguity between value and pointer with interfaces.
- Generics are designed from beginning, so they can be used in methods and our notation makes the code cleaner.
- Fix the ugly mascot.

## Details

- [Basics](./01_basics.md)
  - Rules Conventions
  - Modules
  - Type checking
  - Primitive types
  - Composite types
  - Assignments
  - Expressions
  - Blocks

- [Functions](./02_functions.md)
  - High order
  - Anonymous functions
  - Function generics
  - Captures
  - Pipelines
  - Documentation
  - Lambdas

- [Flow Control](./03_flow-control.md)
  - Ifs
  - For loops
  - Pattern matching
  - With

- [Data types](./04_data-types.md)
  - Structs and Enums
  - Struct Updates
  - Generics
  - Opaque Types
  - Interfaces
  - Alias

- [Error Handling](./05_error-handling.md)
  - Error
  - Results
  - Immediate Return
  - Panics

- [Async](./06_async.md)
  - Promises
  - Goroutines
  - Channels
