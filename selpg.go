package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

var (
	help               bool
	startPage, endPage int
	lineNum            int
	forcePage          bool
	destination        string
)

func init() {
	pflag.BoolVar(&help, "h", false, "help")
	pflag.IntVar(&startPage, "s", -1, "page to start printing")
	pflag.IntVar(&endPage, "e", -1, "page to end printing")
	pflag.IntVar(&lineNum, "l", 72, "number of lines per page")
	pflag.BoolVar(&forcePage, "f", false, "change page when meets \\f")
	pflag.StringVar(&destination, "d", "", "send to where to print")

	pflag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `selpg version: selpg/1.0
	Usage: go run selpg.go -sNumber -eNumber [-l] [-p prefix] [-g directives]
	
	Options:\n`)
	pflag.PrintDefaults()
}

func checkParameters() error {
	if startPage == -1 || endPage == -1 {
		return errors.New("please give start page and end page")
	}
	if startPage > endPage {
		return errors.New("start page number should be <= end page")
	}
	if forcePage && lineNum != 72 {
		return errors.New("can't force page when given line number")
	}
	return nil
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pflag.Parse()

	if help {
		pflag.Usage()
	}

	handle(checkParameters())

	var reader *bufio.Reader
	if pflag.NArg() == 0 {
		reader = bufio.NewReader(os.Stdin)

	} else {
		file, err := os.Open(pflag.Arg(0))
		handle(err)
		reader = bufio.NewReader(file)
	}

	var pageSepe byte
	if forcePage {
		pageSepe = '\f'
	} else {
		pageSepe = '\n'
	}
	cursor := 0
	var line string
	for {
		if line, _ = reader.ReadString(pageSepe); line == "" {
			break
		}
		if cursor >= (startPage-1)*lineNum && cursor <= (endPage)*lineNum {
			if destination == "" {
				fmt.Print(line)
			} else {
			}
		}
		cursor++
	}
	if destination != "" {
		cmd := exec.Command("lp", "-d"+destination)
		_, err := cmd.StdinPipe()
		handle(err)
	}
}
