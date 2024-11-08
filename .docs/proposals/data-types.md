# Data Types

## Context

User defined types is a basic requirement. Algebraic data types is the main building blocks because it is a generalized form for structs and tuples.

The main guidelines for data types:

- Functions-only, no methods! This is a requirement to keep a consistent single-level organization of the module, thus:
  - Data types represents only data.
  - Functions can be chained in a future `|` operator.
  - Incentive pure functions for state-less operations.
  - Easy of refactoring and reorganization.

## Proposal


```
data <TypeName/Constructor>
data <TypeName> = <Constructor>
data <TypeName> = <Constructor> | <Constructor> | ...
```

where constructors:

```
<Name>
<Name>()
<Name>(<Type>, <Type>, ...)
<Name>(<var> <Type>, <var> <Type>, ...)
```