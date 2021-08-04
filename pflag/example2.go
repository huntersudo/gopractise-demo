package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

var (
	flagvar = pflag.Int("flagname", 1234, "help message for flagname")
)

func main() {
	pflag.Parse()

	fmt.Printf("argument number is: %v\n", pflag.NArg())
	fmt.Printf("argument list is: %v\n", pflag.Args())
	fmt.Printf("the first argument is: %v\n", pflag.Arg(0))
}
// [going@dev pflag]$ go run example2.go -h
//Usage of /tmp/go-build862866068/b001/exe/example2:
//      --flagname int   help message for flagname (default 1234)
//pflag: help requested
//exit status 2