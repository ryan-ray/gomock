package main

type Foo interface {
	Bar() error
	Baz(s, z string) (string, error)
	FooBarBaz(a any) (int, float64)
}

type Store[T, R any] interface {
	Do(obj T) R
	Create() Bar
	CreateWithPointer() *Bar
}

type Bar struct{}
