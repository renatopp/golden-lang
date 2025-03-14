# 

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

## Default types

```
let a0 = 'Hello, World!';
let b0 = 42;
let c0 = 3.14;
let d0 = true;
let e0 = []Int{1, 2, 3};
let f0 = [String]Int{'a':a, 'b':2, 'c':3};
let g0 = (1, 'renato')

let a1 String = 'Hello, World!';
let b1 Int = 42;
let c1 Float = 3.14;
let d1 Bool = true;
let e1 []Int = [1, 2, 3];
let f1 [String]Int = ['a':a, 'b':2, 'c':3];
let g1 (Int, String) = (1, 'renato')

let a2 String
let b2 Int
let c2 Float
let d2 Bool
let e2 []Int
let f2 [String]Int
let g2 (Int, String)
```

## Capabilities

```
-- Option 1:
let a = []Int{1, 2, 3}
let *a = []Int{1, 2, 3}
let #a = []Int{1, 2, 3}
let &a = []Int{1, 2, 3}
fn sample(z, *a, #b, &c Int)

-- Option 2:
let mut a = []Int{1, 2, 3}
let box a = []Int{1, 2, 3}
let var a = []Int{1, 2, 3}
fn sample(z, mut a, box b, var Int) var Int
```

## Custom type definition

- ADT
- Interfaces
- Aliases

```
-- ADT
type Tree<T> = Empty | Node(T, Tree<T>, Tree<T>)
type User(
    id Int
    name String
    email String
    age Int
    active Bool
)

-- Interface
type Stringer {
    string() String
}

-- Alias
type Number = Int | Float
```

## Async

```
-- Promises
let calls = []Promise{
  async fetch('https://api.github.com/users/renatopp'),
  async fetch('https://api.github.com/users/r2pdev')
}
let results = await promise.all(calls)

-- Channels
let channel = Channel<Int>(10)
channel.send(42)
channel.receive()

-- Box / Async
let box counter = Box<Int>(0)
counter.set(42) -- atomic

unwrap(box)

```

## Error handling

```
fn withError() Result {
    return Error('Something went wrong')
}

fn main() {
    -- propagates de error up
    withError()!

    -- handle the error
    let result = withError()
    match result {
        Ok => println('Success')
        Error => println('Error')
    }
}
```

## Testing

```
fn testSample(runner) {
    assert.equal(...)
}
```

## Loops and conditionals

```
for value in iterator {}
for value, key in iterator {}
for value, key, index in iterator {}
for condition {}
for {}

if condition {} else if {} else {}

match value {
    1 -> ...
    2 -> ...
    _ -> ...
}

case condition -> ...
case {
    condition -> ...
    condition -> ...
    else -> ...
}
```

## Examples: Algorithms

### Quicksort

```
fn quicksort(arr []Int) []Int {
    if arr.length() <= 1 {
        return arr
    }
    let pivot = arr[arr.length() / 2]
    let left = arr.filter(x: x < pivot)
    let middle = arr.filter(x: x == pivot)
    let right = arr.filter(x: x > pivot)
    return quicksort(left) + middle + quicksort(right)
}

fn main() {
    let arr = []Int{3, 6, 8, 10, 1, 2, 1}
    let sorted = quicksort(arr)
    println(sorted)
}
```

### Sieve of Eratosthenes

```
fn sieve(limit Int) []Int {
    let sieve = List.repeat<Bool>(True, limit + 1)
    sieve[0] = False
    sieve[1] = False

    for i in Int.range(2, Int(math.sqrt(limit)) + 1) {
        if sieve[i] {
            for multiple in Int.rangeStep(i*i, limit + 1, i) {
                sieve[multiple] = False
            }
        }
    }
    
    let primes = []Int{}
    for prime, i in isPrime {
        if prime {
            primes.append(i)
        }
    }
    
    return primes
}

fn main() {
    let primes = sieve(50)
    println(primes)
}
```

### Fibonacci with Memoization

```
fn fibonacci(n Int, memo [Int]Int) Int {
    case {
        n in memo -> return memo[n]
        n <= 2    -> return 1
    }
    
    memo[n] = fibonacci(n - 1, memo) + fibonacci(n - 2, memo)
    return memo[n]    
}

fn main() {
    let result = fibonacci(10, [Int]Int{})
    println(result)
}
```

### Dining Philosophers Problem (Concurrency)
### Monte Carlo Pi Estimation (Parallelism)
### Search Algorithm

## Examples: Data Structures

### Trie
### LRU Cache
### Markov Chain
### Circular Buffer