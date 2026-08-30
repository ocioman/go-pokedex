package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"pokedexcli/cache"
	"strings"
	"time"
)

type Config struct {
	Commands                 map[string]CliCommand
	LocationAreasUrlPrev     string
	LocationAreasUrlNext     string
	PokemonAreaLocationUrl   string
	PokemonAreaLocationParam string
	PokemonToCatchUrl        string
	PokemonToCatchParam      string
	Pokemons                 map[string]Pokemon
	Cache                    *cache.Cache
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Locations struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type PokemonEncounter struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Encounter struct {
	Poke           PokemonEncounter `json:"pokemon"`
	VersionDetails []map[string]any `json:"version_details"`
}

type LocationArea struct {
	Id                   int              `json:"id"`
	Name                 string           `json:"name"`
	GameIndex            int              `json:"game_index"`
	EncounterMethodRates []map[string]any `json:"encounter_method_rates"`
	Locat                Location         `json:"location"`
	Names                []map[string]any `json:"names"`
	PokemonEncounters    []Encounter      `json:"pokemon_encounters"`
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
}

func getLocationAreas(cfg *Config) error {
	var decodedRes Locations

	if cached, ok := cfg.Cache.Get(cfg.LocationAreasUrlNext); ok {
		err := json.Unmarshal(cached, &decodedRes)

		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(cfg.LocationAreasUrlNext)

		if err != nil {
			return err
		}

		resBody, err := io.ReadAll(res.Body)

		defer func() {
			err = res.Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}()

		if err != nil {
			return err
		}

		err = json.Unmarshal(resBody, &decodedRes)

		if err != nil {
			return err
		}

		cfg.Cache.Add(cfg.LocationAreasUrlNext, resBody)
	}

	for _, l := range decodedRes.Results {
		fmt.Println(l.Name)
	}

	cfg.LocationAreasUrlPrev = decodedRes.Previous
	cfg.LocationAreasUrlNext = decodedRes.Next

	return nil
}

func getPreviousLocationAreas(cfg *Config) error {
	if len(cfg.LocationAreasUrlPrev) == 0 {
		return fmt.Errorf("no previous area locations")
	}

	var decodedRes Locations

	if cached, ok := cfg.Cache.Get(cfg.LocationAreasUrlPrev); ok {
		err := json.Unmarshal(cached, &decodedRes)

		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(cfg.LocationAreasUrlPrev)

		if err != nil {
			return err
		}

		resBody, err := io.ReadAll(res.Body)

		defer func() {
			err = res.Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}()

		if err != nil {
			return err
		}

		err = json.Unmarshal(resBody, &decodedRes)

		if err != nil {
			return err
		}

		cfg.Cache.Add(cfg.LocationAreasUrlPrev, resBody)
	}

	for _, l := range decodedRes.Results {
		fmt.Println(l.Name)
	}

	cfg.LocationAreasUrlPrev = decodedRes.Previous
	cfg.LocationAreasUrlNext = decodedRes.Next

	return nil
}

func getPokemonsInLocationArea(cfg *Config) error {
	if cfg.PokemonAreaLocationParam == "" {
		return fmt.Errorf("no area specified")
	}
	var decodedData LocationArea

	url := cfg.PokemonAreaLocationUrl + cfg.PokemonAreaLocationParam

	if cached, ok := cfg.Cache.Get(url); ok {
		err := json.Unmarshal(cached, &decodedData)

		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)

		if err != nil {
			return err
		}

		defer func() {
			err = res.Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}()

		resBody, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}

		err = json.Unmarshal(resBody, &decodedData)

		cfg.Cache.Add(url, resBody)
	}

	if len(decodedData.PokemonEncounters) == 0 {
		return fmt.Errorf("no results in the area %s", cfg.PokemonAreaLocationParam)
	}

	for _, p := range decodedData.PokemonEncounters {
		fmt.Println(p.Poke.Name)
	}

	return nil
}

func catchPokemon(cfg *Config) error {
	var decodedBody Pokemon

	if cfg.PokemonToCatchParam == "" {
		return fmt.Errorf("no pokemon provided")
	}

	if _, ok := cfg.Pokemons[cfg.PokemonToCatchParam]; ok {
		return fmt.Errorf("you already have this pokemon!")
	}

	url := cfg.PokemonToCatchUrl + cfg.PokemonToCatchParam + "/"

	if cached, ok := cfg.Cache.Get(url); ok {
		err := json.Unmarshal(cached, &decodedBody)

		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)

		if err != nil {
			return err
		}

		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("unknown pokemon")
		}

		defer func() {
			err = res.Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}()

		resBody, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}

		err = json.Unmarshal(resBody, &decodedBody)
		cfg.Cache.Add(url, resBody)
	}

	isCaptured := randomWithBias(float64(decodedBody.BaseExperience))

	fmt.Printf("Throwing a Pokeball at %s...\n", cfg.PokemonToCatchParam)

	ticker := time.NewTicker(1 * time.Second)
	var ticked int

	if isCaptured == 0 {
		ticks := rand.Intn(3) + 1

		ticked = 1

		for range ticker.C {
			fmt.Printf("%d...\n", ticked)
			if ticked >= ticks {
				ticker.Stop()
				break
			}
			ticked++
		}
		fmt.Printf("%s escaped!\n", cfg.PokemonToCatchParam)
	} else {
		ticked = 1

		for range ticker.C {
			fmt.Printf("%d...\n", ticked)
			if ticked >= 3 {
				ticker.Stop()
				break
			}
			ticked++
		}

		fmt.Println("Gotcha!")
		cfg.Pokemons[cfg.PokemonToCatchParam] = decodedBody
	}

	return nil
}

func randomWithBias(baseExp float64) int {
	bias := 1 / (1 + math.Exp(-0.01*(baseExp-100)))

	if rand.Float64() < bias {
		return 0
	}

	return 1
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for k, v := range cfg.Commands {
		fmt.Printf("%s: %s\n", k, v.Description)
	}

	fmt.Printf("\n")

	return nil
}

func GetCommands() map[string]CliCommand {
	commands := make(map[string]CliCommand)

	commands["exit"] = CliCommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    commandExit,
	}

	commands["help"] = CliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    commandHelp,
	}

	commands["map"] = CliCommand{
		Name:        "map",
		Description: "Prints the next 20 location areas",
		Callback:    getLocationAreas,
	}

	commands["mapb"] = CliCommand{
		Name:        "mapb",
		Description: "Prints the previous 20 location ares",
		Callback:    getPreviousLocationAreas,
	}

	commands["explore"] = CliCommand{
		Name:        "explore",
		Description: "Prints all the possible pokemon encounters in the location area",
		Callback:    getPokemonsInLocationArea,
	}

	commands["catch"] = CliCommand{
		Name:        "catch",
		Description: "Catches a pokemon",
		Callback:    catchPokemon,
	}
	return commands
}

func CleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Split(text, " ")
}
