# Error Handling

Errors are defined by an interface:

```
interface Error {
  error() String
}
```

## Results

Using ADT, we can return errors by using the `Result` type:

```
fn divide(a, b Int) Result<Int> {
  if b == 0 {
    return Err(errors.New('division by 0'))
  }
  return a/b
}

-- shortcut:

fn divide(a, b Int) Int! { ... }
```

## Immediate Returns

Instead of checking for errors after every action, you can simply propagate the error to be treated outside the function. Just prepend `!` to the expression:

```
fn load_config(path String) Config! {
  os.read(path)!
  | json.parse<Config>()!
  | validate()!
}
```

This is specially useful when used together with the `with` block:

```
fn load_config(path String) Config {
  let res = with {
    os.read(path)!
    | json.parse<Config>()!
    | validate()!
  }
  
  match res {
    Ok(v) -> c
    Err(e) -> {
      println('warning: error loading config, usign default config instead')
      defaultConfig()
    }
  }
}
```

## Panics

Panics cannot be recovered!