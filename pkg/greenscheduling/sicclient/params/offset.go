package sicparams

import (
	"fmt"
)

// Offset represents pagination offset.
type Offset struct {
	Value int // The offset value for pagination
}

// NewOffset creates a new Offset instance.
func NewOffset(value int) *Offset {
	return &Offset{Value: value}
}

// GetValue returns the value representation of the Offset.
// For example: "offset=10"
func (o *Offset) GetValue() string {
	return fmt.Sprintf("%d", o.Value)
}
