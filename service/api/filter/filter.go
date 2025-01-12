package filter

type Filter interface {
	Evaluate(formula string) (string, error)
}
