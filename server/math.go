package server

import (
	"math/big"
)

func computeFibonacci(n int) *big.Int {
	a := big.NewInt(0)
	b := big.NewInt(1)
	for i := 0; i < n; i++ {
		a, b = b, new(big.Int).Add(a, b)
	}
	return a
}

func calculatePi(n int) float64 {
	pi := 0.0
	sign := 1.0

	for i := 0; i < n; i++ {
		term := 4.0 / float64(2*i+1)
		pi += sign * term
		sign *= -1
	}
	return pi
}

// check if is prime
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// count prime numbers up to limit
func countPrime(limit int) int {
	count := 0
	for i := 2; i <= limit; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}
