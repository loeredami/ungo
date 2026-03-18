package ungo

import "fmt"

type Object struct {
	Data      FastMap[string, any]
	Methods   FastMap[string, func(*Object, ...any) Result[any]]
	Prototype *Object
}

func NewObject() *Object {
	return &Object{
		Data:      FastMap[string, any]{},
		Methods:   FastMap[string, func(*Object, ...any) Result[any]]{},
		Prototype: nil,
	}
}

func (o *Object) Get(key string) Result[any] {
	if value, ok := o.Data.Get(key); ok {
		return Result[any]{value: value}
	}
	if o.Prototype != nil {
		return o.Prototype.Get(key)
	}
	return Result[any]{err: fmt.Errorf("key not found: %s", key)}
}

func (o *Object) Set(key string, value any) {
	o.Data.Set(key, value)
}

func (o *Object) Call(method string, args ...any) Result[any] {
	if fn, ok := o.Methods.Get(method); ok {
		return fn(o, args...)
	}
	if o.Prototype != nil {
		return o.Prototype.Call(method, args...)
	}
	return Result[any]{err: fmt.Errorf("method not found: %s", method)}
}

func (o *Object) SetMethod(method string, fn func(*Object, ...any) Result[any]) {
	o.Methods.Set(method, fn)
}

func (o *Object) GetMethod(method string) (func(*Object, ...any) Result[any], bool) {
	if fn, ok := o.Methods.Get(method); ok {
		return fn, true
	}
	if o.Prototype != nil {
		return o.Prototype.GetMethod(method)
	}
	return nil, false
}

func (o *Object) SetPrototype(prototype *Object) {
	o.Prototype = prototype
}

func (o *Object) GetPrototype() *Object {
	return o.Prototype
}

func (o *Object) Has(key string) bool {
	return o.Data.Contains(key)
}

func FromPrototype(prototype *Object) *Object {
	return &Object{
		Data:      FastMap[string, any]{},
		Methods:   FastMap[string, func(*Object, ...any) Result[any]]{},
		Prototype: prototype,
	}
}
