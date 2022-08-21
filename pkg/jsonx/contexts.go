package jsonx

import "fmt"

type ContextReader interface {
	Get(varName string) any
}

type ContextWriter interface {
	Set(varName string, value any)
}

type Context interface {
	ContextReader
	ContextWriter
}

type SimpleContext struct {
	vars map[string]any
}

func NewContext() Context {
	return &SimpleContext{vars: make(map[string]any)}
}

func (s *SimpleContext) Get(key string) any {
	return s.vars[key]
}
func (s *SimpleContext) Set(key string, value any) {
	fmt.Println("set", key, value)
	s.vars[key] = value
}
