# Error Handling

```
type Result<T> = Ok(T) | Error(error() String)
type Option<T> = Some(T) | None

type Config(
    Server(
        host String
        port Int
    )
    Debug bool
)

fn main() {
    let file = os.readFile('config.json')
}





type Error = { error() String }
type Result<T> = T | Error

```