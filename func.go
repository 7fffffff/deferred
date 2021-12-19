package deferred

func FuncErr(fn func() error) (deferMe func(*error)) {
	return func(errPtr *error) {
		err := fn()
		if *errPtr == nil {
			*errPtr = err
		}
	}
}
