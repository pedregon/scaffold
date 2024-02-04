// Package scaffold is a compile-time plugin framework for building extendable applications. Scaffold encourages
// inversion of control via a generic specification. A Registrar may be initialized with a specification.
package scaffold

import (
	"context"
	"sync"
	"time"
)

type (
	// scaffoldold is a plugin manager.
	Scaffold[T any] struct {
		mu      sync.RWMutex
		plugins map[string]*plugin[T]
		app     T
	}
	// Dependency is a Plugin dependency relationship.
	Dependency struct {
		To   string
		From string
	}
	// Info is scaffoldold information about a Plugin.
	Info struct {
		Name         string
		Runtime      time.Duration
		Dependencies []Dependency
	}
)

// New initializes a scaffoldold for an application.
func New[T any](app T) *Scaffold[T] {
	return &Scaffold[T]{
		plugins: make(map[string]*plugin[T]),
		app:     app,
	}
}

// Register registers a Plugin.
func (scaffold *Scaffold[T]) Register(plg Plugin[T]) {
	if plg == nil {
		return
	}
	scaffold.mu.Lock()
	defer scaffold.mu.Unlock()
	scaffold.plugins[plg.String()] = newPlugin(plg)
}

// Load loads registered Plugin(s) for an application.
func (scaffold *Scaffold[T]) Load(ctx context.Context, app T, loaders ...Loader) error {
	c := newContext[T](ctx, scaffold, loaders...)
	scaffold.mu.RLock()
	defer scaffold.mu.RUnlock()
	for _, plg := range scaffold.plugins {
		if err := c.Lazy(plg.String()); err != nil {
			return err
		}
	}
	return nil
}

// Open retrieves Info about a plugin.
func (scaffold *Scaffold[T]) Open(name string) (Info, bool) {
	scaffold.mu.RLock()
	defer scaffold.mu.RUnlock()
	plg, ok := scaffold.plugins[name]
	if !ok {
		return Info{}, false
	}
	return plg.stat(), true
}

func (scaffold *Scaffold[T]) lookup(name string) (*plugin[T], bool) {
	scaffold.mu.RLock()
	defer scaffold.mu.RUnlock()
	plg, ok := scaffold.plugins[name]
	if !ok {
		return nil, false
	}
	return plg, true
}

// String implements fmt.String.
func (dep Dependency) String() string {
	return dep.From + "->" + dep.To
}

// Len returns the number of registered plugins in a Registrar.
func Len[T any](scaffold *Scaffold[T]) (i int) {
	scaffold.mu.RLock()
	defer scaffold.mu.RUnlock()
	for range scaffold.plugins {
		i++
	}
	return
}

// Graph returns a Scaffold dependency graph.
func Graph[T any](scaffold *Scaffold[T]) (deps []Dependency) {
	scaffold.mu.RLock()
	defer scaffold.mu.RUnlock()
	for _, plg := range scaffold.plugins {
		info := plg.stat()
		deps = append(deps, info.Dependencies...)
	}
	return
}
