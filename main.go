package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

// Takes flags to run the different functions.
func main() {
	generateRandomTxs(5)
}

// Takes n and generates n rows in a text file, each containing a random number between 0.01 and 99.99
// with a uniform random distribution, this is written to the file txs.txt
func generateRandomTxs(n int) {
	var outputString string
	// Check if the input number is one or greater.
	if n < 1 {
		fmt.Println("Error, input to generateRandomTxs to low. Please input a number greater than 0.")
		return
	}

	outputFile, err := os.Create("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}
	// Close the file after used
	defer outputFile.Close()

	// Generate the n different values
	for i := 0; i < n; i++ {
		if i != 1 {
			outputString = outputString + "\n"
		}
		outputString = outputString + fmt.Sprintf("%v", math.Round((float64(rand.Intn(9999))/1000)*100)/100)
	}

	// Write to the file
	_, err = outputFile.WriteString(outputString)
	if err != nil {
		fmt.Println("Error writing to txs.txt file.")
		return
	}
}

// Reads the txs.txt file, sums all the numbers and prints the result.
func sum() {

}

// Reads the txs.txt file, generates fees for each row and puts them into the fees.txt file with the same structure.
func generateFees() {

}

// Calculates the earnings of the app provider (70%) and puts it into the earnings.txt file.
func earnings() {

}

// Calculates two numbers, fees sum - fees total and total - total earnings + fees sum
func compare() {

}

// Same as generateRandomTxs, but for a million values.
func generateMillionTxs() {

}
