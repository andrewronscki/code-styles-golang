package utils

import "strconv"

func Page(page string) int64 {
	if page == "" {
		return 1
	}
	p, err := strconv.ParseInt(page, 10, 64)

	if err != nil {
		return 1
	}

	return p
}

func PageSize(size string) int64 {
	if size == "" {
		return 10
	}

	s, err := strconv.ParseInt(size, 10, 64)

	if err != nil {
		return 10
	}

	return s
}

func PreviousPage(page int64) int64 {
	if page <= 1 {
		return 1
	}

	return page - 1
}

func NextPage(page, totalPages int64) int64 {
	next := page + 1
	if next > totalPages && totalPages != 0 {
		return totalPages
	}

	return page + 1
}
