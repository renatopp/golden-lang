# Flow Control

## If's

The common usage of Ifs are similar to other languages:

```
if x < 0 {
  println('Negative')

} else if x > 0 {
  println('Positive')

} else {
  println('Zero!')
}
```

The main difference is the if's can also be used as expression, in this case, it must always contain an `else` case, and return the same type:

```
let x = if value > 0 { 1 } else if value < 0 { -1 } else { 0 }
```

Ifs can also be used as empty switch cases:

```
let x = if {
  value > 0 ->  1
  value < 0 -> -1
  else      ->  0
}
```

## For's

For also have some different ways of usage. The most simple way is the 'forever':

```
for {
  ...
}
```

For can work like a `while` in other languages:

```
for condition {
  ...
}
```

They can also be used as expressions:

```
let x = for {
  break 42
}
```

Finally, they can be used as a `foreach` construct, very similar to what lodash does in their construction:

```
for value in iterator { ... }
for value, key in iterator { ... }
for value, key, i in iterator { ... }
```

## Pattern Matching

Pattern matching can be used as an statement, which makes it similar to a switch case:

```
match x {
  0 -> println(0)
  1 -> println(1)
  x -> println(x)
}
```

and used as expression, which makes it exhaustive:

```
let res = match (x%3, x%5) {
  (0, 0) -> 'FizzBuzz'
  (0, _) -> 'Fizz'
  (_, 0) -> 'Buzz'
  (_, _) _> String(x)
}
```

## With

With is a very special block that also adds a syntax sugar over anonymous functions:

```
let x = with { 1 }
-- is equivalent to
let x = (fn() Int { 1 })()
```