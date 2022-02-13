package main

import (
	"math/rand"
	"testing"
)

// Tests the generateRandomTxs function by running it with a fixed seed
// and testing the outputfile hash for integrity of the output
func TestGenerateRandomTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	generateRandomTxs(1000)
	hash := generateFileHash("txs.txt")
	if hash != "[147 25 97 34 103 197 218 137 95 62 208 93 4 109 185 96 177 115 145"+
		" 33 121 2 89 244 114 233 57 114 166 8 247 156]" {
		t.Error("testGenerateRandomTxs did not generate the expected output.")
	}
}

// Tests the sum function by checking if a predefined input gives the expected output
func TestSum(t *testing.T) {
	if sum() != 50953.66 {
		t.Error("testSum did not generate the expected output.")
	}
}

// Tests the generateFees function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestGenerateFees(t *testing.T) {
	generateFees()
	hash := generateFileHash("fees.txt")
	if hash != "[196 66 18 200 137 14 208 100 162 32 169 160 24 177 19 249 173"+
		" 213 93 17 0 65 4 174 212 129 82 34 199 36 95 16]" {
		t.Error("testGenerateFees did not generate the expected output.")
	}
}

// Tests the earnings function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestEarnings(t *testing.T) {
	earnings()
	hash := generateFileHash("earnings.txt")
	if hash != "[186 146 193 253 129 199 7 183 184 172 157 72 163 29 145 "+
		"152 192 127 248 31 228 30 237 9 243 118 85 184 56 7 69 125]" {
		t.Error("testEarnings did not generate the expected output.")
	}
}

// Tests the compare function with a predefined input checking if it gives the expected output
func TestCompare(t *testing.T) {
	number1, number2 := compare()
	if number1 != 0.06 || number2 != 30572.3 {
		t.Error("testCompare did not generate the expected output.")
	}
}

// Tests the generateMillionTxs function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestGenerateMillionTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	hash := generateFileHash("txs.txt")
	if hash != "[147 25 97 34 103 197 218 137 95 62 208 93 4 109 185 96 177 115 145 "+
		"33 121 2 89 244 114 233 57 114 166 8 247 156]" {
		t.Error("testGenerateMillionTxs did not generate the expected output.")
	}
}
