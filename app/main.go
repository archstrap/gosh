package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// TODO: Uncomment the code below to pass the first stage
	fmt.Print("$ ")

	command, err := bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	commandName := strings.TrimRight(command, "\n")
	fmt.Printf("%s: command not found\n", commandName)

}
