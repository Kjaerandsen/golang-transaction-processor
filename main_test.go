package main

import (
	"math/rand"
	"testing"
)

// Tests the generateRandomTxs function by running it with a fixed seed
// and testing the outputfile hash for integrity of the output
func TestGenerateRandomTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	err := writeToFile(generateRandomTxs(1000), "txs.txt")
	if err != nil {
		t.Error("Error writing to txs.txt: " + err.Error())
	}

	hash := generateFileHash("txs.txt")
	if hash != "[6 243 251 141 40 164 224 157 194 123 85 134 150 93 230 192 64 124"+
		" 231 223 36 73 202 193 208 241 206 159 209 26 52 95]" {
		t.Error("testGenerateRandomTxs did not generate the expected output.")
	}
}

// Tests the sum function by checking if a predefined input gives the expected output
func TestSum(t *testing.T) {
	if sum() != 51005.23 {
		t.Error("testSum did not generate the expected output.")
	}
}

// Tests the generateFees function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestGenerateFees(t *testing.T) {
	generateFees()
	hash := generateFileHash("fees.txt")
	if hash != "[201 7 165 170 38 39 62 218 216 5 12 126 244 189 135 232 70 213"+
		" 112 215 194 167 65 93 115 122 253 86 59 4 143 38]" {
		t.Error("testGenerateFees did not generate the expected output.")
	}
}

// Tests the earnings function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestEarnings(t *testing.T) {
	earnings()
	hash := generateFileHash("earnings.txt")
	if hash != "[6 254 114 199 232 143 53 217 190 61 148 6 22 231 217 198"+
		" 233 94 80 201 27 6 65 103 190 177 169 153 139 189 188 6]" {
		t.Error("testEarnings did not generate the expected output.")
	}
}

// Tests the compare function with a predefined input checking if it gives the expected output
func TestCompare(t *testing.T) {
	number1, number2 := compare()
	if number1 != 0.08 || number2 != 30603.18 {
		t.Error("testCompare did not generate the expected output.")
	}
}

// Tests the generateMillionTxs function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestGenerateMillionTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	err := generateMillionTxs()
	if err != nil {
		t.Error("Error writing to txs.txt: " + err.Error())
	}

	hash := generateFileHash("txs.txt")
	if hash != "[240 176 187 255 29 141 235 56 144 2 225 71 118 163 19 180 228 118 1"+
		" 34 250 158 37 93 128 73 29 62 39 153 161 75]" {
		t.Error("testGenerateMillionTxs did not generate the expected output.")
	}
}
