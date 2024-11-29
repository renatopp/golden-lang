# Mutability

Some notes abount mutability.

## Coding Patterns

Ignoring all discution about safety of mutability, using immutable variables and structures bring us some nice coding patterns that are specially useful for a language without methods, such as chaining functions with pipe expressions:

```
10
| add(2)
| rem(2)
| mult(2)
```

This is only possible because you are force to return the value of interest over and over:

```
data State = Hidden | Flagged | Revealed
data Cell(State)

fn reveal(cell Cell) Cell { Cell(Reaveled) }
fn hide(cell Cell) Cell { Cell(Hidden) }
fn flag(cell Cell) Cell { Cell(Flagged) }

Cell(Hidden)
| reveal()
| hide()
| flag()
```

The same pattern is not as useful when the programmer can use mutable variables as such:

```

fn reveal(mutable cell Cell) {
  cell.0 = Revealed
}

...
```

One possible solution to this is to use a strict ownership system where mutable variables can have a single reference to it, so, if you send the mutable variable to a function, that function must return it so you can used it again.

However, this brings a lot of complexity and hard problemas, such as what to do with cyclic references inside structs?


## Addressing vs Value Capability

We can define that a variable is a pair `{address, value}`, where the first is the memory address the second lives. In javascript, you can define address capability to a variable where `let` and `var` allow you to change the address a variable is pointing to, and `const` doesn't. However, in javascript, you can declare the capability of the value itself, e.g.: can this variable change the value underneath?

Working with immutable x mutable values, we may think in providing the mutablility capability to the value, so a variable could have the following capabilities `{READ|WRITE, READ|WRITE}`, resulting in 4 possible combinations. 

This is an interesting idea, but adding this level of detail seems to add unnecessary complexity to the language.


## Immutable vs Mutable Capabilities

My current option for the mutability system is to treat values as merely data in memory that don't have an inherent capability function. These permissions are only conceeded to the variables:

```
let x = 2     --> x is immutable {READ, READ}, ie: you can't reassign and you can't change its value
let mut x = 2 --> x is mutable (WRITE, WRITE), ie: you can reassign and change the underneath value
```

This simplifies the previous idea of address/value capability pairs.

```
data Point(x, y Int)

data Actor(
  name String
  mut position Point
)

let a = Actor(name: 'Renato', position: Point(2, 3))
a.position = Point(4, 5)
a.position.x = 1
a.position.y = 1
```

In the example above, I show a sideeffect of this system. Even though the variable a is immutable, the value it points contain a mutable field. Thus, allowing changing it's nested properties.

Mutable capability gives permission to the complete subtree. This seems bad, but when you combine this system with private fields and pure functions, event though a subfield from a mutable reference could be changed, you still need to use the immutable patterns. This looks like a good balance.

## Mutability and async

There are three main elements in the mutability vs immutability discussion:

1. Immutability patterns fit pipe expressions
2. Mutability introduce race conditions
3. Mutability can introduce logic errors

Point 2 is the main concern, and there are some possibilities to handle it. Point 3 is important but the rules applied to point 2 will help point 3.

There are multiple options to handle race conditions:

- Affine types where you can have only one reference to a value at a given type.
- Mutex/RWMutex
- STM (software transactional memory)
- CAS (compare-and-swap)
...

My idea is to provide a concurrency control on language level and check at compile time the validity of asynchronous routines.

- immutable: globally safe, value never changes 
- readable: locally safe, value may change if points to mutable
- mutable: locally safe
- shared(?): globally safe, can only be used when temporarily converted to readable/writable, which may happen using locks, stms, cas...


