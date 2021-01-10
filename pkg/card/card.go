package card

type Card struct {
	Id       int64
	Issuer   string
	Balance  int64
	Currency string
	Number   string
	Icon     string
}

type Service struct {
	BankName     string
	IssuerNumber string
	Cards        []*Card
}

func (c *Card) Withdraw(amount int64) bool {
	if c.Balance < amount {
		return false
	}
	c.Balance -= amount
	return true
}

func (c *Card) AddMoney(amount int64) {
	c.Balance += amount
}

func NewService(bankName, issuerNumber string) *Service {
	return &Service{
		BankName:     bankName,
		IssuerNumber: issuerNumber,
		Cards:        []*Card{},
	}
}

func (s *Service) FindCard(number string) *Card {
	for _, c := range s.Cards {
		if c.Number == number {
			return c
		}
	}
	return nil
}

func (s *Service) Issue(issuer string, startBalance int64, currency, number, icon string) *Card {
	var id int64 = 1
	if len(s.Cards) > 0 {
		id = s.Cards[len(s.Cards)-1].Id + 1
	}

	card := &Card{
		Id:       id,
		Issuer:   issuer,
		Balance:  startBalance,
		Currency: currency,
		Number:   number,
		Icon:     icon,
	}
	s.Cards = append(s.Cards, card)

	return card
}
