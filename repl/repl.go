package repl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Commands             map[string]CliCommand
	LocationAreasUrlPrev string
	LocationAreasUrlNext string
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

func getLocationAreas(cfg *Config) error {
	res, err := http.Get(cfg.LocationAreasUrlNext)

	if err != nil {
		return err
	}

	defer func() {
		err = res.Body.Close()

		if err != nil {
			fmt.Println(err)
		}
	}()

	decoder := json.NewDecoder(res.Body)

	var decodedRes Locations

	err = decoder.Decode(&decodedRes)

	if err != nil {
		return err
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
	res, err := http.Get(cfg.LocationAreasUrlPrev)

	if err != nil {
		return err
	}

	defer func() {
		err = res.Body.Close()

		if err != nil {
			fmt.Println(err)
		}
	}()

	decoder := json.NewDecoder(res.Body)

	var decodedRes Locations

	err = decoder.Decode(&decodedRes)

	if err != nil {
		return err
	}

	for _, l := range decodedRes.Results {
		fmt.Println(l.Name)
	}

	cfg.LocationAreasUrlPrev = decodedRes.Previous
	cfg.LocationAreasUrlNext = decodedRes.Next

	return nil
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

	return commands
}

func CleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Split(text, " ")
}
