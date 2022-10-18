package blacklist

type BlackLister interface {
	Validate(data RowData) error
}

type CheckData interface {
}

type RowData struct {
	CheckData
}
