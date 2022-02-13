package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Takes flags to run the different functions.
func main() {
	rand.Seed(time.Now().UnixNano())

	timeStart := time.Now()
	generateRandomTxs(1000)
	fmt.Println(time.Since(timeStart))
	timeStart = time.Now()
	generateFees()
	fmt.Println(time.Since(timeStart))
	timeStart = time.Now()
	fmt.Println(sum())
	fmt.Println(time.Since(timeStart))
	timeStart = time.Now()
	earnings()
	fmt.Println(time.Since(timeStart))
	timeStart = time.Now()
	number1, number2 := compare()
	fmt.Println(number1, number2)
	fmt.Println(time.Since(timeStart))
	timeStart = time.Now()

	generateMillionTxs()
	fmt.Println(time.Since(timeStart))
	timeStart = time.Now()
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
		outputString = outputString + fmt.Sprintf("%v\n", math.Round((float64(rand.Intn(9999))/100)*100)/100)
	}

	// Write to the file
	_, err = outputFile.WriteString(outputString)
	if err != nil {
		fmt.Println("Error writing to txs.txt file.")
		return
	}
}

// Reads the txs.txt file, sums all the numbers and prints the result.
func sum() float64 {
	var sum float64
	var add float64

	// Open the file
	inputFile, err := os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return 0
	}

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// Add the line to the sum
		// Add the line to the add value, if invalid input return an error
		add, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: txs.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return 0
		}
		// Add the add value to the sum
		sum = sum + add
	}

	// Print the result
	return (math.Round(sum * 100)) / 100
}

// Reads the txs.txt file, generates fees for each row and puts them into the fees.txt file with the same structure.
func generateFees() {
	var fee float64

	// Open the input file
	inputFile, err := os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}

	// Open the output file
	outputFile, err := os.Create("fees.txt")
	if err != nil {
		fmt.Println("Error accessing fees.txt file.")
		return
	}
	// Close the file after used
	defer outputFile.Close()

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	// For each line generate the fee (30%)
	for scanner.Scan() {
		// Add the line to the fee value, if invalid input return an error
		fee, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: txs.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return
		}
		// Calculate the fee
		fee = math.Round(fee*0.30*100) / 100
		// Add it to the output document
		_, err = outputFile.WriteString(fmt.Sprintf("%v\n", fee))
		if err != nil {
			fmt.Println("Error writing to fees.txt file.")
			return
		}
	}
}

// Calculates the earnings of the app provider (70%) and puts it into the earnings.txt file.
func earnings() {
	var profit float64

	// Open the input file
	inputFile, err := os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}

	// Open the output file
	outputFile, err := os.Create("earnings.txt")
	if err != nil {
		fmt.Println("Error accessing earnings.txt file.")
		return
	}
	// Close the file after used
	defer outputFile.Close()

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	// For each line generate the earnings (70%)
	for scanner.Scan() {
		// Add the line to the profit value, if invalid input return an error
		profit, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: txs.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return
		}
		// Calculate the earnings
		profit = math.Round(profit*0.70*100) / 100
		// Add it to the output document
		_, err = outputFile.WriteString(fmt.Sprintf("%v\n", profit))
		if err != nil {
			fmt.Println("Error writing to earnings.txt file.")
			return
		}
	}
}

// Calculates two numbers, fees sum - fees total and total - total earnings + fees sum
func compare() (float64, float64) {
	var FEES_SUM float64
	var FEES_TOTAL float64
	var TOTAL_EARNINGS float64
	var inputVal float64

	// Open the input file
	inputFile, err := os.Open("fees.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return 0, 0
	}

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	// For each line read the value
	for scanner.Scan() {
		// Add the line to the inputVal, if invalid input return an error
		inputVal, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: txs.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return 0, 0
		}
		// Add the fee to the total
		FEES_SUM += math.Round(inputVal*0.70*100) / 100
	}
	// Close the file after used
	inputFile.Close()

	// Open the input file
	inputFile, err = os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return 0, 0
	}

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner = bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	// For each line read the value
	for scanner.Scan() {
		// Add the line to the inputVal, if invalid input return an error
		inputVal, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: txs.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return 0, 0
		}
		// Add the fee to the total
		FEES_TOTAL += math.Round(inputVal*0.70*100) / 100
	}
	// Close the file after used
	inputFile.Close()

	// Open the input file
	inputFile, err = os.Open("earnings.txt")
	if err != nil {
		fmt.Println("Error accessing earnings.txt file.")
		return 0, 0
	}

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner = bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	// For each line read the value
	for scanner.Scan() {
		// Add the line to the inputVal, if invalid input return an error
		inputVal, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: earnings.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return 0, 0
		}
		// Add the fee to the total
		FEES_SUM += math.Round(inputVal*0.70*100) / 100
	}
	// Close the file after used
	inputFile.Close()

	return math.Round((FEES_SUM-FEES_TOTAL)*100) / 100, // Number 1
		math.Round((FEES_TOTAL-(TOTAL_EARNINGS-FEES_SUM))*100) / 100 // Number 2
}

// Same as generateRandomTxs, but for a million values.
func generateMillionTxs() {
	outputFile, err := os.Create("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}
	// Close the file after used
	defer outputFile.Close()

	// Generate a million different values
	for i := 0; i < 1000000; i++ {
		// Write to the file
		_, err = outputFile.WriteString(
			fmt.Sprintf("%v\n", math.Round((float64(rand.Intn(9999))/100)*100)/100))
		if err != nil {
			fmt.Println("Error writing to txs.txt file.")
			return
		}
	}

}

// Generates a filehash from the input filename and returns it as a string
func generateFileHash(filename string) string {
	inputFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error accessing " + filename + " file.")
		return ""
	}

	// Create the sha256 hash
	fileHash := sha256.New()

	// Read line by line and calculate the fileHash
	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// Add the line to the fileHash
		fileHash.Write([]byte(scanner.Text()))
	}

	return fmt.Sprintf("%v", fileHash.Sum(nil))
}
