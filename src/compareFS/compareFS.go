package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFileOne(m map[string]int, fileName string) {
	f, err := os.Open(fileName)
	check(err)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Text() != "" {
			m[sc.Text()]--
		}
	}
	f.Close()
}

func compare(m map[string]int, fileName string) {
	f, err := os.Open(fileName)
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Text() != "" {
			m[sc.Text()]++
			if m[sc.Text()] == 0 {
				delete(m, sc.Text())
			}
		}
	}
}

func printResult(m map[string]int) {
	for a, i := range m {
		if i > 0 {
			fmt.Println("ADDED ", a)
		} else if i < 0 {
			fmt.Println("REMOVED ", a)
		}
	}
}

func fileFormat(fileName string) {
	var c = utf8.RuneCountInString(fileName)
	if fileName[c-4:] != ".txt" {
		panic("Wrong file format")
	}
}

func main() {
	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	old := f1.Bool("old", false, "take old data")
	neww := f1.Bool("new", false, "take new data")
	if len(os.Args) == 5 {
		f1.Parse(os.Args[1:])
		f1.Parse(os.Args[3:])
	} else {
		panic("Not correct args.")
	}

	if *old && *neww {
		if os.Args[1] != "--old" && os.Args[1] != "-old" {
			panic("Wrong order of args. Older is first.")
		} else if os.Args[1] == "-old" {
			panic("Use '--old', not '-old'.")
		} else if os.Args[3] == "-new" {
			panic("Use '--new', not '-new'.")
		}

		fileNameOld := os.Args[2]
		fileFormat(fileNameOld)
		fileNameNew := os.Args[4]
		fileFormat(fileNameNew)
		m := map[string]int{}

		readFileOne(m, fileNameOld)
		compare(m, fileNameNew)
		printResult(m)

	} else {
		panic("Use '--old' & '--new' flags for passing path to Args.")
	}
}
