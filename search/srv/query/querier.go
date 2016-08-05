package query

type Querier interface {
	Query() (string, error)
}
