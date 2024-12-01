# Feature Reference

## Module

Modules are first class.

```rust
// repository
module[T](
  redis *Redis
)

fn save(e T) { redis.save(e) } 
fn get(id Int) { redis.get(e) }

//
import "repository"

let redis = redislib.new()
repository(redis).save(user)
// or
let redisRepository = redis(redis)
redisRepository.save(user)


```


