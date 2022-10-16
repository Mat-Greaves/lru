# LRU

The `lru` package provides a go implementation of a generic fixed size, thread-safe LRU cache.

## Documentation

Full docs are availabel on [godocs](https://pkg.go.dev/github.com/Mat-Greaves/lru)

## Example

```
c, _ := New[int, int](128)
for i := 0; i < 256; i++ {
  c.Add(i, i)
}
if c.Len() != 128 {
  panic(fmt.Sprintf("bad len: %d", c.Len()))
}
```
