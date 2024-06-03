package apimaker

type Model interface {
	Save() error
	GetOne(id interface{}) error
	//filter for filter, page filtering - totalCounts, totalPages, list , error
	List(filter Filter, pfilter Pagination) (int, int, interface{}, error)
	Remove(id interface{}) error
}

type (
	AfterSave  func(model Model) error
	BeforeSave func(model Model) error
)
