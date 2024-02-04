package scaffold

import "github.com/pedregon/scaffold/internal/list"

// Tag creates a new TaggedHook with the provided main Hook and optional tags.
func Tag[T Tagger](hook *Hook[T], tags ...string) *TaggedHook[T] {
	return &TaggedHook[T]{
		mainHook[T]{Hook: hook},
		tags,
	}
}

type (
	// Tagger defines an interface for event data structs that support tags/groups/categories/etc.
	//
	// Usually used together with TaggedHook.
	Tagger interface {
		Tags() []string
	}
	// mainHook is a private Hook.
	mainHook[T Tagger] struct {
		*Hook[T]
	}
	// TaggedHook is a proxy Hook for which register Handler(s) are conditionally triggered based on tag.
	//
	// A Hook is registered only if the tags are empty or includes at least one of the event data tag(s).
	TaggedHook[T Tagger] struct {
		mainHook[T]
		tags []string
	}
)

// CanTriggerOn checks if the current TaggedHook can be triggered with the provided event data tags.
func (h *TaggedHook[T]) CanTriggerOn(tags []string) bool {
	if len(h.tags) == 0 {
		return true // match all
	}

	for _, t := range tags {
		if list.ExistInSlice(t, h.tags) {
			return true
		}
	}

	return false
}

// PreAdd registers a new Handler to the Hook by prepending it to the existing queue.
//
// The Handler will be called only if the event data tags satisfy TaggedHook.CanTriggerOn.
func (h *TaggedHook[T]) PreAdd(fn Handler[T]) string {
	return h.mainHook.PreAdd(func(e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(e)
		}
		return nil
	})
}

// Add registers a new Handler to the Hook by appending it to the existing queue.
//
// The Handler will be called only if the event data tags satisfy TaggedHook.CanTriggerOn.
func (h *TaggedHook[T]) Add(fn Handler[T]) string {
	return h.mainHook.Add(func(e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(e)
		}
		return nil
	})
}