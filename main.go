package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/repl"
)

func main() {
	cfg := repl.Config{
		Commands:             repl.GetCommands(),
		LocationAreasUrlNext: "https://pokeapi.co/api/v2/location-area/?limit=20",
	}

	scann := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")

		scann.Scan()

		input := scann.Text()
		args := repl.CleanInput(input)

		if len(args) > 0 {
			command, ok := cfg.Commands[args[0]]

			if !ok {
				fmt.Println("unknown command")
				continue
			}

			err := command.Callback(&cfg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
