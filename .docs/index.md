# Documentation

Only developed features.


## Statement vs Expressions

Statements are lines of code that performs an action but does not produce a value.
Expressions are lines of code that produces a values.

In the module level, golden only accepts statements:

```rust
import <package>
type <Name> = <type expression>
const <name> = <compile-time expression>
fn <name>(<params>) <Type> <block>
```

Inside a block, golden accepts only expressions.