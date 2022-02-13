package main

import (
	"math/rand"
	"testing"
)

func TestGenerateRandomTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	generateRandomTxs(1000)
}

func TestSum(t *testing.T) {
	sum()
}

func TestGenerateFees(t *testing.T) {
	generateFees()
}

func TestEarnings(t *testing.T) {
	earnings()
}

func TestCompare(t *testing.T) {
	compare()
}

func TestGenerateMillionTXs(t *testing.T) {
	rand.Seed(1099293902193012)

}
