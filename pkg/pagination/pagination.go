// Package pagination contains common cursor pagination parameters.
package pagination

import (
	"net/url"
	"strconv"
)

// Params are common cursor pagination parameters.
type Params struct {
	After  string
	Before string
	Limit  int
	Order  string
}

// Values encodes non-zero parameters as query values.
func (p Params) Values() url.Values {
	values := make(url.Values)
	if p.After != "" {
		values.Set("after", p.After)
	}
	if p.Before != "" {
		values.Set("before", p.Before)
	}
	if p.Limit > 0 {
		values.Set("limit", strconvItoa(p.Limit))
	}
	if p.Order != "" {
		values.Set("order", p.Order)
	}
	return values
}

func strconvItoa(value int) string {
	return strconv.Itoa(value)
}
