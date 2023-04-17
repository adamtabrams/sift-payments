package pay

type PaymentListMap map[string]PaymentList

type PaymentList []Payment

type Payment struct {
	Name     string
	Amount   int
	Category string
}

type Summary struct {
	Income     int
	Expenses   int
	Total      int
	Categories map[string]int
}

func (plm PaymentListMap) SummaryMap() map[string]Summary {
	summaryMap := make(map[string]Summary, len(plm))

	for name, paymentList := range plm {
		summaryMap[name] = paymentList.Summary()
	}

	return summaryMap
}

func (pl PaymentList) Summary() Summary {
	s := Summary{
		Categories: make(map[string]int),
	}

	for _, payment := range pl {
		if payment.Amount > 0 {
			s.Income += payment.Amount
		} else {
			s.Expenses -= payment.Amount
		}

		s.Categories[payment.Category] += payment.Amount
	}

	s.Total = s.Income - s.Expenses

	return s
}
