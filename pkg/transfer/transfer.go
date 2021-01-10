package transfer

import (
	"errors"
	"github.com/artrey/bgo-adv-errors/pkg/card"
	"github.com/artrey/bgo-adv-errors/pkg/transaction"
	"strconv"
	"strings"
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

var (
	NonPositiveAmount = errors.New("attempt to transfer negative or zero sum")
	NotEnoughMoney    = errors.New("not enough money on card to transfer")
	CardNotFound      = errors.New("card not found in bank")
	InvalidCardNumber = errors.New("card number is invalid")
)

func checkOwning(cardNumber, issuerNumber string) bool {
	return strings.HasPrefix(cardNumber, issuerNumber)
}

func IsValidCardNumber(number string) bool {
	numberStr := strings.Split(strings.ReplaceAll(number, " ", ""), "")
	if len(numberStr) == 0 {
		return false
	}

	numberDigits := make([]int, 0, len(number))
	for _, val := range numberStr {
		digit, err := strconv.Atoi(val)
		if err != nil {
			return false
		}
		numberDigits = append(numberDigits, digit)
	}

	for ri := len(numberDigits) - 2; ri >= 0; ri -= 2 {
		val := numberDigits[ri] * 2
		if val > 9 {
			val -= 9
		}
		numberDigits[ri] = val
	}

	result := 0
	for _, val := range numberDigits {
		result += val
	}
	return result % 10 == 0
}

func (s *Service) Card2Card(from, to string, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, NonPositiveAmount
	}

	if !IsValidCardNumber(from) || !IsValidCardNumber(to) {
		return 0, InvalidCardNumber
	}

	fromCard := s.CardSvc.FindCard(from)
	if fromCard == nil && checkOwning(from, s.CardSvc.IssuerNumber) {
		return 0, CardNotFound
	}

	toCard := s.CardSvc.FindCard(to)
	if toCard == nil && checkOwning(to, s.CardSvc.IssuerNumber) {
		return 0, CardNotFound
	}

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
	total := amount + commission

	if fromCard != nil && !fromCard.Withdraw(total) {
		return total, NotEnoughMoney
	}

	if toCard != nil {
		toCard.AddMoney(amount)
	}

	s.TransactionSvc.Add(from, to, amount, total)

	return total, nil
}
