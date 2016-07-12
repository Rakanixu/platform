package scan

type Info struct {
	Id          int64
	Type        string
	Description string
	Config      map[string]string
}

type Scanner interface {
	Start()
	Stop()
	Info() (Info, error)
}
