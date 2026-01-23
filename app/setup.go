package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func loadShellRC(path string) error {

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Unable to open %s, please check whether %s exists or not.\n", path, path)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if "" == line || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], "'\"")
			os.Setenv(key, value)
		}

	}

	return nil
}
