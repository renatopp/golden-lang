# Loops

## Context

I want to give users the ability to create loops that can process data in a traditional way, without resorting exclusively to recursion. For loops are the strong candidate to implement this feature.

There is some details in the implementation:

- Loops can be an expression?
- How loops can express with clarity their return value?
- Are loops really useful without variable mutation?

## Proposal

For loops can implement several loop cases, close to what Go does:

```rust
for { ... }
for <condition> { ... }
for <name> in <iterator> { ... }
```

Where `for {}` is an infinite loop, `for <condition> {}` is similar a traditional while block and `for <name> in <iterator> {}` follows Python statements.

