package main

import (
	"math/rand"
	"testing"
)

func TestGenerateRandomTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	generateRandomTxs(1000)
	hash := generateFileHash("txs.txt")
	if hash != "[25 101 25 210 24 41 237 196 131 223 145 95 222 23 223 88 5 47 143"+
		" 244 10 47 32 8 116 104 56 191 220 138 227 18]" {
		t.Error("testGenerateRandomTxs did not generate the expected output.")
	}
}

func TestSum(t *testing.T) {
	if sum() != 50953.66 {
		t.Error("testSum did not generate the expected output.")
	}
}

func TestGenerateFees(t *testing.T) {
	generateFees()
	hash := generateFileHash("fees.txt")
	if hash != "[70 69 106 47 141 60 69 60 139 100 224 240 112 192 188 92 36 205 62"+
		" 122 158 245 29 211 30 82 48 5 246 67 10 211]" {
		t.Error("testGenerateFees did not generate the expected output.")
	}
}

func TestEarnings(t *testing.T) {
	earnings()
	hash := generateFileHash("earnings.txt")
	if hash != "[56 152 172 165 80 53 143 160 229 77 1 149 63 90 106 174 255 208 207"+
		" 35 219 207 202 212 19 253 242 175 229 159 31 157]" {
		t.Error("testEarnings did not generate the expected output.")
	}
}

func TestCompare(t *testing.T) {
	number1, number2 := compare()
	if number1 != -0.04 || number2 != 71335 {
		t.Error("testCompare did not generate the expected output.")
	}
}

func TestGenerateMillionTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	hash := generateFileHash("txs.txt")
	if hash != "[25 101 25 210 24 41 237 196 131 223 145 95 222 23 223 88 5 47 143 "+
		"244 10 47 32 8 116 104 56 191 220 138 227 18]" {
		t.Error("testGenerateMillionTxs did not generate the expected output.")
	}
}
