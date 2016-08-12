package scan

type Info struct {
	Id          int64
	Type        string
	Description string
	Config      map[string]string
}

type Scanner interface {
	Start(map[int64]Scanner, int64)
	Stop()
	Info() (Info, error)
}
