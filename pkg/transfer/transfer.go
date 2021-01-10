package transfer

import (
	"github.com/artrey/bgo-adv-errors/pkg/card"
	"github.com/artrey/bgo-adv-errors/pkg/transaction"
)

type CommissionEvaluator func(val int64) int64

type Commissions struct {
	FromInner        CommissionEvaluator
	ToInner          CommissionEvaluator
	FromOuterToOuter CommissionEvaluator
}

type Service struct {
	CardSvc        *card.Service
	TransactionSvc *transaction.Service
	commissions    Commissions
}

func NewService(cardSvc *card.Service, transactionSvc *transaction.Service, commissions Commissions) *Service {
	return &Service{
		CardSvc:        cardSvc,
		TransactionSvc: transactionSvc,
		commissions:    commissions,
	}
}

func (s *Service) Card2Card(from, to string, amount int64) (total int64, ok bool) {
	fromCard := s.CardSvc.FindCard(from)
	toCard := s.CardSvc.FindCard(to)

	var commission int64 = 0
	if fromCard == nil && toCard == nil {
		commission += s.commissions.FromOuterToOuter(amount)
	} else {
		if toCard != nil {
			commission += s.commissions.ToInner(amount)
		}
		if fromCard != nil {
			commission += s.commissions.FromInner(amount)
		}
	}
	total = amount + commission

	ok = true
	if fromCard != nil {
		ok = fromCard.Withdraw(total)
	}

	if ok && toCard != nil {
		toCard.AddMoney(amount)
	}

	if ok {
		s.TransactionSvc.Add(from, to, amount, total)
	}

	return
}
