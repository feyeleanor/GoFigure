package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const ERROR		= 2

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gofigure [file]\n")
	flag.PrintDefaults()
	os.Exit(ERROR)
}

func main() {
	var error os.Error
	var inode *os.Dir
	var file *os.File

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "no files were specified\n")
		os.Exit(ERROR)
	}

	for _, name := range flag.Args() {
		switch inode, error = os.Stat(name); {
			case error != nil:
			case inode.IsRegular():
				if file, error = os.Open(name, os.O_RDONLY, 0); error == nil {
					io.Copy(os.Stdout, file)
				}
		}
	}
	if error != nil {
		os.Exit(ERROR)
	}
	os.Exit(0)
}
