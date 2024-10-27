package sicparams

import (
	"fmt"
)

// SortKey is a custom type for defining valid sort keys.
type SortKey string

// SortKeys defines the available sort keys.
const (
	SortByEntityID        SortKey = "entityId"
	SortByEntityMake      SortKey = "entityMake"
	SortByEntityModel     SortKey = "entityModel"
	SortByEntityType      SortKey = "entityType"
	SortByEntitySerialNum SortKey = "entitySerialNum"
	SortByEntityProductID SortKey = "entityProductId"
	SortByLocationName    SortKey = "locationName"
	SortByLocationID      SortKey = "locationId"
	SortByLocationCity    SortKey = "locationCity"
	SortByLocationState   SortKey = "locationState"
	SortByLocationCountry SortKey = "locationCountry"
	SortByName            SortKey = "name"
)

// SortOrder defines the available sort orders.
type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

// Sort represents a single sort condition.
type Sort struct {
	Key   SortKey
	Order SortOrder
}

// NewSort creates a new Sort instance.
func NewSort(key SortKey, order SortOrder) (*Sort, error) {
	if !isValidSortKey(key) {
		return nil, fmt.Errorf("invalid sort key: %s", key)
	}
	if !isValidSortOrder(order) {
		return nil, fmt.Errorf("invalid sort order: %s", order)
	}
	return &Sort{Key: key, Order: order}, nil
}

// isValidSortKey checks if the sort key is valid.
func isValidSortKey(key SortKey) bool {
	switch key {
	case SortByEntityID, SortByEntityMake, SortByEntityModel,
		SortByEntityType, SortByEntitySerialNum, SortByEntityProductID,
		SortByLocationName, SortByLocationID, SortByLocationCity,
		SortByLocationState, SortByLocationCountry, SortByName:
		return true
	}
	return false
}

// isValidSortOrder checks if the sort order is valid.
func isValidSortOrder(order SortOrder) bool {
	switch order {
	case Ascending, Descending:
		return true
	}
	return false
}

// GetValue returns the value representation of the Sort.
// For example: "entityId asc"
func (s *Sort) GetValue() string {
	return fmt.Sprintf("%s %s", s.Key, s.Order)
}
