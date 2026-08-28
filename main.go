package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/repl"
)

var commands map[string]repl.CliCommand

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.Description)
	}

	fmt.Printf("\n")

	return nil
}

func main() {
	commands = make(map[string]repl.CliCommand)

	commands["exit"] = repl.CliCommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    commandExit,
	}

	commands["help"] = repl.CliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    commandHelp,
	}

	scann := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")

		scann.Scan()

		input := scann.Text()
		args := repl.CleanInput(input)

		if len(args) > 0 {
			err := commands[args[0]].Callback()

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
