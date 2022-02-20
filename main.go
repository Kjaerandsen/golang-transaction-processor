package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Million          = 1000000
	feesFile         = "fees.txt"
	earningsFile     = "earnings.txt"
	transactionsFile = "txs.txt"
	earningsrate     = 7 // Multiplied by 10 to use with integer functions
	feesrate         = 3 // Multiplied by 10 to use with integer functions
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
	sumt := flag.Bool("sumt", false, "Run the sum function.")
	genm := flag.Bool("genm", false, "Run the generateMillionTxs function.")

	// For the generate flag
	var genValue int
	flag.IntVar(&genValue, "gen", 0, "Run the generateRandomTxs function with x transactions.")

	// Parse the flags
	flag.Parse()

	// Set the random seed to the current time
	rand.Seed(time.Now().UnixNano())

	if *help {
		fmt.Println("This program uses flags to run the different commands, the following flags are availablle:\n" +
			"-perf  : Time the execution of the program.\n" +
			"-fees  : Calculates the fees (transactions * 30%) and outputs them into fees.txt\n" +
			"-earn  : Calculates the earnings (transactions - 30% fees) and outputs them into earnings.txt\n" +
			"-comp  : Calculates and prints out Number1 and Number2 (as specified in the readme)\n" +
			"-sumt  : Calculates the sum of all transactions and prints it out.\n" +
			"-help  : Outputs this message\n" +
			"-gen=x : Generate x transactions, one per line to the txs.txt file\n" +
			"-genm  : Generates a million transactions, one per line to the txs.txt file\n" +
			"Running with no parameters generates a million transactions to txs.txt and runs the comp flag.")
	}

	// If invalid ignore the flag
	if genValue != 0 && genValue < 0 {
		genValue = 0
	} else if genValue != 0 {
		err := writeToFile(generateRandomTxs(genValue), transactionsFile)
		if err != nil {
			panic("Error: " + err.Error())
		}
	}

	// If no flags are provided run the default routine
	if !*fees && !*earn && !*comp && !*help && !*genm && !*sumt && genValue == 0 {
		err := generateMillionTxs()
		if err != nil {
			panic("Error: " + err.Error())
		}
		err = earnings()
		if err != nil {
			panic("Error: " + err.Error())
		}
		err = generateFees()
		if err != nil {
			panic("Error: " + err.Error())
		}
		//earnings()
		//generateFees()
		number1, number2 := compare()
		fmt.Println("Number 1: ", number1, "Number2: ", number2)
		fmt.Println("For help on using the program run the program with the -help parameter or refer to the readme")
	}

	if *fees {
		err := generateFees()
		if err != nil {
			panic("Error: " + err.Error())
		}
	}

	if *genm {
		err := generateMillionTxs()
		if err != nil {
			panic("Error: " + err.Error())
		}
	}

	if *earn {
		err := earnings()
		if err != nil {
			panic("Error: " + err.Error())
		}
	}

	if *sumt {
		fmt.Println("Sum: ", sum())
	}

	if *comp {
		number1, number2 := compare()
		fmt.Println("Number 1: ", number1, "Number2: ", number2)
	}

	// Print the time the program took to execute if the perf flag is used
	if *perf {
		fmt.Println("Time to run: ", time.Since(timeStart))
	}

}

// Takes n and generates n integers in an array, each containing a random number between 1 and 9999
// with a uniform random distribution
func generateRandomTxs(n int) []int {
	var randomTxs = make([]int, n)

	for i := 0; i < n; i++ {
		randomTxs[i] = rand.Intn(9999)
	}

	return randomTxs
}

// Takes an array of integers and writes it into a file as floats (the first two numbers are after the comma)
func writeToFile(n []int, filename string) error {
	var commaValues string
	// Create a buffer for writing each line
	outputLine := strings.Builder{}
	outputLine.Grow(7)

	// Create a buffer for writing to the file
	buf := strings.Builder{}
	// Size it to the contents being written
	// 7 as in characters per line max (line endings + contents should never exceed seven characters)
	buf.Grow(len(n) * 7)

	// Convert the integers to floats in strings and add them to the write-buffer
	for i := 0; i < len(n); i++ {
		outputLine.WriteString(strconv.Itoa(n[i]/100) + ".")
		commaValues = fmt.Sprintf("%02v", strconv.Itoa(n[i]%100))
		outputLine.WriteString(commaValues + "\n")
		buf.WriteString(outputLine.String())
		outputLine.Reset()
	}

	// Create the output file
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error accessing " + filename + " file.")
		return err
	}

	// Write the buffer to the file
	_, err = outputFile.WriteString(buf.String())

	// Close the file
	outputFile.Close()

	return err
}

// Reads the lines in a file and returns each line as an integer in an integer array (remove ".")
func readFromFile(filename string) ([]int, error) {
	var inputString string
	var i = 0
	var err error

	// Count the number of lines
	inputFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error accessing " + filename + " file.")
		return nil, err
	}
	lineCount, err := lineCounter(inputFile)
	if err != nil {
		fmt.Println("Error accessing " + filename + " file.")
		return nil, err
	}
	inputFile.Close()

	// Create an output array of the same size as the input file
	var outputArr = make([]int, lineCount)

	// Open again for reading the values
	inputFile, err = os.Open(filename)
	if err != nil {
		fmt.Println("Error accessing " + filename + " file.")
		return nil, err
	}

	defer inputFile.Close()

	// Reads the file line by line, using code from https://golangdocs.com/golang-read-file-line-by-line
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// Add the line to the inputVal, if invalid input return an error
		inputString = scanner.Text()
		// Remove the dot and convert to integer
		inputString = strings.Replace(inputString, ".", "", 1)
		// Add the int to the output string
		outputArr[i], err = strconv.Atoi(inputString)
		if err != nil {
			fmt.Println("Error: " + filename + " contains invalid values. Please refer to documentation" +
				" to generate a valid file.")
			return nil, err
		}
		i++
	}
	return outputArr, err
}

// Reads the txs.txt file, sums all the numbers and prints the result.
func sum() string {
	totalSum := readFileAndSumLines(transactionsFile)

	sumString := fmt.Sprintf("%v%s%v", totalSum/100, ".", totalSum%100)
	// Print the result
	return sumString
}

// Reads the txs.txt file calculates 30% fees on each transaction and writes the output to fees.txt
func generateFees() error {
	var err error
	var rounding int
	var fee int

	// Read the transaction data from file
	transactions, err := readFromFile(transactionsFile)
	if err != nil {
		fmt.Println("Error reading from txs.txt")
		return err
	}

	// Go through the list of transactions
	for i := 0; i < len(transactions); i++ {
		rounding = 0
		// Calculate the fee
		fee = transactions[i] * feesrate
		rounding = fee % 10
		// If .5 and even odd number round up
		if rounding == 5 {
			if (fee/10)%2 == 1 {
				fee += 10
			}
			// If more than .5 round up
		} else if rounding > 5 {
			fee += 10
		}

		transactions[i] = fee / 10
	}

	// write to file
	err = writeToFile(transactions, feesFile)
	if err != nil {
		fmt.Println("Error writing to fees.txt")
		return err
	}

	return err
}

// Reads the txs.txt file and calculates 70% earnings on each transaction, finally writes the output to earnings.txt
func earnings() error {
	var earning int
	var rounding int
	var err error

	// Read the transactions from file
	transactions, err := readFromFile(transactionsFile)
	if err != nil {
		fmt.Println("Error reading from txs.txt")
		return err
	}

	// Go through the list of transactions
	for i := 0; i < len(transactions); i++ {
		rounding = 0
		// Calculate the earnings
		earning = transactions[i] * earningsrate
		// If .5 and odd number round up
		rounding = earning % 10
		if rounding == 5 {
			if (earning/10)%2 == 1 {
				earning += 10
			}
			// If more than .5 round up
		} else if rounding > 5 {
			earning += 10
		}

		transactions[i] = earning / 10
	}

	// write to file
	err = writeToFile(transactions, earningsFile)
	if err != nil {
		fmt.Println("Error writing to earnings.txt")
		return err
	}

	return err
}

// Calculates two numbers, fees sum - fees total and total - total earnings + fees sum
func compare() (string, string) {
	var number1 string
	var number2 string
	feesSum := readFileAndSumLines(feesFile)
	total := readFileAndSumLines(transactionsFile)
	totalEarnings := readFileAndSumLines(earningsFile)
	var test int

	// Rounds the last digit according to bankers rounding
	totalMod10 := total % 10
	feesTotal := total * feesrate
	if totalMod10 == 5 {
		feesTotal = feesTotal/10 + 1
	} else if totalMod10 > 5 {
		feesTotal = feesTotal/10 + 1
	} else {
		feesTotal = feesTotal / 10
	}

	test = feesSum - feesTotal
	// If the value is negative, fix the formatting
	if test > 0 {
		number1 = fmt.Sprintf("%v.%v", test/100, test%100)
	} else {
		test = -test
		number1 = fmt.Sprintf("-%v.%v", test/100, test%100)
	}

	test = total - (totalEarnings + feesSum)
	if test > 0 {
		number2 = fmt.Sprintf("%v.%v", test/100, test%100)
	} else {
		test = -test
		number2 = fmt.Sprintf("-%v.%v", test/100, test%100)
	}

	return number1, number2
}

// Same as generateRandomTxs, but for a million values.
func generateMillionTxs() error {
	err := writeToFile(generateRandomTxs(Million), transactionsFile)
	return err
}

// Generates a hash from the input filenames contents and returns it as a string
func generateFileHash(filename string) string {
	var inputString string
	var inputVal int

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
		// Add the line to the inputVal, if invalid input return an error
		inputString = scanner.Text()
		// Remove the dot and convert to integer
		inputString = strings.Replace(inputString, ".", "", 1)
		inputVal, err = strconv.Atoi(inputString)
		if err != nil {
			fmt.Println("Error: " + filename + " contains invalid values. Please refer to documentation" +
				" to generate a valid file.")
			return ""
		}
		// Add the fee to the total
		// Add the line to the fileHash
		fileHash.Write([]byte(fmt.Sprintf("%v", inputVal)))
	}

	return fmt.Sprintf("%v", fileHash.Sum(nil))
}

// Takes an input filename, and reads the file line by line, adding the number on the line to a total sum and returns it
func readFileAndSumLines(filename string) int {
	var totalSum int
	var inputString string
	var inputVal int
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
		inputString = scanner.Text()
		// Remove the dot and convert to integer
		inputString = strings.Replace(inputString, ".", "", 1)
		inputVal, err = strconv.Atoi(inputString)
		if err != nil {
			fmt.Println("Error: " + filename + " contains invalid values. Please refer to documentation" +
				" to generate a valid file.")
			return 0
		}
		// Add the fee to the total
		totalSum = totalSum + inputVal
	}

	return totalSum
}

// Function that counts lines in a file, retrieved from:
// https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
