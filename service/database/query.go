package database

type QueryParameters struct {
	Offset int
	Limit  int
	Filter string
}

func NewQueryParameters(page int, size int, filter string) QueryParameters {
	return QueryParameters{
		(page - 1) * size,
		size,
		filter,
	}
}
