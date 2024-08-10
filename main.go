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
		
		go line.AppendHistory(input)
		c.History = append(c.History, input)

		commandName, args := ParseInput(input)
		command, exists := cliCommandsTable[commandName]
		if exists {
			err := command.Callback(&c, args)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to execute command '%s': %w", commandName, err).Error())
			}
			if command.Name == "exit" {
				break
			}
		} else {
			PrintUnknown(commandName)
		}

		line.AppendHistory(input)
		c.History = append(c.History, input)
		fmt.Println()
	}
}
