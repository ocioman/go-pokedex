package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/cache"
	"pokedexcli/repl"
	"time"
)

func main() {
	cfg := repl.Config{
		Commands:               repl.GetCommands(),
		LocationAreasUrlNext:   "https://pokeapi.co/api/v2/location-area/",
		Cache:                  cache.NewCache(10 * time.Second),
		PokemonAreaLocationUrl: "https://pokeapi.co/api/v2/location-area/",
	}

	scann := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")

		scann.Scan()

		input := scann.Text()
		args := repl.CleanInput(input)

		if len(args) > 0 {
			if len(args) > 1 && args[0] == "explore" {
				cfg.PokemonAreaLocationParam = args[1]
			}

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
