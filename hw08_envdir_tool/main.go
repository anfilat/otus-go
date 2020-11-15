package main

import (
	"fmt"
	"os"
)

// как в envdir.
const selfReturnCode = 111

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		printFatal("go-envdir: usage: go-envdir dir child")
	}
	env, err := ReadDir(args[0])
	if err != nil {
		printFatal(fmt.Errorf("go-envdir: fatal: %w", err))
	}
	code := RunCmd(args[1:], env)
	if code == selfReturnCode {
		printFatal(fmt.Errorf("go-envdir: unable to run %s", args[1]))
	}
	os.Exit(code)
}

func printFatal(a interface{}) {
	// log.Fatal не подходит, потому что добавляет перед сообщением дополнительную информацию,
	// а для утилиты это неправильное поведение
	_, _ = fmt.Fprintln(os.Stderr, a)
	os.Exit(selfReturnCode)
}
