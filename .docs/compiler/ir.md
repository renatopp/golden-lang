# Golden Intermediate Representation

GIR provide a subset of Golden, or as close as possible to Golden, that represents more "purely" an imperative language, converting complex, abstract or shortcut structures into simple flat definitions. The goal is to represent the source language into a version of it that is more compatible to C, Go, assembly or bytecode.

The GIR subset is:

## Modules vs packages

Package is now the only superstructure making modules lose their meaning. In practice, modules are merged into a single block of code. 

## Variables

Variables are defined by `let x0 = 0`. GIR uses SSA, so every variable has only one assignment. The only exception to SSA form is mutable variables that are assigned from an inner scope, which is mutated by `set x0 = 3`. SSA help us handle redeclaration in a compatible way.

## Functions

Functions are defined only in the root level of the package, thus, it cannot be nested. It should follow the complete form `fn(param Type, param Type) Type {}`, with explicity returns.

## Data types

ADT should keep its meaning.