package transfer

import (
	"github.com/artrey/bgo-adv-errors/pkg/card"
	"github.com/artrey/bgo-adv-errors/pkg/transaction"
	"math"
	"testing"
)

func TestService_Card2Card(t *testing.T) {
	type fields struct {
		CardSvc        *card.Service
		TransactionSvc *transaction.Service
		commissions    Commissions
	}
	type args struct {
		from   string
		to     string
		amount int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantTotal int64
		wantOk    bool
	}{
		{
			name: "Inner success",
			fields: fields{
				CardSvc: &card.Service{
					BankName: "Tinkoff",
					Cards: []*card.Card{
						{
							Id:       1,
							Issuer:   "Visa",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0001",
							Icon:     "...",
						},
						{
							Id:       2,
							Issuer:   "MasterCard",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0002",
							Icon:     "...",
						},
					},
				},
				TransactionSvc: transaction.NewService(),
				commissions: Commissions{
					FromInner: func(val int64) int64 {
						return int64(math.Max(float64(val*5/1000), 10_00))
					},
					ToInner: func(val int64) int64 {
						return 0
					},
					FromOuterToOuter: func(val int64) int64 {
						return int64(math.Max(float64(val*15/1000), 30_00))
					},
				},
			},
			args: args{
				from:   "0001",
				to:     "0002",
				amount: 500_00,
			},
			wantTotal: 510_00,
			wantOk:    true,
		},
		{
			name: "Inner fail",
			fields: fields{
				CardSvc: &card.Service{
					BankName: "Tinkoff",
					Cards: []*card.Card{
						{
							Id:       1,
							Issuer:   "Visa",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0001",
							Icon:     "...",
						},
						{
							Id:       2,
							Issuer:   "MasterCard",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0002",
							Icon:     "...",
						},
					},
				},
				TransactionSvc: transaction.NewService(),
				commissions: Commissions{
					FromInner: func(val int64) int64 {
						return int64(math.Max(float64(val*5/1000), 10_00))
					},
					ToInner: func(val int64) int64 {
						return 0
					},
					FromOuterToOuter: func(val int64) int64 {
						return int64(math.Max(float64(val*15/1000), 30_00))
					},
				},
			},
			args: args{
				from:   "0001",
				to:     "0002",
				amount: 1000_00,
			},
			wantTotal: 1010_00,
			wantOk:    false,
		},
		{
			name: "Inner-outer success",
			fields: fields{
				CardSvc: &card.Service{
					BankName: "Tinkoff",
					Cards: []*card.Card{
						{
							Id:       1,
							Issuer:   "Visa",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0001",
							Icon:     "...",
						},
					},
				},
				TransactionSvc: transaction.NewService(),
				commissions: Commissions{
					FromInner: func(val int64) int64 {
						return int64(math.Max(float64(val*5/1000), 10_00))
					},
					ToInner: func(val int64) int64 {
						return 0
					},
					FromOuterToOuter: func(val int64) int64 {
						return int64(math.Max(float64(val*15/1000), 30_00))
					},
				},
			},
			args: args{
				from:   "0001",
				to:     "0002",
				amount: 500_00,
			},
			wantTotal: 510_00,
			wantOk:    true,
		},
		{
			name: "Inner-outer fail",
			fields: fields{
				CardSvc: &card.Service{
					BankName: "Tinkoff",
					Cards: []*card.Card{
						{
							Id:       1,
							Issuer:   "Visa",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0001",
							Icon:     "...",
						},
					},
				},
				TransactionSvc: transaction.NewService(),
				commissions: Commissions{
					FromInner: func(val int64) int64 {
						return int64(math.Max(float64(val*5/1000), 10_00))
					},
					ToInner: func(val int64) int64 {
						return 0
					},
					FromOuterToOuter: func(val int64) int64 {
						return int64(math.Max(float64(val*15/1000), 30_00))
					},
				},
			},
			args: args{
				from:   "0001",
				to:     "0002",
				amount: 1000_00,
			},
			wantTotal: 1010_00,
			wantOk:    false,
		},
		{
			name: "Outer-inner success",
			fields: fields{
				CardSvc: &card.Service{
					BankName: "Tinkoff",
					Cards: []*card.Card{
						{
							Id:       1,
							Issuer:   "Visa",
							Balance:  1000_00,
							Currency: "RUB",
							Number:   "0001",
							Icon:     "...",
						},
					},
				},
				TransactionSvc: transaction.NewService(),
				commissions: Commissions{
					FromInner: func(val int64) int64 {
						return int64(math.Max(float64(val*5/1000), 10_00))
					},
					ToInner: func(val int64) int64 {
						return 0
					},
					FromOuterToOuter: func(val int64) int64 {
						return int64(math.Max(float64(val*15/1000), 30_00))
					},
				},
			},
			args: args{
				from:   "0002",
				to:     "0001",
				amount: 1000_00,
			},
			wantTotal: 1000_00,
			wantOk:    true,
		},
		{
			name: "Outer success",
			fields: fields{
				CardSvc: &card.Service{
					BankName: "Tinkoff",
					Cards:    []*card.Card{},
				},
				TransactionSvc: transaction.NewService(),
				commissions: Commissions{
					FromInner: func(val int64) int64 {
						return int64(math.Max(float64(val*5/1000), 10_00))
					},
					ToInner: func(val int64) int64 {
						return 0
					},
					FromOuterToOuter: func(val int64) int64 {
						return int64(math.Max(float64(val*15/1000), 30_00))
					},
				},
			},
			args: args{
				from:   "0002",
				to:     "0001",
				amount: 1000_00,
			},
			wantTotal: 1030_00,
			wantOk:    true,
		},
	}
	for _, tt := range tests {
		s := &Service{
			CardSvc:        tt.fields.CardSvc,
			TransactionSvc: tt.fields.TransactionSvc,
			commissions:    tt.fields.commissions,
		}
		gotTotal, gotOk := s.Card2Card(tt.args.from, tt.args.to, tt.args.amount)
		if gotTotal != tt.wantTotal {
			t.Errorf("Card2Card() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
		}
		if gotOk != tt.wantOk {
			t.Errorf("Card2Card() gotOk = %v, want %v", gotOk, tt.wantOk)
		}
	}
}
