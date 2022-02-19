package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	sumt := flag.Bool("sumt", false, "Run the sumt function.")
	genm := flag.Bool("genm", false, "Run the generateMillionTxs function.")

	// For the generate flag
	var genValue int
	flag.IntVar(&genValue, "gen", 0, "Run the generateRandomTxs function with x transactions.")

	// Parse the flags
	flag.Parse()

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
		err := writeToFile(generateRandomTxs(genValue), "txs.txt")
		if err != nil {
			panic("Error: " + err.Error())
		}
	}

	// If no flags are provided run the default routine
	if !*fees && !*earn && !*comp && !*help && !*genm && !*sumt && genValue == 0 {
		/*err := generateMillionTxs()
		if err != nil {
			panic("Error: " + err.Error())
		}*/
		number1, number2 := compare()
		fmt.Println("Number 1: ", number1, "Number2: ", number2)
		fmt.Println("For help on using the program run the program with the -help parameter or refer to the readme")
	}

	if *fees {
		generateFees()
	}

	if *genm {
		err := generateMillionTxs()
		if err != nil {
			panic("Error: " + err.Error())
		}
	}

	if *earn {
		earnings()
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

func writeToFile(n []int, filename string) error {
	// Create a buffer for writing each line
	outputLine := strings.Builder{}
	outputLine.Grow(7)

	// Create a buffer for writing to the file
	buf := strings.Builder{}
	// Size it to the contents being written
	// 7 as in characters per line max (line endings + contents should never exceed seven characters)
	buf.Grow(len(n) * 7)

	// Convert the integers to floats in strings and add them to the write buffer
	for i := 0; i < len(n); i++ {
		outputLine.WriteString(strconv.Itoa(n[i]/100) + "." + strconv.Itoa(n[i]%100) + "\n")
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

// Reads the lines in a file and returns each line as an integer (remove ".")
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
	sum := readFileAndSumLines("txs.txt")

	sumString := fmt.Sprintf("%v%s%v", sum/100, ".", sum%100)
	// Print the result
	return sumString
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
func compare() (int, int) {
	feesSum := readFileAndSumLines("fees.txt")
	total := readFileAndSumLines("txs.txt")
	feesTotal := total * 3 / 10
	totalEarnings := readFileAndSumLines("earnings.txt")

	return feesSum - feesTotal, // Number 1
		total - (totalEarnings - feesSum) // Number 2
}

// Same as generateRandomTxs, but for a million values.
func generateMillionTxs() error {
	err := writeToFile(generateRandomTxs(1000000), "txs.txt")
	return err
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
