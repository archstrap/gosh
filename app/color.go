package main

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

func Debug(color string, content any) {
	fmt.Printf("%s%v%s\n", color, content, Reset)
}
