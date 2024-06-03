package apimaker

type Form interface {
	Bind(Model) error
}

type Filter interface {
	GetFilters() map[string]interface{}
}
