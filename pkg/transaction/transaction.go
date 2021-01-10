package transaction

type Transaction struct {
	Id     int64
	From   string
	To     string
	Amount int64
	Total  int64
}

type Service []*Transaction

func NewService() *Service {
	return &Service{}
}

func (s *Service) Add(from, to string, amount, total int64) *Transaction {
	var id int64 = 1
	if len(*s) > 0 {
		id = (*s)[len(*s)-1].Id + 1
	}
	transaction := &Transaction{
		Id:     id,
		From:   from,
		To:     to,
		Amount: amount,
		Total:  total,
	}
	*s = append(*s, transaction)
	return transaction
}
