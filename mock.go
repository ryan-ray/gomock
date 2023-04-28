package main

// IMPORTANT!
// Do not edit or delete any code in this file unless you really, truly,
// madly, deeply know what you are doing. Depending on your setup, any changes
// you make here will most likely be overwritten when these mocks are
// regenerated.
type FooMock struct {
	BarFunc func() error
	BazFunc func(s string, z string) (string, error)
	FooBarBazFunc func(a any) (int, float64)
}

// NewFooMock creates a new instance of FooMock.
// This instance implements the Foo interface.
func NewFooMock() *FooMock {
	return &FooMock{
		BarFunc: func() error { return errors.New("not implemented") },
		BazFunc: func(s string, z string) (string, error) { return "", errors.New("not implemented") },
		FooBarBazFunc: func(a any) (int, float64) { return 0, 0 },
		
	}
}

// Bar function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the FooMock.BarFunc variable as functionality requires.
func (m FooMock) Bar() error {
	return m.BarFunc()
}

// Baz function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the FooMock.BazFunc variable as functionality requires.
func (m FooMock) Baz(s string, z string) (string, error) {
	return m.BazFunc(s, z)
}

// FooBarBaz function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the FooMock.FooBarBazFunc variable as functionality requires.
func (m FooMock) FooBarBaz(a any) (int, float64) {
	return m.FooBarBazFunc(a)
}

type StoreMock[T any, R any] struct {
	DoFunc func(obj T) (error, R, int, float64, string)
	CreateFunc func() Bar
	CreateWithPointerFunc func() *Bar
}

// NewStoreMock creates a new instance of StoreMock.
// This instance implements the Store interface.
func NewStoreMock[T any, R any]() *StoreMock[T, R] {
	return &StoreMock[T, R]{
		DoFunc: func(obj T) (error, R, int, float64, string) { return errors.New("not implemented"), *new(R), 0, 0, "" },
		CreateFunc: func() Bar { return *new(Bar) },
		CreateWithPointerFunc: func() *Bar { return new(Bar) },
		
	}
}

// Do function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the StoreMock.DoFunc variable as functionality requires.
func (m StoreMock[T, R]) Do(obj T) (error, R, int, float64, string) {
	return m.DoFunc(obj)
}

// Create function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the StoreMock.CreateFunc variable as functionality requires.
func (m StoreMock[T, R]) Create() Bar {
	return m.CreateFunc()
}

// CreateWithPointer function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the StoreMock.CreateWithPointerFunc variable as functionality requires.
func (m StoreMock[T, R]) CreateWithPointer() *Bar {
	return m.CreateWithPointerFunc()
}

