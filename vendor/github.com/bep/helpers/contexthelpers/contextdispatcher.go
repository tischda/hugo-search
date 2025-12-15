// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package contexthelpers

import "context"

// ContextDispatcher is a generic interface for setting and getting values from context.Context.
type ContextDispatcher[T any] interface {
	// Set stores the value in the context and returns a new context.
	Set(ctx context.Context, value T) context.Context

	// Get retrieves the value from the context. If the value is not set, it returns the zero value of T.
	Get(ctx context.Context) T

	// Lookup retrieves the value from the context and returns a boolean indicating if the value was found.
	Lookup(ctx context.Context) (T, bool)
}

// NewContextDispatcher creates a new ContextDispatcher with the given key.
func NewContextDispatcher[T any, R comparable](key R) ContextDispatcher[T] {
	return keyInContext[T, R]{
		id: key,
	}
}

type keyInContext[T any, R comparable] struct {
	zero T
	id   R
}

func (f keyInContext[T, R]) Get(ctx context.Context) T {
	v := ctx.Value(f.id)
	if v == nil {
		return f.zero
	}
	return v.(T)
}

func (f keyInContext[T, R]) Lookup(ctx context.Context) (T, bool) {
	v := ctx.Value(f.id)
	if v == nil {
		return f.zero, false
	}
	return v.(T), true
}

func (f keyInContext[T, R]) Set(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, f.id, value)
}
