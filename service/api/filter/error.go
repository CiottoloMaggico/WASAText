package filter

type FilterError struct {
	Err    string
	Detail string
}

func (e *FilterError) Error() string {
	return e.Err
}

func NewFilterError(err string, detail string) *FilterError {
	return &FilterError{err, detail}
}
