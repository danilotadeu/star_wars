package generic

type Pagination struct {
	Count        int64  `json:"count"`
	NextPage     *int64 `json:"next_page"`
	PreviousPage *int64 `json:"previous_page"`
}

func MakePagination(page int64) (*int64, *int64) {
	nPage := page + 1
	nextPage := &nPage

	nPrevious := page - 1
	previousPage := &nPrevious

	if *previousPage < 0 {
		previousPage = nil
	}

	return nextPage, previousPage
}
