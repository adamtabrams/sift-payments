package sift

type Summary struct {
	Income     int
	Expenses   int
	Total      int
	Categories map[string]int
}

func (s Summary) Add(r Record, category string) Summary {
	res := s

	if s.Categories == nil {
		res.Categories = make(map[string]int)
	}

	if r.Amount > 0 {
		res.Income += r.Amount
	} else {
		res.Expenses -= r.Amount
	}

	res.Categories[category] += r.Amount
	res.Total = res.Income - res.Expenses

	return res
}
