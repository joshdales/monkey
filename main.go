package main

import (
	"bufio"
	"flag"
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

func main() {
	flag.Parse()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Printf("Using the %s engine\n", *engine)
	fmt.Printf("Feel free to type in commands\n")

	scanner := bufio.NewScanner(os.Stdin)
	if *engine == "vm" {
		repl.StartVm(scanner, os.Stdout)
	} else {
		repl.StartEval(scanner, os.Stdout)
	}
}
