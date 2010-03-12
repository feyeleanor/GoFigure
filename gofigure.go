package main

import (
	"mustache"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const ERROR		= 2

// main operation modes
var write       = flag.Bool("w", false, "write result to (source) file instead of stdout")
var input_file	= flag.String("i", "", "specifies the source file instead of stdin")

type Dictionary map[string]string
func (d *Dictionary) define(rule string) {
	f := strings.Split(rule, "->", 0)
	if len(f) != 2 {
		fmt.Fprintf(os.Stderr, "%v\n", f)
		fmt.Fprintf(os.Stderr, "template rule must be of the form 'pattern -> replacement'\n")
		os.Exit(ERROR)
	}
	(*d)[strings.TrimSpace(f[0])] = strings.TrimSpace(f[1])
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gofigure [flags] [pattern -> replacement]\n")
	flag.PrintDefaults()
	os.Exit(ERROR)
}

func main() {
	var rendered_content string
	var template []uint8
	var error os.Error
	var inode *os.Dir

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "no template rules defined\n")
		os.Exit(ERROR)
	}

	rules := make(Dictionary)
	for i := 0; i < flag.NArg(); i++ {
		rules.define(flag.Args()[i])
	}

	if *input_file == "" {
		if template, error = ioutil.ReadAll(os.Stdin); error == nil {
			rendered_content, error = mustache.Render(string(template), rules)
		}
	} else {
		switch inode, error = os.Stat(*input_file); {
			case error != nil:
			case inode.IsRegular():
				rendered_content, error = mustache.RenderFile(*input_file, rules)
			default:
		}
	}
	if error != nil {
		os.Exit(ERROR)
	} else {
		if *write {
			error = ioutil.WriteFile(*input_file, []byte(rendered_content), 0)
		} else {
			_, error = os.Stdout.Write([]byte(rendered_content))
		}
		if error != nil { os.Exit(ERROR) }
	}
	os.Exit(0)
}
