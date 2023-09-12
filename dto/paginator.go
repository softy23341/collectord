package dto

// NewDefaultPagePaginator TBD
func NewDefaultPagePaginator() *PagePaginator {
	return &PagePaginator{
		Page: 0,
		Cnt:  16,
	}
}

// NewPagePaginator TBD
func NewPagePaginator(page, cnt int16) *PagePaginator {
	return &PagePaginator{
		Page: page,
		Cnt:  cnt,
	}
}

// PagePaginator TBD
type PagePaginator struct {
	Page int16
	Cnt  int16
}

// RangePaginator TBD
type RangePaginator struct {
	ID       *int64
	Distance int32
	Include  bool
}

// Offset TBD
func (p *PagePaginator) Offset() int64 {
	return int64(p.Page) * int64(p.Cnt)
}

// Limit TBD
func (p *PagePaginator) Limit() int16 {
	return p.Cnt
}

// NextPage RBD
func (p *PagePaginator) NextPage() {
	p.Page++
}
