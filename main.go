package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"io"
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
		generateRandomTxs(genValue)
	}

	// If no flags are provided run the default routine
	if !*fees && !*earn && !*comp && !*help && !*genm && !*sumt && genValue == 0 {
		generateMillionTxs()
		number1, number2 := compare()
		fmt.Println(number1, number2)
		fmt.Println("For help on using the program run the program with the -help parameter or refer to the readme")
	}

	if *fees {
		generateFees()
	}

	if *genm {
		generateMillionTxs()
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

	// Print the time the program took to execute if the perf flag is used
	if *perf {
		fmt.Println(time.Since(timeStart))
	}
}

// Takes n and generates n rows in a text file, each containing a random number between 0.01 and 99.99
// with a uniform random distribution, this is written to the file txs.txt
func generateRandomTxs(n int) {
	transactionsLeft := n
	loopSize := 100
	var outputString string

	outputFile, err := os.Create("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}
	// Close the file after used
	defer outputFile.Close()

	// Write 100 lines at a time to reduce the number of system calls
	for {
		if transactionsLeft == 0 {
			break
		}

		if transactionsLeft < 100 {
			loopSize = transactionsLeft
			transactionsLeft = 0
		} else {
			transactionsLeft -= 100
		}

		for i := 0; i < loopSize; i++ {
			outputString = outputString +
				fmt.Sprintf("%v\n", math.Round((float64(rand.Intn(9999))/100)*100)/100)
		}

		_, err = outputFile.WriteString(outputString)
		if err != nil {
			fmt.Println("Error writing to txs.txt file.")
			return
		}
		outputString = ""
	}
}

// Reads the txs.txt file, sums all the numbers and prints the result.
func sum() float64 {
	sum := readFileAndSumLines("txs.txt")

	// Print the result
	return (math.Round(sum * 100)) / 100
}

// Reads the txs.txt file, generates fees for each row and puts them into the fees.txt file with the same structure.
func generateFees() {
	var fee float64
	var output string
	bufferSize := 100
	lineBuffer := 0

	// Open the input file
	inputFile, err := os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}

	// Count the number of lines
	lineCount, err := lineCounter(inputFile)
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
	}
	if lineCount < bufferSize {
		bufferSize = lineCount
	}
	inputFile.Close()

	// Open the input file again to read from the start
	inputFile, err = os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}

	defer inputFile.Close()

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
		// Add it to the output buffer
		output = output + fmt.Sprintf("%v\n", fee)
		lineBuffer = lineBuffer + 1
		if lineBuffer == bufferSize {
			// Add it to the output document
			_, err = outputFile.WriteString(output)
			if err != nil {
				fmt.Println("Error writing to fees.txt file.")
				return
			}
			// Resets the buffer, and sets it size to the smallest value between lineCount (lines left) or 100.
			lineCount -= lineBuffer
			lineBuffer = 0
			output = ""
			if lineCount < bufferSize {
				bufferSize = lineCount
			}
		}
	}
}

// Calculates the earnings of the app provider (70%) and puts it into the earnings.txt file.
func earnings() {
	var profit float64
	var output string
	bufferSize := 100
	lineBuffer := 0

	// Open the input file
	inputFile, err := os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}

	// Count the number of lines
	lineCount, err := lineCounter(inputFile)
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
	}
	if lineCount < bufferSize {
		bufferSize = lineCount
	}
	inputFile.Close()

	// Open the input file again to read from the start
	inputFile, err = os.Open("txs.txt")
	if err != nil {
		fmt.Println("Error accessing txs.txt file.")
		return
	}

	defer inputFile.Close()

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
		profit = math.Round(profit*70) / 100
		// Add it to the output buffer
		output = output + fmt.Sprintf("%v\n", profit)
		lineBuffer = lineBuffer + 1
		if lineBuffer == bufferSize {
			// Add it to the output document
			_, err = outputFile.WriteString(output)
			if err != nil {
				fmt.Println("Error writing to earnings.txt file.")
				return
			}
			// Resets the buffer, and sets it size to the smallest value between lineCount (lines left) or 100.
			lineCount -= lineBuffer
			lineBuffer = 0
			output = ""
			if lineCount < bufferSize {
				bufferSize = lineCount
			}
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
	generateRandomTxs(1000000)
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
	var totalSum float64
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
		totalSum += inputVal
	}

	return totalSum
}

// Function that counts lines, retrieved from:
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
