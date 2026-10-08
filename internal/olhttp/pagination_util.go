package olhttp

// paginationPages uses zero for an ellipsis and keeps the first/last pages visible.
func paginationPages(page, total, size uint32) []uint32 {
	if total == 0 {
		return nil
	}
	page = min(max(page, 1), total)
	slots := min(max(size, 5), 7)
	if total <= slots {
		pages := make([]uint32, 0, total)
		for p := uint32(1); p <= total; p++ {
			pages = append(pages, p)
		}
		return pages
	}
	radius := (slots - 4) / 2
	start := max(uint32(2), page-min(page-1, radius))
	end := min(total-1, page+min(total-page, radius))
	if page <= slots-3 {
		start = 2
		end = slots - 2
	}
	if page >= total-(slots-4) {
		start = total - (slots - 3)
		end = total - 1
	}
	pages := []uint32{1}
	if start == 3 {
		pages = append(pages, 2)
	} else if start > 3 {
		pages = append(pages, 0)
	}
	for p := start; p <= end; p++ {
		pages = append(pages, p)
	}
	if end == total-2 {
		pages = append(pages, total-1)
	} else if end < total-2 {
		pages = append(pages, 0)
	}
	return append(pages, total)
}
