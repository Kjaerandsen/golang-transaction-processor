package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestGenerateRandomTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	generateRandomTxs(1000)
	hash := generateFileHash("txs.txt")
	if hash == "&{[1053271813 2703628384 378130007 3008585538 1052052335 2538720549 704867873 2417958192] "+
		"[48 51 53 48 50 49 46 56 56 55 56 46 57 52 51 52 46 57 52 57 48 46 52 53 56 50 46 53 55 55 46 56 50 "+
		"50 48 46 54 50 52 46 55 51 55 56 46 53 52 51 48 46 56 56 53 54 46 55 52 53 48 46 57 53 55 46] 19 4755 "+
		"false}" {
		fmt.Println("YES")
	} else {
		t.Error("testGenerateRandomTxs did not generate the expected output.")
	}
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
