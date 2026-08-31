package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/cache"
	"pokedexcli/repl"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	cfg := repl.Config{
		Commands:               repl.GetCommands(),
		LocationAreasUrlNext:   "https://pokeapi.co/api/v2/location-area/",
		Cache:                  cache.NewCache(10 * time.Second),
		PokemonAreaLocationUrl: "https://pokeapi.co/api/v2/location-area/",
		Pokemons:               make(map[string]repl.Pokemon),
		PokemonToCatchUrl:      "https://pokeapi.co/api/v2/pokemon/",
	}

	color.Red(`
▄▖ 
▌ ▛▌
▙▌▙▌`)

	fmt.Println(`▄▖
▌▌█▌▚▘
▙▘▙▖▞▖
`)

	scann := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")

		scann.Scan()

		input := scann.Text()
		args := repl.CleanInput(input)

		if len(args) > 0 {
			if len(args) > 1 && args[0] == "explore" {
				cfg.PokemonAreaLocationParam = strings.TrimSpace(args[1])
			} else if len(args) > 1 && args[0] == "catch" {
				cfg.PokemonToCatchParam = strings.TrimSpace(args[1])
			} else if len(args) > 1 && args[0] == "inspect" {
				cfg.InspectPokemonName = strings.TrimSpace(args[1])
			}

			command, ok := cfg.Commands[strings.TrimSpace(args[0])]

			if !ok {
				fmt.Println("\nunknown command\n")
				continue
			}

			err := command.Callback(&cfg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
