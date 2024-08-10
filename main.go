package main

import (
	"fmt"
	"github.com/peterh/liner"
)

func main() {

	cliCommandsTable := getCliCommandsTable()
	c := *getInitConfig()

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	for {
		fmt.Println()
		input, err := line.Prompt("Pokedex > ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				break
			}
			fmt.Println("Error reading line:", err)
			continue
		}

		text := CleanInput(input)
		command, exists := cliCommandsTable[text]
		if exists {
			err := command.Callback(&c)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to execute command '%s': %w", text, err).Error())
			}
			if command.Name == "exit" {
				break
			}
		} else {
			PrintUnknown(text)
		}

		line.AppendHistory(text)
		c.History = append(c.History, text)
		fmt.Println()
	}
	
}
