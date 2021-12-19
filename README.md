# error helper for deferred functions

What if you want to defer Close() on an io.WriteCloser, but you don't want to
throw away the error?

```go
// Always calls Close() and returns the Close() error if no writes fail
func writeAndClose(output io.WriteCloser, records [][]byte) (err error) {
	defer deferred.FuncErr(output.Close)(&err)
	for _, record := range records {
		_, err = output.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}
```

The named result parameter is necessary

## See Also

Joe Shaw: [Don't defer Close() on writable files](https://www.joeshaw.org/dont-defer-close-on-writable-files/)