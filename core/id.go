package core

import "github.com/oklog/ulid/v2"

// NewID generates a new Unique ID.
//
// Returns a new Unique ID.
func NewID() string {
	return ulid.Make().String()
}
