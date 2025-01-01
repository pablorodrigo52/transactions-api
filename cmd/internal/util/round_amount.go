package util

import "math"

// RoundPurchaseAmount rounds the purchase amount to two decimal places (nearest cent)
func RoundPurchaseAmount(amount float32) float32 {
	return float32(math.Round(float64(amount*100)) / 100)
}
