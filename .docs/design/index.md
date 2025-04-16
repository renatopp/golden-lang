# Language Design Description

Golden language is being designed to """improve""" Golang in a general way. It tries to follow the same design principles adding more syntactical and semantical consistency, safety and ergonomics improvements. It removes completely the presence of `nil`s, which changes drastically the whole foundational structure of the language.

## Core Features

- Algebraic Data Types (ADT) as a form of eliminate `nil`s and invalid states.
- ADT-based error propagation, by using `Result` object and `!` operator to immediately return the error.
- Promise-like layer upon goroutines allowing more flexibility and control over async.
- Differentiating exception and panics (on is recoverable and other is not).
- Reference only, no more ambiguity between value and pointer with interfaces.

## Basics

- Modules
- Type checking
- Primitive types
- Composite types
- Assignments
- Expressions
- Blocks

## Functions

- High order
- Anonymous functions
- Function generics
- Captures
- Pipelines
- Documentation

## Flow control

- Ifs
- For loops
- Pattern matching
- With

## Data types

- Structs and Enums
- Struct Updates
- Generics
- Opaque Types
- Interfaces
- Alias

## Error Handling

- Error
- Results
- Immediate Return
- Exceptions
- Recovering
- Panics

## Async

- Promises
- Goroutines
- Channels
