package webcomponents

func getPaginationOrder(page, totalPages, size uint32) []uint32 {
	if totalPages == 0 {
		return nil
	}

	remaining := int(size) - 1
	left := remaining / 2

	if left > int(page)-1 {
		left = int(page) - 1
	}

	remaining -= left

	right := remaining
	if int(page)+right > int(totalPages) {
		right = int(totalPages) - int(page)
	}

	remaining -= right

	if remaining > 0 && left < int(page)-1 {
		need := int(page) - 1 - left
		if remaining < need {
			left += remaining
		} else {
			left += need
		}
	}

	var result []uint32

	// first page
	if page != 1 {
		result = append(result, 1)
	}

	// left side
	for i := int(page) - left + 1; i < int(page); i++ {
		result = append(result, uint32(i))
	}

	// current page
	result = append(result, page)

	// right side
	for i := int(page) + 1; i <= int(page)+right-1; i++ {
		result = append(result, uint32(i))
	}

	// last page
	if page < totalPages {
		result = append(result, totalPages)
	}

	return result
}
