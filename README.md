# Exire - General purpose TTL

Exire is a general purpose TTL (time to live) eviction library. This means that
Exire unlike most expiring caches can be used to expire things other than cache
items. This may be items on disk, in a database, or something else entirely.

## Example

```go
package main

func main() {
  lockingMu := sync.Mutex{}

  ttl := exire.New(
    exire.WithTTL(30 * time.Minute),
    exire.WithPreLock(func(ctx context.Context) error {
      lockingMu.Lock()
      defer lockingMu.Unlock()
    })
  )
}
```
