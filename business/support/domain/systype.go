package domain

import (
	"github.com/shopspring/decimal"
)

var (
	MAX_RATE, _ = decimal.NewFromString("0.00375")
	MIN_RATE, _ = decimal.NewFromString("-0.00375")

	DB_PRECISION                int32 = 5
	DecimalMax                        = decimal.New(1, 19)
	DecimalSimilarZero                = decimal.New(1, -18)
	DecimalPFive                      = decimal.New(5, -1)
	Decimal6P                         = decimal.New(1, -6)
	Decimal8P                         = decimal.New(1, -8)
	Decimal10P                        = decimal.New(1, -10)
	Decimal12P                        = decimal.New(1, -12)
	DecimalOne                        = decimal.New(1, 0)
	DecimalTen                        = decimal.New(1, 1)
	DecimalMinusOne                   = decimal.New(-1, 0)
	DecimalTow                        = decimal.New(2, 0)
	DecimalHundredMillion             = decimal.New(1, 8)
	DecimalLow                        = decimal.New(1, -13)
	DecimalDBUnit                     = decimal.New(-1, -9)
	Decimal100                        = decimal.New(1, 2)
	DecimalDebtRatio                  = decimal.NewFromFloat(0.81)
	DecimalEductibleMarginRatio       = decimal.NewFromFloat(0.95)
)

func DecimalReduceDBPrecision(vol decimal.Decimal) decimal.Decimal {
	return vol.Sub(Decimal8P)
}

func DecimalABS(item decimal.Decimal) decimal.Decimal {
	if item.Sign() >= 0 {
		return item
	}
	return decimal.Zero.Sub(item)
}

func CheckVol(vol decimal.Decimal) bool {
	return vol.GreaterThan(DecimalSimilarZero) && vol.LessThan(DecimalMax)
}

func CheckPrice(price decimal.Decimal) bool {
	return price.GreaterThan(DecimalSimilarZero) && price.LessThan(DecimalMax)
}
