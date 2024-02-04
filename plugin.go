package scaffold

import (
	"sync"
	"time"
)

type (
	// Plugin is a compile-time plugin.
	Plugin[T any] interface {
		// String returns a unique name.
		String() string
		// Mount mounts and/or extends an application.
		Mount(c Context[T]) error
	}
)

// newPlugin initializes a Plugin wrapper.
func newPlugin[T any](plg Plugin[T]) *plugin[T] {
	return &plugin[T]{Plugin: plg, dependencies: make(map[string]Dependency)}
}

type (
	// plugin is a proxy Plugin for Scaffold.
	plugin[T any] struct {
		Plugin[T]
		mu           sync.RWMutex
		runtime      time.Duration
		dependencies map[string]Dependency
	}
)

// stat safely returns Info.
func (plg *plugin[T]) stat() Info {
	plg.mu.RLock()
	defer plg.mu.RUnlock()
	info := Info{
		Name:    plg.String(),
		Runtime: plg.runtime,
	}
	for _, dep := range plg.dependencies {
		info.Dependencies = append(info.Dependencies, dep)
	}
	return info
}

// updateRuntime safely updates the load time.
func (plg *plugin[T]) updateRuntime(d time.Duration) {
	plg.mu.Lock()
	defer plg.mu.Unlock()
	plg.runtime = d
}

// dependsOn safely adds a Plugin dependency.
func (plg *plugin[T]) dependsOn(name string) {
	plg.mu.Lock()
	defer plg.mu.Unlock()
	dep := Dependency{
		To:   name,
		From: plg.String(),
	}
	plg.dependencies[dep.String()] = dep
}