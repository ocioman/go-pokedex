package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"pokedexcli/cache"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/qeesung/image2ascii/convert"
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
	InspectPokemonName       string
	Cache                    *cache.Cache
	PokemonsBuffer           []Pokemon
	PokemonsChannel          chan any
	SaveFile                 *os.File
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

type PokemonType struct {
	Name string `json:"name"`
}
type TypeSlot struct {
	Ptype PokemonType `json:"type"`
}

type Sprites struct {
	FrontDefault string `json:"front_default"`
}

type Pokemon struct {
	Id             int        `json:"id"`
	Name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	IsDefault      bool       `json:"is_default"`
	Order          int        `json:"order"`
	Weight         int        `json:"weight"`
	Types          []TypeSlot `json:"types"`
	PSprites       Sprites    `json:"sprites"`
	Stats          []StatSlot `json:"stats"`
	aSCIIArt       string
}

type Stat struct {
	Name string `json:"name"`
}

type StatSlot struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	PStat    Stat `json:"stat"`
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

	fmt.Printf("\n")
	for _, l := range decodedRes.Results {
		fmt.Println(l.Name)
	}
	fmt.Printf("\n")

	cfg.LocationAreasUrlPrev = decodedRes.Previous
	cfg.LocationAreasUrlNext = decodedRes.Next

	return nil
}

func getPreviousLocationAreas(cfg *Config) error {
	if len(cfg.LocationAreasUrlPrev) == 0 {
		return fmt.Errorf("\nno previous area locations\n")
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

	fmt.Printf("\n")
	for _, l := range decodedRes.Results {
		fmt.Println(l.Name)
	}
	fmt.Printf("\n")

	cfg.LocationAreasUrlPrev = decodedRes.Previous
	cfg.LocationAreasUrlNext = decodedRes.Next

	return nil
}

func getPokemonsInLocationArea(cfg *Config) error {
	if cfg.PokemonAreaLocationParam == "" {
		return fmt.Errorf("\nno area specified\n")
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
		return fmt.Errorf("\nno results in the area %s\n", cfg.PokemonAreaLocationParam)
	}

	fmt.Printf("\n")
	for _, p := range decodedData.PokemonEncounters {
		fmt.Println(p.Poke.Name)
	}
	fmt.Printf("\n")

	cfg.PokemonAreaLocationParam = ""

	return nil
}

func catchPokemon(cfg *Config) error {
	var decodedBody Pokemon

	if cfg.PokemonToCatchParam == "" {
		return fmt.Errorf("\nno pokemon provided\n")
	}

	if _, ok := cfg.Pokemons[cfg.PokemonToCatchParam]; ok {
		return fmt.Errorf("\nyou already have this pokemon\n")
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

		defer func() {
			err = res.Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}()

		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("\nunknown pokemon\n")
		}

		resBody, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}

		err = json.Unmarshal(resBody, &decodedBody)
		cfg.Cache.Add(url, resBody)
	}

	isCaptured := randomWithBias(float64(decodedBody.BaseExperience))

	fmt.Printf("\nThrowing a Pokeball at %s...\n", cfg.PokemonToCatchParam)

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

		color.Red(fmt.Sprintf("%s escaped!\n\n", cfg.PokemonToCatchParam))

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

		color.Green("Gotcha!\n\n")

		var err error

		decodedBody.aSCIIArt, err = getASCIIart(decodedBody.PSprites.FrontDefault)

		if err != nil {
			return err
		}

		var empty any
		cfg.Pokemons[cfg.PokemonToCatchParam] = decodedBody
		cfg.PokemonsBuffer = append(cfg.PokemonsBuffer, decodedBody)
		cfg.PokemonsChannel <- empty
	}

	cfg.PokemonToCatchParam = ""

	return nil
}

func inspectPokemon(cfg *Config) error {
	if cfg.InspectPokemonName == "" {
		return fmt.Errorf("\nno pokemon provided\n")
	}

	if inspected, ok := cfg.Pokemons[cfg.InspectPokemonName]; ok {
		var err error

		if inspected.aSCIIArt == "" {
			inspected.aSCIIArt, err = getASCIIart(inspected.PSprites.FrontDefault)

			if err != nil {
				return err
			}
		}
		fmt.Print(inspected.aSCIIArt)
		fmt.Println("Name: ", inspected.Name)
		fmt.Println("Height: ", inspected.Height, "dm")
		fmt.Println("Weight: ", inspected.Weight, "hg")

		fmt.Println("Stats: ")

		for _, s := range inspected.Stats {
			fmt.Printf("\t-%s:", s.PStat.Name)
			fmt.Printf("\n\t\t-Base %s: %d", s.PStat.Name, s.BaseStat)
			fmt.Printf("\n\t\t-Effort: %d\n", s.Effort)
		}

		fmt.Println("Types: ")

		for _, t := range inspected.Types {
			fmt.Printf("\t-%s\n", t.Ptype.Name)
		}
		fmt.Print("\n")
	} else {
		return fmt.Errorf("\nyou haven't caught this pokemon yet\n")
	}

	cfg.InspectPokemonName = ""

	return nil
}

func getASCIIart(spriteUrl string) (string, error) {
	res, err := http.Get(spriteUrl)

	if err != nil {
		return "", err
	}

	defer func() {
		err = res.Body.Close()

		if err != nil {
			fmt.Println(err)
		}
	}()

	spriteImg, err := png.Decode(res.Body)

	spriteImg = imaging.Resize(spriteImg, 96, 72, imaging.Lanczos)

	if err != nil {
		return "", err
	}

	converter := convert.NewImageConverter()

	convertOption := convert.Options{
		FixedWidth:      96,
		FixedHeight:     72,
		FitScreen:       false,
		StretchedScreen: true,
		Colored:         true,
		Reversed:        false,
	}

	spriteASCII := converter.Image2ASCIIString(spriteImg, &convertOption)
	return spriteASCII, nil
}

func inspectPokedex(cfg *Config) error {
	if len(cfg.Pokemons) == 0 {
		return fmt.Errorf("\nyour pokedex is empty\n")
	}

	fmt.Println("\nYour Pokedex: ")

	for _, p := range cfg.Pokemons {
		fmt.Printf("\t-%s\n", p.Name)
	}

	fmt.Printf("\n")

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

	if len(cfg.PokemonsBuffer) > 0 {
		err := savePokedex(cfg.SaveFile, cfg)

		if err != nil {
			fmt.Println(err)
		}
	}

	err := cfg.SaveFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	close(cfg.PokemonsChannel)

	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Printf("\n")
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for _, v := range cfg.Commands {
		fmt.Printf("%s: %s\n", v.Name, v.Description)
	}

	fmt.Printf("\n")

	return nil
}

func savePokedex(outStream *os.File, cfg *Config) error {
	bw := bufio.NewWriter(outStream)

	defer func() {
		err := bw.Flush()

		if err != nil {
			fmt.Println(err)
		}
	}()

	stat, err := cfg.SaveFile.Stat()
	if err != nil {
		fmt.Println(err)
	}

	if stat.Size() == 0 {
		_, err = bw.WriteRune('[')

		if err != nil {
			return err
		}
	} else {
		err = cfg.SaveFile.Truncate(stat.Size() - 1)
		if err != nil {
			fmt.Println(err)
		}
		_, err = bw.WriteRune(',')

		if err != nil {
			return err
		}
	}

	var i int

	for _, p := range cfg.PokemonsBuffer {
		var encoded []byte

		encoded, err = json.Marshal(p)

		if err != nil {
			return err
		}

		_, err = bw.Write(encoded)

		if err != nil {
			return err
		}

		if i < len(cfg.PokemonsBuffer)-1 {
			_, err = bw.WriteRune(',')

			if err != nil {
				return err
			}
		}

		i++
	}

	_, err = bw.WriteRune(']')

	if err != nil {
		return err
	}

	return nil
}

func WritePokemonsBufferLoop(cfg *Config) error {
	var bufSize int

	for range cfg.PokemonsChannel {
		bufSize++
		if bufSize == 5 {
			err := savePokedex(cfg.SaveFile, cfg)
			cfg.PokemonsBuffer = cfg.PokemonsBuffer[:0]
			bufSize = 0
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func LoadSave(cfg *Config, save *os.File) error {
	decoder := json.NewDecoder(save)

	var decodedPokemons []Pokemon

	err := decoder.Decode(&decodedPokemons)

	if err != nil {
		return err
	}

	for _, p := range decodedPokemons {
		cfg.Pokemons[p.Name] = p
	}

	return nil
}

func GetCommands() map[string]CliCommand {
	commands := make(map[string]CliCommand)

	commands["exit"] = CliCommand{
		Name:        "exit",
		Description: "Exits the Pokedex",
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
		Name:        "explore <location_area>",
		Description: "Prints all the possible pokemon encounters in the location area",
		Callback:    getPokemonsInLocationArea,
	}

	commands["catch"] = CliCommand{
		Name:        "catch <pokemon>",
		Description: "Catches a pokemon",
		Callback:    catchPokemon,
	}

	commands["inspect"] = CliCommand{
		Name:        "inspect <pokemon>",
		Description: "Inspects a Pokemon in your Pokedex based on its name",
		Callback:    inspectPokemon,
	}

	commands["pokedex"] = CliCommand{
		Name:        "pokedex",
		Description: "Prints all the Pokemons in your Pokedex",
		Callback:    inspectPokedex,
	}
	return commands
}

func CleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Split(text, " ")
}
