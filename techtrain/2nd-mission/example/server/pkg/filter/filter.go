// Package filter implements decorator handlers
package filter

import (
	"net/http"
)

// Decorator is the func type for all decorator handlers
type Decorator func(http.Handler) http.Handler

// Wrap invokes all decorators in order
// e.g. (h, decoA, decoB, decoC),
// running order would be decoA, decoB, decoC
func Wrap(h http.Handler, decorators ...Decorator) http.Handler {
	for i := range decorators {
		// reverse the order because it is the outer-most wrapper
		h = decorators[len(decorators)-1-i](h)
	}
	return h
}
