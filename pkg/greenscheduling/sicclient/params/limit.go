package sicparams

import (
	"fmt"
)

// Limit represents the maximum number of results to return.
type Limit struct {
	Value int // The limit value for pagination
}

// NewLimit creates a new Limit instance.
func NewLimit(value int) *Limit {
	return &Limit{Value: value}
}

// GetValue returns the value representation of the Limit.
// For example: "limit=10"
func (l *Limit) GetValue() string {
	return fmt.Sprintf("%d", l.Value)
}
