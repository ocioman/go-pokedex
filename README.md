# go-pokedex

## Italiano

CLI in Go che usa [PokeAPI](https://pokeapi.co/) per esplorare aree, incontrare Pokémon e gestire un piccolo Pokédex locale.

### Requisiti

- Go 1.26+

### Avvio

```bash
go run .
```

Dopo l'avvio, comparirà il prompt:

```text
Pokedex >
```

### Comandi disponibili

- `help`: mostra l'elenco dei comandi
- `exit`: chiude l'applicazione
- `map`: mostra le prossime aree
- `mapb`: mostra le aree precedenti
- `explore <area>`: mostra i Pokémon incontrabili in un'area
- `catch <pokemon>`: prova a catturare un Pokémon
- `inspect <pokemon>`: mostra i dettagli di un Pokémon catturato
- `pokedex`: mostra tutti i Pokémon catturati

### Esempio rapido

```text
Pokedex > map
Pokedex > explore canalave-city-area
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > pokedex
```

---

## English

Go CLI that uses [PokeAPI](https://pokeapi.co/) to explore areas, encounter Pokémon, and manage a small local Pokédex.

### Requirements

- Go 1.26+

### Run

```bash
go run .
```

After startup, the prompt will appear:

```text
Pokedex >
```

### Available commands

- `help`: show the list of commands
- `exit`: close the application
- `map`: show the next areas
- `mapb`: show the previous areas
- `explore <area>`: show encounterable Pokémon in an area
- `catch <pokemon>`: try to catch a Pokémon
- `inspect <pokemon>`: show details for a caught Pokémon
- `pokedex`: show all caught Pokémon

### Quick example

```text
Pokedex > map
Pokedex > explore canalave-city-area
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > pokedex
```
