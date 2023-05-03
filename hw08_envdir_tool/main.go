package main

import (
	"log"
	"os"
)

var (
	argErrormsg = `not enough command line arguments.
	(Should be "./go-env path/to/env/ command arg1 arg2 ...")`

	errorEnvDirCode = 111
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal(argErrormsg)
	}
	dir := os.Args[1]
	environment, err := ReadDir(dir)
	if err != nil {
		os.Exit(errorEnvDirCode)
	}

	os.Exit(RunCmd(os.Args[2:], environment))
}
