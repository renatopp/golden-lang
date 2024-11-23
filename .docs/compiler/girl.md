# Golden Intermediate Representation Language

GIRL is high level intermediate language used to represent the Golden language in a simpler way, removing complex data structures and simplifying flows of controls in a way that is more compatible to other languages, in particular, C.

Main characteristics:

- Package only, there is no notion of modules anymore.
- Type and function definitions occur only in the package scope (flatify types and functions)
  - Closures and internal type definitions should be extracted
- Names are treated as SSA
- Expressions are broke down into single operator ones
- No more information about mutability