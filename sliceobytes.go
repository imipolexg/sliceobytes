// sliceobytes generates the go code for a slice of bytes ([]byte)
// representing a file, encoded in base 16, suitable for embedding
// small amounts of data into source code, like unit tests.  Large
// amounts of data cause the go compiler to gobble up tremendous
// amounts of memory when compiling, so keep it under a few
// kilobytes. Could also easily be adapted to an array of C's char
// type, etc.
package main

import (
	"fmt"
	"io"
	"os"
)

const BUFMAX = 1024 // Look at us, so memory conscious!
const WIDTH = 8     // Print 8 bytes and then newline

func main() {
	args := os.Args[:]
	if len(args) < 2 {
		fmt.Printf("Bad mojo! Do %s <filename>\n", args[0])
		return
	}

	fi, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}

	// We need to know the file size so we can know when to stop
	// printing commas (a,b,c<--). I guess it's supposed to be a
	// rotten and nasty thing, but I wish golang had a ?: ternary
	// operator
	fInfo, err := fi.Stat()
	if err != nil {
		panic(err)
	}
	f_sz := fInfo.Size()

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Printf("[]byte{")

	buf := make([]byte, BUFMAX)
	pos := int64(0) // Our position in the file as a whole
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// i is the index only into the current buffer
		for i := 0; i < n; i++ {
			if pos < (f_sz - 1) {
				fmt.Printf("0x%x, ", buf[i])
			} else {
				fmt.Printf("0x%x", buf[i])
			}
			if pos == WIDTH-2 {
				fmt.Printf("\n")
			} else if pos > WIDTH && (pos+2)%WIDTH == 0 {
				fmt.Printf("\n")
			}
			pos++
		}
	}
	fmt.Printf("}\n")
}
