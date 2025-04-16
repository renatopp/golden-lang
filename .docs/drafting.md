# Language Draft

Focus:

- Syntactical and semantic consistency
- Nil and Async safety

Other considerations:

- Iterator/Generator patterns
- Error handling: error codes, custom data, retries, recovery, timeout
- Logging: structured, levels, transporters, formatters
- Performance: caching, benchmarking, object pooling
- Security: input validation
- Async: thread pool, synchronization, queuing
- Data: serialization, conversion (DTO?)
- Testing: unit, mock, integration, fixtures, load
- Deploying: feature flags, monitoring, alerts, build tools
- GRPC
- Localization: datetime, timezone, currency, language

## Rules and Conventions

Comments starts with `--`.

Naming rules:

- `lowerCase` for public values
- `UpperCase` for public types
- `_under` for private values
- `_Under` for private types
- package name is the directory name
- files starting with _ non export to the package

## Type System

- Strong and Static typing with explicit notations
- Productivity with safety guards
- No implicit conversions, except by subtyping/interfaces

```
-- Primitives
Int      0
Float    0f
Bool     true
String   'string'
Void     ()
Byte     Byte(0)

-- Composite
Tuples   (1, 2)
Structs  {a:1, b:2}
Lists    []Int{1,2,3}
Maps     [String]Int{'a':1, 'b':2}

-- Custom Definitions
enum Value { Empty, Mine, Neighbor(Int) }
struct Cell { value Value }
interface Stringer { ToString() }
alias Number { Float, Int }
```

Golden uses ADT as a solution for removing nils and invalid states from code.

Type notation follows the Golang rules:

```
fn add(a Int, b Int) Int { ... }
-- is the same as:
fn add(a, b Int) Int { ... }
```

Generics are present and should be allowed in methods.

```
fn compute<T>(a, b T) T { ... }

struct MyStuff {
  fn compute<T>(a, b T) T { ... }
}
```

Type/Pattern matching:

```
match value {
  Struct { x:_, y:_, .. }
  Tuple(x, _, ..)
  _
}
```

In the future, we may add refinement types and row types.

Types are not first-citizen, thus, they cannot be used as value.

## Expressions

Golden is an expression-first language, which means that every block do generate a value:

```
let x = if condition { 5 } else { 10 }
let y = loop { break 42 }
let z = { a := 2; a*a }
fn add(a, b Int) Int { a + b }
```

## Examples

Load config file and process stuff.

```
result := with {
  os.read('file')!
  | json.parse<Config>()!
  | (x: x.tools)()
  | iter.map(x: load(x))
  | promise.all()
  | promise.wait()
  | result.errors()!
}

match result {
  Ok(v)  -> report_ok(v)
  Err(e) -> report_error(e)
}
```


```
opaque struct Object {
  _x Float
  _y Float

  fn x() Float { _x }
  fn set_x(v Float) { _x = v }

  fn y() Float { _y }
  fn set_y(v Float) { _y = v }

  fn position() (Float, Float) { return (_x, _y) }
  fn set_position(a, b Float) { _x = a; _y = b }
}

o := Object{}

```