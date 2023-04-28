package main

type Foo interface {
	Bar() error
	Baz(s, z string) (string, error)
	FooBarBaz(a any) (int, float64)
}

type Store[T, R any] interface {
	Do(obj T) (error, R, int, float64, string)
	Create() Bar
	CreateWithPointer() *Bar
}

type Bar struct{}
