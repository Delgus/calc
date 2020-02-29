package calc

import (
	"math"
	"math/big"
)

const (
	// MoneyPrec - precition for money
	MoneyPrec int = 2
	// QuantityPrec - precition for quantity
	QuantityPrec int = 6
)

// SumPrice sum for price
func SumPrice(x, y float64) float64 {
	return sum(x, y, MoneyPrec)
}

// SumPrice2 sum for price
func SumPrice2(x, y float64) float64 {
	return round(x+y, MoneyPrec)
}

// SumQuantity sum for calc
func SumQuantity(x, y float64) float64 {
	return sum(x, y, QuantityPrec)
}

func sum(x, y float64, prec int) float64 {
	bx := big.NewFloat(x)
	by := big.NewFloat(y)
	f, _ := bx.SetMode(big.ToNearestAway).Add(bx, by).Float64()
	return round(f, prec)
}

func round(x float64, prec int) float64 {
	precLimiter := math.Pow10(prec)
	return math.Round(x*precLimiter) / precLimiter
}
