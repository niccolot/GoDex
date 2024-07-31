package main

import (
	"fmt"
	"strings"
)

func PrintPrompt() {

	fmt.Print("\n\nPokedex > ")
}

func PrintUnknown(text string) {

	fmt.Printf("'%s' command not found", text)
}

func CleanInput(text string) string {
	// removes trailing whitespaces and lowercases the command

	outCmd := strings.TrimSpace(text)
	outCmd = strings.ToLower(outCmd)
	return outCmd
}