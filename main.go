package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Takes flags to run the different functions.
func main() {
	timeStart := time.Now()
	// Flag implementation:
	perf := flag.Bool("perf", false, "Time the execution of the program.")
	fees := flag.Bool("fees", false, "Run the generateFees function.")
	earn := flag.Bool("earn", false, "Run the earnings function.")
	comp := flag.Bool("comp", false, "Run the compare function.")
	help := flag.Bool("help", false, "Print out the available functions.")
	sumt := flag.Bool("sum", false, "Run the sumt function.")
	genm := flag.Bool("genm", false, "Run the generateMillionTxs function.")

	// For the generate flag
	var genValue int
	flag.IntVar(&genValue, "gen", 0, "Run the generateRandomTxs function with x transactions.")

	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	//defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	if *help {
		fmt.Println("This program uses flags to run the different commands, the following flags are availablle:\n" +
			"-perf : Time the execution of the program.\n" +
			"-fees : Calculates the fees (transactions * 30%) and outputs them into fees.txt\n" +
			"-earn : Calculates the earnings (transactions - 30% fees) and outputs them into earnings.txt\n" +
			"-comp : Calculates and prints out Number1 and Number2 (as specified in the readme)\n" +
			"-sumt  : Calculates the sum of all transactions and prints it out.\n" +
			"-help : Outputs this message\n" +
			"-gen x: Generate x transactions, one per line to the txs.txt file\n" +
			"-genm : Generates a million transactions, one per line to the txs.txt file\n" +
			"Running with no parameters generates a million transactions to txs.txt and runs the comp flag.")
	}

	// If invalid ignore the flag
	if genValue != 0 && genValue < 0 {
		genValue = 0
	} else {
		generateRandomTxs(genValue)
	}

	// If no flags run the generateMillionTxs function, compare function and show that the help parameter is available
	if !*fees && !*earn && !*comp && !*help && !*genm && !*sumt && genValue == 0 {
		generateMillionTxs()
		number1, number2 := compare()
		fmt.Println(number1, number2)
		fmt.Println("For help on using the program run the program with the -help parameter or refer to the readme")
	}

	if *fees {
		generateFees()
	}

	if *earn {
		earnings()
	}

	if *sumt {
		fmt.Println(sum())
	}

	if *comp {
		number1, number2 := compare()
		fmt.Println(number1, number2)
	}

	/*
		generateRandomTxs(1000000)
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
	*/

	// Print the time the program took to execute if the perf flag is used
	if *perf {
		fmt.Println(time.Since(timeStart))
	}
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
		// Write to the file
		_, err = outputFile.WriteString(
			fmt.Sprintf("%v\n", math.Round((float64(rand.Intn(9999))/100)*100)/100))
		if err != nil {
			fmt.Println("Error writing to txs.txt file.")
			return
		}
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
	feesSum := readFileAndSumLines("fees.txt")
	total := readFileAndSumLines("txs.txt")
	feesTotal := total * 0.3
	totalEarnings := readFileAndSumLines("earnings.txt")

	return math.Round((feesSum-feesTotal)*100) / 100, // Number 1
		math.Round((total-(totalEarnings-feesSum))*100) / 100 // Number 2
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
		// The input value is converted to a float, then later to string again to ignore all os specific line endings.
		inputVal, err := strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			fmt.Println("Error: txs.txt contains invalid values. Please refer to documentation to generate" +
				" a valid file.")
			return ""
		}
		// Add the line to the fileHash
		fileHash.Write([]byte(fmt.Sprintf("%v", inputVal)))
	}

	return fmt.Sprintf("%v", fileHash.Sum(nil))
}

// Takes an input filename, and reads the file line by line, adding the number on the line to a total sum
func readFileAndSumLines(filename string) float64 {
	var sum float64
	var inputVal float64
	inputFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error accessing " + filename + " file.")
		return 0
	}

	defer inputFile.Close()

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// Add the line to the inputVal, if invalid input return an error
		inputVal, err = strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			// TODO: Edit this error message
			fmt.Println("Error: " + filename + " contains invalid values. Please refer to documentation" +
				" to generate a valid file.")
			return 0
		}
		// Add the fee to the total
		sum += inputVal
	}

	return sum
}
