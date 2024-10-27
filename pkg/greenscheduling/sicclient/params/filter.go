package sicparams

import (
	"fmt"
	"strings"
)

// FilterKey is a custom type for defining valid filter keys.
type FilterKey string

// FilterKeys defines the available filter keys.
const (
	EntityID        FilterKey = "entityId"
	EntityMake      FilterKey = "entityMake"
	EntityModel     FilterKey = "entityModel"
	EntityType      FilterKey = "entityType"
	EntitySerialNum FilterKey = "entitySerialNum"
	EntityProductID FilterKey = "entityProductId"
	LocationName    FilterKey = "locationName"
	LocationID      FilterKey = "locationId"
	LocationCity    FilterKey = "locationCity"
	LocationState   FilterKey = "locationState"
	LocationCountry FilterKey = "locationCountry"
	Name            FilterKey = "name"
)

// FilterOperator is a custom type for defining valid filter operators.
type FilterOperator string

// FilterOperators defines the available filter operators.
const (
	Equals   FilterOperator = "eq"
	Contains FilterOperator = "contains"
	In       FilterOperator = "in"
)

// Filter represents a single filter condition.
type Filter struct {
	Key      FilterKey
	Operator FilterOperator
	Value    interface{} // Change Value to interface{} to allow flexibility
}

// NewFilter creates a new Filter instance.
func NewFilter(key FilterKey, operator FilterOperator, value interface{}) (*Filter, error) {
	if !isValidKey(key) {
		return nil, fmt.Errorf("invalid filter key: %s", key)
	}
	if !isValidOperator(operator) {
		return nil, fmt.Errorf("invalid filter operator: %s", operator)
	}
	return &Filter{Key: key, Operator: operator, Value: value}, nil
}

// isValidKey checks if the filter key is valid.
func isValidKey(key FilterKey) bool {
	switch key {
	case EntityID, EntityMake, EntityModel, EntityType, EntitySerialNum,
		EntityProductID, LocationName, LocationID,
		LocationCity, LocationState, LocationCountry, Name:
		return true
	}
	return false
}

// isValidOperator checks if the filter operator is valid.
func isValidOperator(operator FilterOperator) bool {
	switch operator {
	case Equals, Contains, In:
		return true
	}
	return false
}

// GetValue returns the value representation of the Filter.
// For example: "entityId eq 'value'" or "contains(entityMake, 'value')"
func (f *Filter) GetValue() string {
	if f.Operator == In {
		// If the value is a slice, format it as a string
		values := f.Value.([]string) // Assuming Value is a slice of strings for "in" operator
		return fmt.Sprintf("%s in (%s)", f.Key, formatInValues(values))
	}
	if f.Operator == Contains {
		return fmt.Sprintf("contains(%s, '%v')", f.Key, f.Value)
	}
	return fmt.Sprintf("%s %s '%v'", f.Key, f.Operator, f.Value)
}

// formatInValues formats a slice of strings into the SQL-like syntax
func formatInValues(values []string) string {
	for i := range values {
		values[i] = "'" + values[i] + "'"
	}
	return strings.Join(values, ", ")
}
