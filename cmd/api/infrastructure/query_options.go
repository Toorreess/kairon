package infrastructure

type QueryOpts struct {
	QueryString string
	QueryParams map[string]any
	Offset      int
	Limit       int
	OrderBy     string
	Order       string
	RangeBy     string
	RangeSlice  []any
}

func NewDefaultQueryOpts() QueryOpts {
	return QueryOpts{
		Offset: 0,
		Limit:  20,
	}
}
