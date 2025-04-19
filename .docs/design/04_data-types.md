# Data Types

Golden uses Algebraic Data Types (ADT's):

```
enum State { Hidden, Revealed, Flagged }

struct Cell {
  state State

  fn reveal() {
    state = Revealed
  }

  fn flag() {
    state = match state {
      Hidden  -> Flagged
      Flagged -> Hidden
      _       -> Revealed
    }
  }
}
```

Notice that methods are declared INSIDE the type declarations. Compared to Go, this makes the usage of Generics a lot cleaner.

## Generics

Similar to functions, you can use generics in types.

```
struct Frame<T> {
  index Int
  value T
}
```

Another important aspect is the generics can be used in methods:

```
struct Ease {
  fn linear<T>(x T) T {
    return x
  }
}
```

## Opaque Types

Opaque types are similar to gleam: they cannot be instantiate outside of the module, adding more control to the developer:

```
opaque struct Vector {
  _x Int
  _y Int
}

fn NewVector() Vector {
  Vector{_x:0, _y:0}
}
```

## Struct and Enum updates

```
struct Vector{
  x Int
  y Int
  z Int
}

let a = Vector{}
let b = Vector{x:1, ...a}
```

## Interfaces

Same as golang.

```
interface Stringer {
  ToString() String
}
```

## Type Alias

Alias can be used as unions:

```
alias Number {
  Int
  Float
}
```