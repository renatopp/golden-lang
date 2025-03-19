# Language Draft

Focus:

- Iterator/Generator focus
- Functional programming patterns
- Async safety
- Syntactical and semantic consistency

Other considerations:

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

Basis:

- Function only - no method allowed!
- ADTs
- Immutability as default

## Custom Data Types

### ADTs

```
-- General form:
type TypeName = Constructor1 | Constructor2 | ...

-- Unit definition
type Unit = A | B()

-- Tuple definition
type Tuple = C(Int) | D(String, Int)

-- Record definition
type Record = E(a Int, b String) | F(a Int)

-- Shortcut
type FullLength = FullLength(Int, String)
type Shortcut(Int, String)
```

### Type Alias

```
-- When using existing types
type Number = Int | Float
```

### Generics

```
type Tree<T> = Leaf(T) | Node(T, Tree<T>, Tree<T>)
```

## Functions and Traits

### Basic Functions

```
fn main() { ... }
fn foo(a Int, b String) Int { ... }
fn sum(a, b Int) Int {
    return a + b
}
```

### Open Functions and Traits

```
-- json.gold file
open fn toJson<T>(value T) String
open fn fromJson<T>(json String) T { ... } -- default implementation
trait Serializable { toJson, fromJson }
```

```
-- main.gold file
type MyType

impl json.toJson(value MyType) String {
    return "..."
}

impl json.fromJson(json String) MyType {
    return MyType
}

fn run(s json.Serializable) {
    json.ToJson(s)
}
```

## Capabilities

```
let a = 1 -- immutable
let *b = 2 -- mutable
let @c = 3 -- wrapped

type Struct(
    a Int
    *b Int
    @c Int
)
```

## Async

```
-- run function asynchroniously
fn func() Int {
    return 1
}

let p = async func() -- returns a Promise<Int>
let result = await p -- result is an Int

-- communicate via channel
fn consumer(c Channel<Int>) {
    for i in c {
        print(i)
    }
}
async consumer()

let c = Channel<Int>(10)
c.send(1)

-- access shared memory
let @a = 0

for i in int.range(10) {
    async fn() {
        unwrap a as c {
            c += 1
        } -- how to define unwrapping with read only? imu as read, mut as write?
    }()
}
```

## Error Handling

```
fn open() Result<String, Error> { ... }
```

## Loops and Conditionals

```
for {} -- inifinite loop
for condition {} -- while loop
for i in range(10) {} -- iterator loop

if condition {} 
else if {}
else {}

if {
    condition -> ...
    condition -> ...
    else      -> ...
} -- switch case?

match value {
    1 -> ...
    _ -> ...
} -- pattern matcher

```

## Practical Examples


### Async Fetch Http

You need to build a function that takes a list of API endpoints, fetches data from all of them concurrently, and
returns a consolidated JSON object. Each API returns a JSON response with a data field (e.g., { "data": [1, 2, 3] }).
The function should combine the data arrays from all responses into a single array, sort it in ascending order, and
return the result as a JSON object with a result field (e.g., { "result": [1, 1, 2, 3, 4, 5, ...] }). Handle errors
gracefully (e.g., skip failed APIs) and ensure the process is optimized for performance by fetching data concurrently.
Optionally, add a timeout and retry mechanism for robustness.

```
import 'fetch'

type Response(
    data List<Int>
)

let endpoints = List<String>([
    "https://api.example.com/data1",
    "https://api.example.com/data2",
    "https://api.example.com/data3"
])

fn main() {
    endpoints
    | list.map(url: async retrieve(url))
    | await promise.all()
    | list.flatten()
    | list.sortBy(int.compare)
}

fn retrieve(url String) List<Int> {
    fetch.get(url)
    | result.mapOk(res: res.body)
    | result.mapOk(json.parse<Response>) 
    | result.mapOk(res: res.data)
    | result.or(err: List<Int>())
}
```

### 

```

```