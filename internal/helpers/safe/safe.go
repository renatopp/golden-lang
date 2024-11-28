package safe

type Optional[T any] struct {
	value   T
	present bool
}

func (o *Optional[T]) Or(v T) T {
	if o.present {
		return o.value
	}
	return v
}

func (o *Optional[T]) Has() bool {
	return o.present
}

func (o *Optional[T]) Unwrap() T {
	if !o.present {
		panic("unwrap of None")
	}
	return o.value
}

func (o *Optional[T]) If(fn func(T)) {
	if o.present {
		fn(o.value)
	}
}

func (o *Optional[T]) IfElse(fn func(T), elseFn func()) {
	if o.present {
		fn(o.value)
	} else {
		elseFn()
	}
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{value, true}
}

func None[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}
