# Golang 02: simple currency processing
CLI program which deals with transaction processing.

## Usage:
To use the program run it in the command line with parameters deciding which functions to run.
The followin parameters are allowed:
- -perf  : Time the execution of the program.
- -fees  : Calculates the fees (transactions * 30%) and outputs them into fees.txt.
- -earn  : Calculates the earnings (transactions - 30% fees) and outputs them into earnings.txt.
- -comp  : Calculates and prints out Number1 and Number2 (as specified in the readme).
- -sumt  : Calculates the sum of all transactions and prints it out.
- -help  : Outputs this message.
- -gen=x : Generate x transactions, one per line to the txs.txt file. X here is an input integer
- -genm  : Generates a million transactions, one per line to the txs.txt file.
Running with no parameters generates a million transactions (runs generateMillionTxs function) to txs.txt and runs the comp flag.
It also prints out a message refering to the help flag for functionality instructions.

## Format:
All transactions are in euros with a minimum value of 0.01 and a maximum of 99.99

## Functions:
The program has x main functions:

1. generateRandomTxs:
Generates x (input integer) amount of transactions, one per line to the txs.txt file.
2. generateFees:
Reads the txs.txt file and generates a new file fees.txt showing the transaction fee for each transaction (30%), one per line.
3. earnings:
Reads the txs.txt file and generates a new file earnings.txt showing the earnings for each transaction (70%), one per line.
4. compare: 
Calculates two numbers: 
reads all the fees from fees.txt file, sums them up (FEES_SUM). Then, it reads all the txs from txs.txt, sums them up, and calculates the fee on the total sum (FEES_TOTAL). The function should return NUMBER1 = (FEES_SUM - FEES_TOTAL)
reads all the fees from fees.txt, sums them up (FEES_SUM). Reads all the earnings from earnings.txt and sums them up (TOTAL_EARNINGS). It also sums up all the transactions from txs.txt file (TOTAL). Then, it reports a number that is  NUMBER2 = TOTAL - (TOTAL_EARNINGS + FEES_SUM).
5. generateMillionTxs: 
Same as generateRandomTxs, but a set value of a million transactions are generated.
6. sum:
Reads the txs.txt file, sums all the numbers and writes it to the terminal.

It also has some helper functions:
generateFileHash: 
Reads the file line by line and creates a hash of all the values (as integers). This is used for unit testing.
readFileAndSumLines: 
Reads the file line by line and sums the float values of each line. This is used for the compare function internals and sum.
readFromFile: 
Reads the contents of the file, removes the end of line characters and ".", and returns the contents as an integer array.
writeToFile: 
Takes an integer array and writes it to file as floats (adds "." to the third last character)
lineCounter: 
Reads the file and returns the amount of lines in the file. This is used in the generateFees and earnings functions.

## Implementation notes:
The unit tests use a fixed seed to generate the files. The file contents are then hashed and compared to correct values.
The functions that return variables instead are compared to precalculated values from the defined seed.
Writes are buffered to 100 lines at a time, as i found the bottleneck got moved to the application at that point (string functions mostly).
File reading line by line uses code from: https://golangdocs.com/golang-read-file-line-by-line
The lineCounter function uses code from: https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently

## Profiling Notes:
The first major bottleneck was in generating transactions because my first implementation simply wrote all the transactions into a string and then wrote the string to the file. This didn't scale well.
The first optimization was then to write each transaction line by line instead.

Then using golang profiling i noted the following:
- generateMillionTxs: Took around 3 seconds, heavily bottlenecked by input output and system calls.
- generateFees:       Took around 3 seconds, heavily bottlenecked by input output and system calls.
- earnings:           Took around 3 seconds, heavily bottlenecked by input output and system calls.
- comp:               Took around 300 milliseconds, mostly bottlenecked by string functions, cpu and memory allocation. Fast enough.

Seconds round of optimizations then focused on buffering the output, so instead of writing line by line, one float at a time i buffered 100 floats together in a string and then wrote that to the file. This reduced the runtime to around 400-500 milliseconds. Which i found acceptable.
If i were to improve it further i would probably look into setting predefined variable sizes and using more basic containers for the buffering to reduce the amount of time spent reallocating memory for the string functions.

## Time-tracking notes:
Issues, time estimated and time spent:
1. Basic setup / outline of the application: Estimated 30m, spent 14m.
2. sum function:                             Estimated 30m, spent 17m.
3. generateFees function:                    Estimated 30m, spent 12m.
4. earnings function:                        Estimated 30m, spent 5m.
5. generateRandomTxs:                        Estimated 2h, spent 1h.
6. Performance profiling and optimizations:  Estimated 2h, spent 1h 25m.
7. compare function:                         Estimated 1h, spent 29m.
8. generateMillionTxs:                       Estimated 2h, spent 15m.
9. Unit tests:                               Estimated 1h, spent 1h 10m.
10. Documentation:                           Estimated 1h, spent 45m.
11. main function with flag implementation:  Estimated 30m, spent 30m.

Total time estimate: 11h 30m
Total time spent: 6h 22m
The estimated difficulty was pretty accurate, but i was too lenient on the time estimates.

## Second delivery notes:
The program now uses integers instead of floats, and runs much faster. A whole unit test round which includes writing a million records twice and all the functions, plus reading all the files takes less than a second on my computer.
New helper functions have been created which makes some functions a lot shorter in length, (readfromfile).
I had some trouble with rounding due to an error with writing to file (1.00 got written as 1 instead).
Consts are implemented for some values used in several places.
Unit tests now test more generate values.

## Time tracking for the second delivery:
1. generateRandomTxs function rewrite: Estimated 30m, spent 30m.
2. New write to file function with buffer: Estimated 45m, spent 45m.
3. Unit test update & upgrade: Estimated 30m, spent 15m.
4. Better rounding algorithm.: Estimated 45m, spent 1h 45m.
5. New ReadFromFile functions: Estimated 45m, spent 45m.
6. generateFees and earnings rewrite: Estimated 45m, spent 45m.
7. Final touches and consts: Estimated 1h, spent 30m.

Total time estimate: 5h
Total time spent:    5h 15m
Combined time spent: 11h 37m
