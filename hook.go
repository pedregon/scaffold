package scaffold

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

const defaultCharacterSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// pseudorandomString generates a pseudorandom string with the specified length and character set.
func pseudorandomString(length int, characters string) string {
	b := make([]byte, length)
	max := len(characters)

	for i := range b {
		b[i] = characters[rand.Intn(max)]
	}

	return string(b)
}

// generateHookId generates a pseudorandom identifier.
func generateHookId() string {
	return pseudorandomString(8, defaultCharacterSet)
}

type (
	// Handler defines a Hook handler function.
	Handler[T any] func(e T) error
	// handlerPair defines a pair of string id and Handler.
	handlerPair[T any] struct {
		id      string
		handler Handler[T]
	}
	// Hook defines a concurrent safe structure for handling event hooks (aka. callbacks propagation).
	Hook[T any] struct {
		mu       sync.RWMutex
		handlers []*handlerPair[T]
	}
)

// PreAdd registers a new Handler to the Hook by prepending it to the existing queue.
func (h *Hook[T]) PreAdd(fn Handler[T]) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := generateHookId()
	// minimize allocations by shifting the slice
	h.handlers = append(h.handlers, nil)
	copy(h.handlers[1:], h.handlers)
	h.handlers[0] = &handlerPair[T]{id, fn}

	return id
}

// Add registers a new Handler to the Hook by appending it to the existing queue.
func (h *Hook[T]) Add(fn Handler[T]) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := generateHookId()
	h.handlers = append(h.handlers, &handlerPair[T]{id, fn})

	return id
}

// Remove removes a single Hook Handler by its id.
func (h *Hook[T]) Remove(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := len(h.handlers) - 1; i >= 0; i-- {
		if h.handlers[i].id == id {
			h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
			return
		}
	}
}

// Reset removes all registered Handler(s).
func (h *Hook[T]) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.handlers = nil
}

// Trigger executes all registered Hook Handlers one by one.
//
// Optionally, this method allows calling one off Handler(s) that will be
// temporarily appended to the Handler(s) queue.
//
// The execution stops when:
// - ErrStopPropagation is returned in one of the Handler(s)
// - any non-nil error is returned in one of the Handler(s)
//
// This method may be called at any time.
func (h *Hook[T]) Trigger(data T, oneOffHandlers ...Handler[T]) error {
	h.mu.RLock()
	handlers := make([]*handlerPair[T], 0, len(h.handlers)+len(oneOffHandlers))
	handlers = append(handlers, h.handlers...)
	// append the one off handlers
	for i, oneOff := range oneOffHandlers {
		handlers = append(handlers, &handlerPair[T]{
			id:      fmt.Sprintf("@%d", i),
			handler: oneOff,
		})
	}

	// unlock is not deferred to avoid deadlocks in case Trigger
	// is called recursively by the handlers
	h.mu.RUnlock()

	for _, item := range handlers {
		err := item.handler(data)
		if err == nil {
			continue
		}

		if errors.Is(err, ErrStopPropagation) {
			return nil
		}

		return err
	}

	return nil
}
