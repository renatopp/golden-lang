# Functions

Functions are also quite similar to Go, but they are declared with `fn` keyword. They are first-class citizens and can be anonymous:

```
fn add(a, b Int) Int {
  return a + b
}

fn double(func Fn(Int, Int) Int) Fn(Int, Int) Int {
  fn (a, b Int) Int { -- anonymous functions
    func(a, b) * 2
  }
}

fn main() {
  let d = double(add)
  d(2, 4) -- 12

  add(2, 4) -- 6
}
```

## Generics

Functions support parametric polymorphism:

```
fn twice<T>(v T, f Fn(T)) {
  f(v)
  f(v)
}
```

By default, the generic type has no constraint as default, but you can pass interfaces that limit the T:"

```
fn serialize<T json.Serializable>(v T) T {
  ...
  return v
}
```

## Captures

Function captures is a syntax sugar of anonymous functions:

```
fn add(a, b Int) Int { a + b }

fn main() {
  let add_one = fn(a Int) Int { add(a, 1) }
  let add_two = add(_, 2) -- capture!
}
```

## Pipelines

The pipe operator is a form of chain different functions sequentially, passing the return of the previous function as implicit parameter to the next:

```
fn main() {
  add(1, 2)
  | add(3)
  | add(4)
  | println() -- prints '10'
}
```

## Lambdas

Lambdas are a shortcut for anonymous functions when they are explicit used as values, for example:

```
fn each<T>(list []T, f Fn(T)) { ... }

fn main() {
  list := []Int { 1, 2, 3 }
  each(list, x: println(x))
}
```

Lambdas follow losely the javascript syntax:

```
:0 -- same as fn() { 0 }
x: x + 1 -- same as fn(x) { x + 1 }
(x, y): x + y -- same as fn(x, y) { x + y }
```