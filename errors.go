package scaffold

import "errors"

var (
	ErrPluginNotRegistered       error = errors.New("plugin not registered")
	ErrSelfReferentialDependency error = errors.New("self-referential plugin dependency")
	ErrCircularDependency        error = errors.New("circular plugin dependency")
	ErrMissingDependency         error = errors.New("missing plugin dependency")
	ErrStopPropagation           error = errors.New("event hook propagation stopped")
)
