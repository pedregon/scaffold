package scaffold

import (
	"context"
	"fmt"
	"github.com/pedregon/scaffold/internal/stack"
	"time"
)

var (
	// interface guard context.Context.
	_ context.Context = Context[any]{}
)

type (
	// Context is a transient context for loading Plugin(s) registered in Scaffold.
	Context[T any] struct {
		context.Context
		loaders  []Loader
		stack    *stack.Stack[string]
		scaffold *Scaffold[T]
	}
)

// newContext initializes a Context with a variadic of Loader(s).
func newContext[T any](ctx context.Context, scaffold *Scaffold[T], loaders ...Loader) Context[T] {
	return Context[T]{
		Context:  ctx,
		loaders:  loaders,
		stack:    stack.New[string](),
		scaffold: scaffold,
	}
}

// Deadline implements context.Context.
func (c Context[T]) Deadline() (time.Time, bool) {
	return c.Context.Deadline()
}

// Done implements context.Context.
func (c Context[T]) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err implements context.Context.
func (c Context[T]) Err() error {
	if c.stack.Err() != nil {
		return c.stack.Err()
	}
	return c.Context.Err()
}

// Value implements context.Context.
func (c Context[T]) Value(key any) any {
	return c.Context.Value(key)
}

// Set is equivalent to context.WithValue.
func (c Context[T]) Set(key, val any) {
	c.Context = context.WithValue(c.Context, key, val)
}

// Lazy lazy loads a Plugin by name. Use this for eagerly loading Plugin dependencies.
func (c Context[T]) Lazy(name string) error {
	if len(c.loaders) == 0 {
		return c.load(name)
	}
	for _, loader := range c.loaders {
		if err := loader(name, c.load); err != nil {
			return err
		}
	}
	return nil
}

func (c Context[T]) load(name string) error {
	if err := c.Err(); err != nil {
		return err
	}
	plg, exist := c.scaffold.lookup(name)
	if !exist {
		err := ErrPluginNotRegistered
		if c.stack.Size() > 0 {
			err = ErrMissingDependency
		}
		err = fmt.Errorf("%w, %s", err, name)
		return err
	}
	if plg.stat().Runtime > 0 {
		return nil
	}
	if current, ok := c.stack.Peek(); ok && current == name {
		err := fmt.Errorf("%w, %s", ErrSelfReferentialDependency, name)
		return err
	}
	if c.stack.Has(name) {
		err := fmt.Errorf("%w, %s", ErrCircularDependency, name)
		return err
	}
	c.stack.Push(name)
	index := c.stack.Size() - 1
	start := time.Now()
	if err := plg.Mount(c); err != nil {
		err = fmt.Errorf("%w, %s", err, name)
		c.stack.Catch(err)
		return err
	}
	plg.updateRuntime(time.Since(start))
	for {
		if err := c.stack.Err(); err != nil {
			err = fmt.Errorf("%w, %s", err, name)
			return err
		}
		if c.stack.Size()-1 == index {
			break
		}
		if top, ok := c.stack.Pop(); ok {
			plg.dependsOn(top)
		}
	}
	return nil
}
