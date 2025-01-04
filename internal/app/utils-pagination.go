package app

type PaginationOptions struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
}

func (opt PaginationOptions) Offset() int32 {
	return (opt.Page - 1) * opt.PageSize
}

func (opt PaginationOptions) Limit() int32 {
	return opt.PageSize
}

type PaginationRestrictions struct {
	MaxPageSize  int32
	MaxPageCount int32
}

func (r PaginationRestrictions) Validate(page, pageSize int32) PaginationOptions {
	if page < 1 {
		page = 1
	} else if page > r.MaxPageCount {
		page = r.MaxPageCount
	}

	if pageSize < 1 {
		pageSize = 1
	} else if pageSize > r.MaxPageSize {
		pageSize = r.MaxPageSize
	}

	return PaginationOptions{Page: page, PageSize: pageSize}
}

var (
	DefaultPaginationRestrictions = PaginationRestrictions{
		MaxPageSize:  100,
		MaxPageCount: 1000,
	}
)

func IsCommonPageSize(pageSize int32) bool {
	return pageSize == 20 || pageSize == 25 || pageSize == 30 || pageSize == 40 || pageSize == 50 || pageSize == 75 || pageSize == 100
}
