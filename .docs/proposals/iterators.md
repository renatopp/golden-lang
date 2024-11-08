# Iterators

## Context

Iterators are a well-established feature in modern programming languages, providing a natural way to streamline data processing, including handling potentially infinite data sequences. The goal of this module is to offer a lazy data processing mechanism, allowing efficient, on-demand computation. This module will also paviment the road for future features, such as Python-like for loops and support for synchronous and asynchronous data processing.

## Proposal

```
type Iterator<T>(Fn() Iteration<T>)
type Iteration<T> =
| Stop
| Iteration(T)

fn range(to Int) Iterator<Int> {
  let *x = -1
  Iterator(fn() Iteration<Int> {
    x += 1
    if x < to { Iteration(x) } else { Stop }
  })
}

...
```
