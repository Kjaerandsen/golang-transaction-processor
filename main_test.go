package main

import (
	"math/rand"
	"testing"
	"time"
)

// Tests the generateRandomTxs function by running it with a fixed seed
// and testing the outputfile hash for integrity of the output
func TestGenerateRandomTxs(t *testing.T) {
	rand.Seed(1099293902193012)
	err := writeToFile(generateRandomTxs(1000), transactionsFile)
	if err != nil {
		t.Error("Error writing to txs.txt: " + err.Error())
	}

	hash := generateFileHash(transactionsFile)
	if hash != "[232 177 111 35 146 228 191 133 121 113 151 98 106 220 145 95 48 217"+
		" 199 141 211 132 129 211 93 72 187 80 113 82 3 125]" {
		t.Error("testGenerateRandomTxs did not generate the expected output.")
	}

	// Set the random seed to a new random value
	rand.Seed(time.Now().UnixNano())
	// Generate and test progressively bigger numbers
	testRange := [3]int{100, 10000, Million}
	for i := 0; i < len(testRange); i++ {
		// Check if the values are within the allowed range
		transactions := generateRandomTxs(testRange[i])
		for x := 0; x < len(transactions); x++ {
			if transactions[x] > 10000 || transactions[x] < 0 {
				t.Error("TestGenerateRandomTxs generated an invalid number")
			}
		}
	}

}

// Tests the sum function by checking if a predefined input gives the expected output
func TestSum(t *testing.T) {
	if sum() != "50953.66" {
		t.Error("testSum did not generate the expected output.")
	}
}

// Tests the generateFees function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestGenerateFees(t *testing.T) {
	err := generateFees()
	if err != nil {
		t.Error("testGenerateFees did not generate the expected output")
	}
	// Check if the values are within the allowed range
	transactions, err := readFromFile(feesFile)
	if err != nil {
		t.Error("testGenerateFees, unable to read the file fees.txt")
	}
	for i := 0; i < len(transactions); i++ {
		if transactions[i] > 10000 || transactions[i] < 0 {
			t.Error("TestGenerateRandomTxs generated an invalid number")
		}
	}
	hash := generateFileHash(feesFile)
	if hash != "[84 30 54 38 142 63 87 254 171 62 56 149 214 171 218 104 228 96"+
		" 82 254 189 131 238 88 159 96 67 72 255 216 95 64]" {
		t.Error("testGenerateFees did not generate the expected output.")
	}
}

// Tests the earnings function by checking if a predefined input gives the expected output
// which is tested using a hash of the contents of the output file
func TestEarnings(t *testing.T) {
	err := earnings()
	if err != nil {
		t.Error("testEarnings did not generate the expected output")
	}
	// Check if the values are within the allowed range
	transactions, err := readFromFile(earningsFile)
	if err != nil {
		t.Error("testEarnings, unable to read the file earnings.txt")
	}
	for i := 0; i < len(transactions); i++ {
		if transactions[i] > 10000 || transactions[i] < 0 {
			t.Error("TestGenerateRandomTxs generated an invalid number")
		}
	}
	hash := generateFileHash(earningsFile)
	if hash != "[46 124 31 89 251 58 87 190 162 210 0 87 137 241 105 101 65"+
		" 48 18 38 206 101 155 128 161 48 117 143 20 161 205 149]" {
		t.Error("testEarnings did not generate the expected output.")
	}
}

// Tests the compare function with a predefined input checking if it gives the expected output
func TestCompare(t *testing.T) {
	number1, number2 := compare()
	if number1 != "0.5" || number2 != "-0.5" {
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

	hash := generateFileHash(transactionsFile)
	if hash != "[125 100 214 116 55 106 43 254 224 115 13 67 64 240 34 4 203 174 1 213"+
		" 121 218 59 163 179 103 220 250 90 139 15 224]" {
		t.Error("testGenerateMillionTxs did not generate the expected output.")
	}
}
