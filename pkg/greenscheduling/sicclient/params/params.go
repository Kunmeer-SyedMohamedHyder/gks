package sicparams

import "net/url"

// Param is an interface that defines methods for adding and retrieving parameters.
type Param interface {
	GetValue() string
}

// Params is a struct that implements the Param interface to manage all parameters.
type Params struct {
	filters []Filter
	sorts   []Sort
	offset  *Offset
	limit   *Limit
}

// NewParams creates a new Params instance.
func NewParams() *Params {
	return &Params{
		filters: []Filter{},
		sorts:   []Sort{},
		offset:  nil,
		limit:   nil,
	}
}

// AddFilter adds a filter to the Params.
func (p *Params) AddFilter(filter *Filter) {
	p.filters = append(p.filters, *filter)
}

// AddSort adds a sort condition to the Params.
func (p *Params) AddSort(sort *Sort) {
	p.sorts = append(p.sorts, *sort)
}

// AddOffset sets the offset in Params.
func (p *Params) AddOffset(offset *Offset) {
	p.offset = offset
}

// AddLimit sets the limit in Params.
func (p *Params) AddLimit(limit *Limit) {
	p.limit = limit
}

// ToQueryParams converts filters and sorts into url.Values for API call.
func (p *Params) ToQueryParams() url.Values {
	queryParams := url.Values{}

	// Add filters to query parameters
	for _, filter := range p.filters {
		queryParams.Add("filter", filter.GetValue())
	}

	// Add sorts to query parameters
	for _, sort := range p.sorts {
		queryParams.Add("sort", sort.GetValue())
	}

	// Add offset to query parameters
	if p.offset != nil {
		queryParams.Add("offset", p.offset.GetValue())
	}

	// Add limit to query parameters
	if p.limit != nil {
		queryParams.Add("limit", p.limit.GetValue())
	}

	return queryParams
}
