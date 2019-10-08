package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

var (
	help        bool
	startPage   int
	endPage     int
	lineNum     int
	forcePage   bool
	destination string
)

func init() {
	pflag.BoolVarP(&help, "h", "h", false, "help")
	pflag.IntVarP(&startPage, "s", "s", -1, "page to start printing")
	pflag.IntVarP(&endPage, "e", "e", -1, "page to end printing")
	pflag.IntVarP(&lineNum, "l", "l", 72, "number of lines per page")
	pflag.BoolVarP(&forcePage, "f", "f", false, "change page when meets \\f")
	pflag.StringVarP(&destination, "d", "d", "", "send to where to print")

	pflag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `selpg version: selpg/1.0
	Usage: go run selpg.go -sNumber -eNumber [-lNumber]|[-f] [-dDestionation]	
	Options:
`)
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
	var data string
	for {
		if line, _ = reader.ReadString(pageSepe); line == "" {
			break
		}
		if cursor >= (startPage-1)*lineNum && cursor < (endPage)*lineNum {
			if destination == "" {
				fmt.Print(line)
			} else {
				data += line
			}
		}
		cursor++
	}
	if destination != "" {
		cmd := exec.Command("lp", "-d"+destination)
		stdin, err := cmd.StdinPipe()
		handle(err)
		io.WriteString(stdin, data)
		err = cmd.Run()
		handle(err)
		stdin.Close()
	}
}
