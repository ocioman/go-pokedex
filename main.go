package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/cache"
	"pokedexcli/repl"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func main() {
	cfg := repl.Config{
		Commands:               repl.GetCommands(),
		LocationAreasUrlNext:   "https://pokeapi.co/api/v2/location-area/",
		Cache:                  cache.NewCache(6 * time.Minute),
		PokemonAreaLocationUrl: "https://pokeapi.co/api/v2/location-area/",
		Pokemons:               make(map[string]repl.Pokemon),
		PokemonToCatchUrl:      "https://pokeapi.co/api/v2/pokemon/",
		PokemonsBuffer:         make([]repl.Pokemon, 0),
		PokemonsChannel:        make(chan struct{}),
		Mu:                     new(sync.RWMutex),
	}

	if _, statErr := os.Stat("save.json"); statErr == nil {
		save, err := os.Open("save.json")

		if err != nil {
			fmt.Println(err)
		}

		err = repl.LoadSave(&cfg, save)

		if err = save.Close(); err != nil {
			fmt.Println(err)
		}

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(statErr)
	}

	defer func() {
		err := cfg.Commands["commandExit"].Callback(&cfg)

		if err != nil {
			fmt.Println(err)
		}
	}()

	var err error

	cfg.SaveFile, err = os.OpenFile("save.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		fmt.Println(err)
	}

	scann := bufio.NewScanner(os.Stdin)

	go func() {
		err = repl.WritePokemonsBufferLoop(&cfg)
		if err != nil {
			fmt.Println(err)
		}
	}()

	color.Red(`
▄▖ 
▌ ▛▌
▙▌▙▌`)

	fmt.Println(`▄▖
▌▌█▌▚▘
▙▘▙▖▞▖
`)

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

			err = command.Callback(&cfg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
