package main

// IMPORTANT!
// Do not edit or delete any code in this file unless you really, truly,
// madly, deeply know what you are doing. Depending on your setup, any changes
// you make here will most likely be overwritten when these mocks are
// regenerated.
type FooMock struct {
	barFunc func() error
	bazFunc func(s string, z string) (string, error)
	fooBarBazFunc func(a any) (int, float64)
}

// NewFooMock creates a new instance of FooMock.
// This instance implements the Foo interface.
// 
// All functions are autogenerated and will return default or nil values where
// applicable. Do not edit this implementation directly unless you know what
// you are doing.
func NewFooMock[]() *FooMock[] {
	return &FooMock{
		barFunc: func() error { return errors.New("not implemented") },
		bazFunc: func(s string, z string) (string, error) { return "", errors.New("not implemented") },
		fooBarBazFunc: func(a any) (int, float64) { return 0, 0 },
		
	}
}

// Bar function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the FooMock.barFunc variable as functionality requires.
func (m FooMock) Bar() error {
	return m.barFunc()
}

// Baz function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the FooMock.bazFunc variable as functionality requires.
func (m FooMock) Baz(s string, z string) (string, error) {
	return m.bazFunc(s, z)
}

// FooBarBaz function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the FooMock.fooBarBazFunc variable as functionality requires.
func (m FooMock) FooBarBaz(a any) (int, float64) {
	return m.fooBarBazFunc(a)
}

type StoreMock[T any, R any] struct {
	doFunc func(obj T) R
	createFunc func() Bar
	createWithPointerFunc func() *Bar
}

// NewStoreMock creates a new instance of StoreMock.
// This instance implements the Store interface.
// 
// All functions are autogenerated and will return default or nil values where
// applicable. Do not edit this implementation directly unless you know what
// you are doing.
func NewStoreMock[T any, R any]() *StoreMock[T, R] {
	return &StoreMock[T, R]{
		doFunc: func(obj T) R { return *new(R) },
		createFunc: func() Bar { return *new(Bar) },
		createWithPointerFunc: func() *Bar { return new(Bar) },
		
	}
}

// Do function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the StoreMock.doFunc variable as functionality requires.
func (m StoreMock[T, R]) Do(obj T) R {
	return m.doFunc(obj)
}

// Create function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the StoreMock.createFunc variable as functionality requires.
func (m StoreMock[T, R]) Create() Bar {
	return m.createFunc()
}

// CreateWithPointer function mock implementation.
// 
// This function is auto generated. Do not edit this directly.
// Override the StoreMock.createWithPointerFunc variable as functionality requires.
func (m StoreMock[T, R]) CreateWithPointer() *Bar {
	return m.createWithPointerFunc()
}
