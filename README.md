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
Reads the file line by line and creates a hash of all the float values. This is used for unit testing.
readFileAndSumLines: 
Reads the file line by line and sums the float values of each line. This is used for the compare function internals and sum.
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

Seconds round of optimizations then focused on buffering the output, so instead of writing line by line, one float at a time i buffered 100 floats together in a string and then wrote that to the file. This reduced the runtime to around 400-500 milliseconds.

## Time-tracking notes:

The estimated difficulty was pretty accurate, but i was too lenient on the time estimates.
